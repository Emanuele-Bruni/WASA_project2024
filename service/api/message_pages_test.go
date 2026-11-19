package api_test

import (
	"WASA_project2024/service/database"
	"WASA_project2024/service/media"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"strings"
	"testing"
)

func testImageUpload() media.Upload {
	var encoded bytes.Buffer
	_ = png.Encode(&encoded, image.NewRGBA(image.Rect(0, 0, 2, 2)))
	return media.Upload{MimeType: "image/png", Data: base64.StdEncoding.EncodeToString(encoded.Bytes()), Name: "test.png"}
}

type pageResponse struct {
	Messages       []database.Message `json:"messages"`
	DeletedIDs     []int              `json:"deletedIds"`
	Cursor         string             `json:"cursor"`
	Before         *string            `json:"before"`
	HasMore        bool               `json:"hasMore"`
	HasMoreChanges bool               `json:"hasMoreChanges"`
}

func readPage(t *testing.T, payload []byte) pageResponse {
	t.Helper()
	var page pageResponse
	if err := json.Unmarshal(payload, &page); err != nil {
		t.Fatal(err)
	}
	return page
}
func TestMessagePaginationAndIncrementalUpdates(t *testing.T) {
	env := newTestAPI(t)
	alice, bob, carol := env.login("alice"), env.login("bobby"), env.login("carol")
	chat := env.createConversation(alice, bob)
	path := fmt.Sprintf("/conversations/%d/messages", chat)
	ids := make([]int, 0, 123)
	for i := 0; i < 123; i++ {
		_, id, err := env.db.CreateMessage(fmt.Sprintf("message %d", i), chat, alice)
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, id)
	}
	response, payload := env.request("GET", path, alice, nil)
	env.expectStatus(response, payload, 200)
	first := readPage(t, payload)
	if len(first.Messages) != 50 || !first.HasMore || first.Before == nil || first.Messages[0].MessageId != ids[73] {
		t.Fatalf("invalid first page: count=%d", len(first.Messages))
	}
	// Removing the anchor must not invalidate the opaque timestamp/ID cursor.
	if err := env.db.RemoveMessage(ids[73]); err != nil {
		t.Fatal(err)
	}
	response, payload = env.request("GET", path+"?before="+*first.Before, alice, nil)
	env.expectStatus(response, payload, 200)
	second := readPage(t, payload)
	if len(second.Messages) != 50 || second.Messages[0].MessageId != ids[23] || !second.HasMore {
		t.Fatal("older page is wrong")
	}
	response, payload = env.request("GET", path+"?before="+*second.Before, alice, nil)
	env.expectStatus(response, payload, 200)
	third := readPage(t, payload)
	if len(third.Messages) != 23 || third.HasMore || third.Before != nil {
		t.Fatal("last page is wrong")
	}
	seen := map[int]bool{}
	for _, page := range []pageResponse{first, second, third} {
		for _, message := range page.Messages {
			if seen[message.MessageId] {
				t.Fatal("duplicate across pages")
			}
			seen[message.MessageId] = true
		}
	}
	_, newID, err := env.db.CreateMessage("new", chat, bob)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err = env.db.SetReaction(bob, "👍", ids[122]); err != nil {
		t.Fatal(err)
	}
	if _, _, err = env.db.SetReaction(bob, "✨", ids[0]); err != nil {
		t.Fatal(err)
	}
	cursor := first.Cursor
	changed := map[int]bool{}
	deleted := map[int]bool{}
	batches := 0
	for {
		response, payload = env.request("GET", path+"?since="+cursor+"&limit=2", alice, nil)
		env.expectStatus(response, payload, 200)
		delta := readPage(t, payload)
		batches++
		for _, message := range delta.Messages {
			changed[message.MessageId] = true
		}
		for _, id := range delta.DeletedIDs {
			deleted[id] = true
		}
		cursor = delta.Cursor
		if !delta.HasMoreChanges {
			break
		}
		if batches > 10 {
			t.Fatal("changes do not advance")
		}
	}
	if batches < 2 || len(changed) != 3 || !changed[newID] || !changed[ids[0]] || !changed[ids[122]] || !deleted[ids[73]] {
		t.Fatalf("lost changes: changed=%v deleted=%v", changed, deleted)
	}
	response, payload = env.request("GET", path+"?since="+cursor, alice, nil)
	env.expectStatus(response, payload, 200)
	empty := readPage(t, payload)
	if len(empty.Messages) != 0 || len(empty.DeletedIDs) != 0 || empty.HasMoreChanges {
		t.Fatal("idle polling redownloaded history")
	}
	other := env.createConversation(alice, carol)
	for _, query := range []string{"?limit=0", "?limit=101", "?limit=1&limit=2", "?unknown=1", "?since=-1", "?since=9999999999999", "?before=bad", "?before=" + *first.Before + "&since=0"} {
		response, payload = env.request("GET", path+query, alice, nil)
		env.expectStatus(response, payload, 400)
	}
	response, payload = env.request("GET", fmt.Sprintf("/conversations/%d/messages?before=%s", other, *first.Before), alice, nil)
	env.expectStatus(response, payload, 400)
}
func TestImagesAreValidatedSharedAndProtected(t *testing.T) {
	env := newTestAPI(t)
	alice, bob, carol := env.login("alice"), env.login("bobby"), env.login("carol")
	chat := env.createConversation(alice, bob)
	path := fmt.Sprintf("/conversations/%d/messages", chat)
	upload := testImageUpload()
	response, payload := env.request("POST", path, alice, map[string]interface{}{"type": "image", "image": upload})
	env.expectStatus(response, payload, 201)
	var sent struct {
		Message database.Message `json:"message"`
	}
	if err := json.Unmarshal(payload, &sent); err != nil {
		t.Fatal(err)
	}
	if sent.Message.Type != "image" || sent.Message.Image == nil || sent.Message.Text != "" {
		t.Fatal("image metadata missing")
	}
	if bytes.Contains(payload, []byte(upload.Data)) {
		t.Fatal("binary image embedded in JSON")
	}
	url := sent.Message.Image.URL
	response, binary := env.request("GET", url, bob, nil)
	env.expectStatus(response, binary, 200)
	original, _ := base64.StdEncoding.DecodeString(upload.Data)
	if !bytes.Equal(binary, original) || response.Header.Get("Cache-Control") != "private, max-age=31536000, immutable" || response.Header.Get("Vary") != "Authorization" {
		t.Fatal("invalid media or cache contract")
	}
	request, _ := http.NewRequest("GET", env.server.URL+url, nil)
	request.Header.Set("Authorization", "Bearer "+env.tokens[bob])
	request.Header.Set("If-None-Match", response.Header.Get("ETag"))
	conditional, err := env.server.Client().Do(request)
	if err != nil {
		t.Fatal(err)
	}
	data, _ := io.ReadAll(conditional.Body)
	conditional.Body.Close()
	env.expectStatus(conditional, data, 304)
	response, payload = env.request("GET", url, carol, nil)
	env.expectStatus(response, payload, 403)
	response, payload = env.request("GET", url, 0, nil)
	env.expectStatus(response, payload, 401)
	response, payload = env.request("GET", "/media/999999", alice, nil)
	env.expectStatus(response, payload, 404)
	destination := env.createConversation(alice, carol)
	response, payload = env.request("POST", fmt.Sprintf("%s/%d/forward/%d", path, sent.Message.MessageId, destination), alice, nil)
	env.expectStatus(response, payload, 201)
	var forwarded struct {
		Message database.Message `json:"message"`
	}
	if err = json.Unmarshal(payload, &forwarded); err != nil {
		t.Fatal(err)
	}
	if forwarded.Message.Image.ImageID != sent.Message.Image.ImageID || !forwarded.Message.Forwarded {
		t.Fatal("forward duplicated image or lost metadata")
	}
	if err = env.db.RemoveMessage(sent.Message.MessageId); err != nil {
		t.Fatal(err)
	}
	response, payload = env.request("GET", url, carol, nil)
	env.expectStatus(response, payload, 200)
	bad := upload
	bad.MimeType = "image/jpeg"
	for _, body := range []interface{}{
		map[string]interface{}{"type": "image", "image": bad},
		map[string]interface{}{"type": "text", "text": "text", "image": upload},
		map[string]string{"type": "image"}, map[string]string{"text": "missing type"},
		map[string]interface{}{"type": "image", "text": "caption", "image": upload},
		map[string]interface{}{"type": "image", "image": map[string]string{"mimeType": "image/png", "data": upload.Data, "unknown": "field"}},
	} {
		response, payload = env.request("POST", path, alice, body)
		env.expectStatus(response, payload, 400)
	}
	response, payload = env.request("POST", path, alice, map[string]string{"type": "text", "text": "[foto]ordinary text"})
	env.expectStatus(response, payload, 201)
	if !strings.Contains(string(payload), `"type":"text"`) {
		t.Fatal("legacy prefix still changes message type")
	}
}

func TestConversationPreviewsFollowMessageActivity(t *testing.T) {
	env := newTestAPI(t)
	alice, bob, carol := env.login("alice"), env.login("bobby"), env.login("carol")
	first, second := env.createConversation(alice, bob), env.createConversation(alice, carol)
	_, id, err := env.db.CreateMessage(strings.Repeat("x", 1000), first, alice)
	if err != nil {
		t.Fatal(err)
	}
	response, payload := env.request("GET", "/conversations", alice, nil)
	env.expectStatus(response, payload, 200)
	var list struct {
		Conversations []database.Conversation `json:"conversations"`
	}
	if err = json.Unmarshal(payload, &list); err != nil {
		t.Fatal(err)
	}
	if list.Conversations[0].ConversationId != first || list.Conversations[0].LastMessage.MessageID != id || len(list.Conversations[0].LastMessage.Text) != 200 || list.Conversations[1].LastMessage != nil {
		t.Fatal("preview/order incorrect")
	}
	image, err := media.Decode(testImageUpload())
	if err != nil {
		t.Fatal(err)
	}
	_, imageID, err := env.db.CreateImageMessage(image, second, alice)
	if err != nil {
		t.Fatal(err)
	}
	response, payload = env.request("GET", "/conversations", alice, nil)
	env.expectStatus(response, payload, 200)
	if err = json.Unmarshal(payload, &list); err != nil {
		t.Fatal(err)
	}
	if list.Conversations[0].ConversationId != second || list.Conversations[0].LastMessage.Type != "image" || list.Conversations[0].LastMessage.Text != "" {
		t.Fatal("image preview incorrect")
	}
	if err = env.db.RemoveMessage(imageID); err != nil {
		t.Fatal(err)
	}
	response, payload = env.request("GET", "/conversations", alice, nil)
	env.expectStatus(response, payload, 200)
	if err = json.Unmarshal(payload, &list); err != nil {
		t.Fatal(err)
	}
	if list.Conversations[0].ConversationId != first || list.Conversations[1].LastMessage != nil {
		t.Fatal("deleted latest message remained in preview")
	}
}

func TestImageDeduplicationDoesNotShareFilenames(t *testing.T) {
	env := newTestAPI(t)
	alice, bob := env.login("alice"), env.login("bobby")
	chat := env.createConversation(alice, bob)
	path := fmt.Sprintf("/conversations/%d/messages", chat)
	first := testImageUpload()
	first.Name = "first-private-name.png"
	second := first
	second.Name = "my-own-name.png"
	var one, two struct {
		Message database.Message `json:"message"`
	}
	response, payload := env.request("POST", path, alice, map[string]interface{}{"type": "image", "image": first})
	env.expectStatus(response, payload, 201)
	if err := json.Unmarshal(payload, &one); err != nil {
		t.Fatal(err)
	}
	response, payload = env.request("POST", path, bob, map[string]interface{}{"type": "image", "image": second})
	env.expectStatus(response, payload, 201)
	if err := json.Unmarshal(payload, &two); err != nil {
		t.Fatal(err)
	}
	if one.Message.Image.ImageID != two.Message.Image.ImageID || one.Message.Image.Name != first.Name || two.Message.Image.Name != second.Name {
		t.Fatal("binary deduplication leaked or overwrote a per-message filename")
	}
}
