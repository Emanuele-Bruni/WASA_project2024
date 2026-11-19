package api_test

import (
	"WASA_project2024/service/api"
	"WASA_project2024/service/database"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

type testAPI struct {
	t      *testing.T
	server *httptest.Server
	db     database.AppDatabase
	dbconn *sql.DB
	users  map[string]int
	tokens map[int]string
}

func newTestAPI(t *testing.T) *testAPI {
	t.Helper()
	dbconn, err := database.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open temporary database: %v", err)
	}
	db, err := database.New(dbconn)
	if err != nil {
		dbconn.Close()
		t.Fatalf("initialize temporary database: %v", err)
	}
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	router, err := api.New(api.Config{Logger: logger, Database: db})
	if err != nil {
		dbconn.Close()
		t.Fatalf("initialize API: %v", err)
	}

	env := &testAPI{t: t, server: httptest.NewServer(router.Handler()), db: db, dbconn: dbconn, users: make(map[string]int), tokens: make(map[int]string)}
	t.Cleanup(func() {
		env.server.Close()
		_ = router.Close()
		_ = dbconn.Close()
	})
	return env
}

func (env *testAPI) request(method, path string, userID int, body interface{}) (*http.Response, []byte) {
	env.t.Helper()
	var requestBody io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			env.t.Fatalf("encode request: %v", err)
		}
		requestBody = bytes.NewReader(encoded)
	}
	req, err := http.NewRequest(method, env.server.URL+path, requestBody)
	if err != nil {
		env.t.Fatalf("create request: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if userID != 0 {
		req.Header.Set("Authorization", "Bearer "+env.tokens[userID])
	}

	resp, err := env.server.Client().Do(req)
	if err != nil {
		env.t.Fatalf("perform request: %v", err)
	}
	defer resp.Body.Close()
	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		env.t.Fatalf("read response: %v", err)
	}
	return resp, payload
}

func (env *testAPI) expectStatus(resp *http.Response, payload []byte, want int) {
	env.t.Helper()
	validateContractResponse(env.t, resp, payload)
	if resp.StatusCode != want {
		env.t.Fatalf("unexpected status: got %d, want %d; body=%q", resp.StatusCode, want, payload)
	}
	if want == 204 && len(payload) != 0 {
		env.t.Fatalf("204 response must be empty: %s", payload)
	}
	if want >= 400 {
		if resp.Header.Get("Content-Type") != "application/json" {
			env.t.Fatalf("non-JSON error: %s", payload)
		}
		var result struct {
			Error struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.Unmarshal(payload, &result); err != nil || result.Error.Code == "" || result.Error.Message == "" {
			env.t.Fatalf("invalid error envelope: %s", payload)
		}
	}

}

func (env *testAPI) login(name string) int {
	env.t.Helper()
	resp, payload := env.request(http.MethodPost, "/login", 0, map[string]string{"name": name})
	env.expectStatus(resp, payload, http.StatusOK)
	var result struct {
		User  database.User `json:"user"`
		Token string        `json:"token"`
	}
	if err := json.Unmarshal(payload, &result); err != nil {
		env.t.Fatalf("decode login response: %v", err)
	}
	if result.User.Id == 0 {
		env.t.Fatal("login returned an empty identifier")
	}
	if len(result.Token) != 43 {
		env.t.Fatal("login did not return a session token")
	}
	env.tokens[result.User.Id] = result.Token
	env.users[name] = result.User.Id
	return result.User.Id
}

func (env *testAPI) createConversation(authorID, recipientID int) int {
	env.t.Helper()
	resp, payload := env.request(http.MethodPost, fmt.Sprintf("/users/%d", recipientID), authorID, nil)
	env.expectStatus(resp, payload, http.StatusCreated)
	var result struct {
		Conversation database.Conversation `json:"conversation"`
	}
	if err := json.Unmarshal(payload, &result); err != nil {
		env.t.Fatalf("decode conversation response: %v", err)
	}
	return result.Conversation.ConversationId
}

func TestLoginUsersAndConversations(t *testing.T) {
	env := newTestAPI(t)
	alice := env.login("alice")
	bob := env.login("bobby")

	if repeated := env.login("alice"); repeated != alice {
		t.Fatalf("repeated login changed identifier: got %d, want %d", repeated, alice)
	}

	resp, payload := env.request(http.MethodGet, "/users", alice, nil)
	env.expectStatus(resp, payload, http.StatusOK)
	var users struct {
		Users []database.User `json:"users"`
	}
	if err := json.Unmarshal(payload, &users); err != nil {
		t.Fatalf("decode users response: %v", err)
	}
	if len(users.Users) != 1 || users.Users[0].Id != bob {
		t.Fatalf("unexpected users list: %#v", users.Users)
	}

	conversationID := env.createConversation(alice, bob)
	resp, payload = env.request(http.MethodGet, "/conversations", alice, nil)
	env.expectStatus(resp, payload, http.StatusOK)
	if !bytes.Contains(payload, []byte(fmt.Sprintf(`"conversationId":%d`, conversationID))) {
		t.Fatalf("created conversation missing from list: %s", payload)
	}

	resp, payload = env.request(http.MethodPost, fmt.Sprintf("/users/%d", bob), alice, nil)
	env.expectStatus(resp, payload, http.StatusConflict)
}

func TestGroupMessageForwardAndReactionFlows(t *testing.T) {
	env := newTestAPI(t)
	alice := env.login("alice")
	bob := env.login("bobby")
	carol := env.login("carol")

	firstConversation := env.createConversation(alice, bob)
	secondConversation := env.createConversation(alice, carol)

	resp, payload := env.request(http.MethodPost, "/conversations", alice, map[string]interface{}{
		"title":   "Team",
		"userIds": []int{bob, carol},
	})
	env.expectStatus(resp, payload, http.StatusCreated)
	if !bytes.Contains(payload, []byte(`"isGroup":true`)) {
		t.Fatalf("group response does not identify a group: %s", payload)
	}

	resp, payload = env.request(http.MethodPost, fmt.Sprintf("/conversations/%d/messages", firstConversation), alice, map[string]string{"type": "text", "text": "hello"})
	env.expectStatus(resp, payload, http.StatusCreated)
	var sent struct {
		Message database.Message `json:"message"`
	}
	if err := json.Unmarshal(payload, &sent); err != nil {
		t.Fatalf("decode message response: %v", err)
	}
	if sent.Message.MessageId == 0 || sent.Message.Text != "hello" {
		t.Fatalf("unexpected sent message: %#v", sent.Message)
	}

	forwardPath := fmt.Sprintf("/conversations/%d/messages/%d/forward/%d", firstConversation, sent.Message.MessageId, secondConversation)
	resp, payload = env.request(http.MethodPost, forwardPath, alice, nil)
	env.expectStatus(resp, payload, http.StatusCreated)

	reactionPath := fmt.Sprintf("/conversations/%d/messages/%d/comments", firstConversation, sent.Message.MessageId)
	resp, payload = env.request(http.MethodPost, reactionPath, bob, map[string]string{"emoji": "👍"})
	env.expectStatus(resp, payload, http.StatusCreated)
	var reacted struct {
		Comment database.Comment `json:"comment"`
	}
	if err := json.Unmarshal(payload, &reacted); err != nil {
		t.Fatalf("decode reaction response: %v", err)
	}
	if reacted.Comment.CommentId == 0 || reacted.Comment.Emoji != "👍" {
		t.Fatalf("unexpected reaction: %#v", reacted.Comment)
	}

	resp, payload = env.request(http.MethodPost, reactionPath, bob, map[string]string{"emoji": "✨"})
	env.expectStatus(resp, payload, http.StatusOK)
	if !bytes.Contains(payload, []byte(`"emoji":"✨"`)) {
		t.Fatalf("reaction was not updated: %s", payload)
	}

	resp, payload = env.request(http.MethodGet, fmt.Sprintf("/conversations/%d/messages", secondConversation), alice, nil)
	env.expectStatus(resp, payload, http.StatusOK)
	if !bytes.Contains(payload, []byte(`"text":"hello"`)) {
		t.Fatalf("forwarded message missing from destination: %s", payload)
	}
}

func TestConversationAuthorization(t *testing.T) {
	env := newTestAPI(t)
	alice := env.login("alice")
	bob := env.login("bobby")
	outsider := env.login("carol")
	conversationID := env.createConversation(alice, bob)

	resp, payload := env.request(http.MethodGet, fmt.Sprintf("/conversations/%d/messages", conversationID), outsider, nil)
	env.expectStatus(resp, payload, http.StatusForbidden)

	t.Run("outsider cannot send a message", func(t *testing.T) {
		resp, payload = env.request(http.MethodPost, fmt.Sprintf("/conversations/%d/messages", conversationID), outsider, map[string]string{"type": "text", "text": "forbidden"})
		env.expectStatus(resp, payload, http.StatusForbidden)

		messages, err := env.db.GetConversationMessages(conversationID)
		if err != nil {
			t.Fatalf("read messages after rejected request: %v", err)
		}
		if len(messages) != 0 {
			t.Fatalf("rejected request still created %d message(s)", len(messages))
		}
	})
}

func TestEmptyCollectionsAreJSONArrays(t *testing.T) {
	env := newTestAPI(t)
	alice := env.login("alice")

	resp, payload := env.request(http.MethodGet, "/users", alice, nil)
	env.expectStatus(resp, payload, http.StatusOK)
	if !bytes.Contains(payload, []byte(`"users":[]`)) {
		t.Fatalf("empty users must be encoded as arrays: %s", payload)
	}

	resp, payload = env.request(http.MethodGet, "/conversations", alice, nil)
	env.expectStatus(resp, payload, http.StatusOK)
	if !bytes.Contains(payload, []byte(`"conversations":[]`)) {
		t.Fatalf("empty conversations must be encoded as an array: %s", payload)
	}

	bob := env.login("bobby")
	conversationID := env.createConversation(alice, bob)
	resp, payload = env.request(http.MethodGet, fmt.Sprintf("/conversations/%d/messages", conversationID), alice, nil)
	env.expectStatus(resp, payload, http.StatusOK)
	if !bytes.Contains(payload, []byte(`"messages":[]`)) {
		t.Fatalf("empty messages must be encoded as an array: %s", payload)
	}
}

func TestMessageAndConversationOrdering(t *testing.T) {
	env := newTestAPI(t)
	alice := env.login("alice")
	bob := env.login("bobby")
	carol := env.login("carol")

	firstConversation := env.createConversation(alice, bob)
	secondConversation := env.createConversation(alice, carol)
	resp, payload := env.request(http.MethodGet, "/conversations", alice, nil)
	env.expectStatus(resp, payload, http.StatusOK)
	var conversations struct {
		Items []database.Conversation `json:"conversations"`
	}
	if err := json.Unmarshal(payload, &conversations); err != nil {
		t.Fatalf("decode conversations: %v", err)
	}
	if len(conversations.Items) != 2 || conversations.Items[0].ConversationId != secondConversation || conversations.Items[1].ConversationId != firstConversation {
		t.Fatalf("empty conversations are not ordered by descending id: %#v", conversations.Items)
	}

	for _, text := range []string{"first", "second", "third"} {
		resp, payload = env.request(http.MethodPost, fmt.Sprintf("/conversations/%d/messages", firstConversation), alice, map[string]string{"type": "text", "text": text})
		env.expectStatus(resp, payload, http.StatusCreated)
	}
	resp, payload = env.request(http.MethodGet, fmt.Sprintf("/conversations/%d/messages", firstConversation), alice, nil)
	env.expectStatus(resp, payload, http.StatusOK)
	var messages struct {
		Items []database.Message `json:"messages"`
	}
	if err := json.Unmarshal(payload, &messages); err != nil {
		t.Fatalf("decode messages: %v", err)
	}
	if len(messages.Items) != 3 || messages.Items[0].Text != "first" || messages.Items[1].Text != "second" || messages.Items[2].Text != "third" {
		t.Fatalf("messages are not ordered by ascending id: %#v", messages.Items)
	}
}

func TestCrossConversationMessageReferencesAreRejected(t *testing.T) {
	env := newTestAPI(t)
	alice := env.login("alice")
	bob := env.login("bobby")
	carol := env.login("carol")
	firstConversation := env.createConversation(alice, bob)
	secondConversation := env.createConversation(alice, carol)

	resp, payload := env.request(http.MethodPost, fmt.Sprintf("/conversations/%d/messages", firstConversation), alice, map[string]string{"type": "text", "text": "belongs to first"})
	env.expectStatus(resp, payload, http.StatusCreated)
	var sent struct {
		Message database.Message `json:"message"`
	}
	if err := json.Unmarshal(payload, &sent); err != nil {
		t.Fatalf("decode sent message: %v", err)
	}

	tests := []struct {
		name   string
		method string
		path   string
		body   interface{}
	}{
		{"forward", http.MethodPost, fmt.Sprintf("/conversations/%d/messages/%d/forward/%d", secondConversation, sent.Message.MessageId, firstConversation), nil},
		{"react", http.MethodPost, fmt.Sprintf("/conversations/%d/messages/%d/comments", secondConversation, sent.Message.MessageId), map[string]string{"emoji": "👍"}},
		{"delete", http.MethodDelete, fmt.Sprintf("/conversations/%d/messages/%d", secondConversation, sent.Message.MessageId), nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, payload := env.request(test.method, test.path, alice, test.body)
			env.expectStatus(resp, payload, http.StatusForbidden)
		})
	}

	messages, err := env.db.GetConversationMessages(firstConversation)
	if err != nil {
		t.Fatalf("read source conversation: %v", err)
	}
	if len(messages) != 1 || messages[0].MessageId != sent.Message.MessageId || len(messages[0].Comments) != 0 {
		t.Fatalf("rejected references mutated source message: %#v", messages)
	}
	destinationMessages, err := env.db.GetConversationMessages(secondConversation)
	if err != nil {
		t.Fatalf("read other conversation: %v", err)
	}
	if len(destinationMessages) != 0 {
		t.Fatalf("rejected forward created messages: %#v", destinationMessages)
	}
}

func (env *testAPI) rawRequest(method, path, authorization string) (*http.Response, []byte) {
	env.t.Helper()
	req, err := http.NewRequest(method, env.server.URL+path, nil)
	if err != nil {
		env.t.Fatal(err)
	}
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	resp, err := env.server.Client().Do(req)
	if err != nil {
		env.t.Fatal(err)
	}
	defer resp.Body.Close()
	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		env.t.Fatal(err)
	}
	return resp, payload
}

func TestEveryApplicationRouteRequiresSession(t *testing.T) {
	env := newTestAPI(t)
	alice := env.login("alice")
	paths := []struct{ method, path string }{
		{"POST", "/logout"}, {"PUT", "/settings/username"}, {"PUT", "/settings/photo"}, {"GET", "/settings/photo/1"},
		{"GET", "/users"}, {"POST", "/users/2"}, {"GET", "/conversations"}, {"POST", "/conversations"},
		{"GET", "/conversations/1"}, {"GET", "/conversations/1/messages"}, {"POST", "/conversations/1/messages"},
		{"POST", "/conversations/1/messages/1/forward/2"}, {"POST", "/conversations/1/messages/1/comments"},
		{"DELETE", "/conversations/1/messages/1/comments/1"}, {"DELETE", "/conversations/1/messages/1"},
		{"PUT", "/conversations/1/members"}, {"DELETE", "/conversations/1/members"}, {"PUT", "/conversations/1/name"}, {"PUT", "/conversations/1/photo"},
	}
	for _, endpoint := range paths {
		t.Run(endpoint.method+endpoint.path, func(t *testing.T) {
			for _, auth := range []string{"", fmt.Sprint(alice), fmt.Sprintf("Bearer %d", alice), env.tokens[alice], "Basic " + env.tokens[alice], "Bearer " + string(bytes.Repeat([]byte("x"), 43)), "Bearer " + env.tokens[alice] + " extra"} {
				resp, payload := env.rawRequest(endpoint.method, endpoint.path, auth)
				env.expectStatus(resp, payload, http.StatusUnauthorized)
				if resp.Header.Get("WWW-Authenticate") != "Bearer" {
					t.Fatal("missing Bearer challenge")
				}
			}
		})
	}
}

func TestSessionsAreIndependentAndRevokedAtLogout(t *testing.T) {
	env := newTestAPI(t)
	alice := env.login("alice")
	firstToken := env.tokens[alice]
	env.login("alice")
	secondToken := env.tokens[alice]
	if firstToken == secondToken {
		t.Fatal("repeated login reused token")
	}
	resp, payload := env.rawRequest("GET", "/users", "Bearer "+firstToken)
	env.expectStatus(resp, payload, http.StatusOK)
	resp, payload = env.rawRequest("POST", "/logout", "Bearer "+firstToken)
	env.expectStatus(resp, payload, http.StatusNoContent)
	resp, payload = env.rawRequest("GET", "/users", "Bearer "+firstToken)
	env.expectStatus(resp, payload, http.StatusUnauthorized)
	resp, payload = env.rawRequest("GET", "/users", "Bearer "+secondToken)
	env.expectStatus(resp, payload, http.StatusOK)
	if _, err := env.dbconn.Exec("DELETE FROM users_table WHERE id=?", alice); err != nil {
		t.Fatal(err)
	}
	resp, payload = env.rawRequest("GET", "/users", "Bearer "+secondToken)
	env.expectStatus(resp, payload, http.StatusUnauthorized)
}

func TestCentralizedPermissions(t *testing.T) {
	env := newTestAPI(t)
	alice, bob, outsider := env.login("alice"), env.login("bobby"), env.login("carol")
	chat := env.createConversation(alice, bob)
	otherChat := env.createConversation(bob, outsider)
	resp, payload := env.request("POST", fmt.Sprintf("/conversations/%d/messages", chat), alice, map[string]string{"type": "text", "text": "original"})
	env.expectStatus(resp, payload, http.StatusCreated)
	var sent struct {
		Message database.Message `json:"message"`
	}
	if err := json.Unmarshal(payload, &sent); err != nil {
		t.Fatal(err)
	}
	msg := sent.Message.MessageId
	resp, payload = env.request("POST", fmt.Sprintf("/conversations/%d/messages/%d/comments", chat, msg), alice, map[string]string{"emoji": "👍"})
	env.expectStatus(resp, payload, http.StatusCreated)
	var reacted struct {
		Comment database.Comment `json:"comment"`
	}
	if err := json.Unmarshal(payload, &reacted); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		method, path string
		user, status int
	}{
		{"GET", fmt.Sprintf("/conversations/%d", chat), outsider, 403},
		{"PUT", fmt.Sprintf("/conversations/%d/name", chat), outsider, 403},
		{"PUT", fmt.Sprintf("/conversations/%d/photo", chat), outsider, 403},
		{"PUT", fmt.Sprintf("/conversations/%d/members", chat), outsider, 403},
		{"DELETE", fmt.Sprintf("/conversations/%d/members", chat), outsider, 403},
		{"PUT", fmt.Sprintf("/conversations/%d/name", chat), alice, 403},
		{"DELETE", fmt.Sprintf("/conversations/%d/messages/%d", chat, msg), bob, 403},
		{"DELETE", fmt.Sprintf("/conversations/%d/messages/%d/comments/%d", chat, msg, reacted.Comment.CommentId), bob, 403},
		{"POST", fmt.Sprintf("/conversations/%d/messages/%d/forward/%d", chat, msg, otherChat), alice, 403},
		{"GET", fmt.Sprintf("/settings/photo/%d", bob), alice, 403},
		{"GET", "/conversations/999999", alice, 404}, {"GET", "/conversations/nope", alice, 400},
		{"GET", fmt.Sprintf("/conversations/%d/messages/", chat), alice, 404},
	}
	for _, test := range tests {
		resp, payload = env.request(test.method, test.path, test.user, nil)
		env.expectStatus(resp, payload, test.status)
	}
	resp, payload = env.request("DELETE", fmt.Sprintf("/conversations/%d/messages/%d/comments/%d", chat, msg, reacted.Comment.CommentId), alice, nil)
	env.expectStatus(resp, payload, http.StatusNoContent)
	resp, payload = env.request("DELETE", fmt.Sprintf("/conversations/%d/messages/%d", chat, msg), alice, nil)
	env.expectStatus(resp, payload, http.StatusNoContent)
}

func TestGroupPermissionsAfterMembershipChanges(t *testing.T) {
	env := newTestAPI(t)
	alice, bob, carol, dave := env.login("alice"), env.login("bobby"), env.login("carol"), env.login("david")
	resp, payload := env.request("POST", "/conversations", alice, map[string]interface{}{"title": "Team", "userIds": []int{bob, carol}})
	env.expectStatus(resp, payload, 201)
	var created struct {
		Conversation database.Conversation `json:"conversation"`
	}
	if err := json.Unmarshal(payload, &created); err != nil {
		t.Fatal(err)
	}
	path := fmt.Sprintf("/conversations/%d", created.Conversation.ConversationId)
	resp, payload = env.request("PUT", path+"/members", dave, map[string]int{"userId": dave})
	env.expectStatus(resp, payload, 403)
	resp, payload = env.request("PUT", path+"/members", bob, map[string]int{"userId": dave})
	env.expectStatus(resp, payload, 200)
	resp, payload = env.request("PUT", path+"/name", dave, map[string]string{"title": "Updated"})
	env.expectStatus(resp, payload, 200)
	resp, payload = env.request("PUT", path+"/photo", dave, map[string]interface{}{"image": testImageUpload()})
	env.expectStatus(resp, payload, 200)
	resp, payload = env.request("DELETE", path+"/members", dave, nil)
	env.expectStatus(resp, payload, 204)
	resp, payload = env.request("GET", path, dave, nil)
	env.expectStatus(resp, payload, 403)
	resp, payload = env.request("POST", path+"/messages", dave, map[string]string{"type": "text", "text": "cannot send"})
	env.expectStatus(resp, payload, 403)
	resp, payload = env.request("GET", path, bob, nil)
	env.expectStatus(resp, payload, 200)
}

func TestForwardMetadataIsVisibleAcrossSessionsAndSurvivesSourceDeletion(t *testing.T) {
	env := newTestAPI(t)
	alice, bob, carol := env.login("alice"), env.login("bobby"), env.login("carol")
	source := env.createConversation(alice, bob)
	destination := env.createConversation(alice, carol)
	resp, payload := env.request("POST", fmt.Sprintf("/conversations/%d/messages", source), alice, map[string]string{"type": "text", "text": "persisted forward"})
	env.expectStatus(resp, payload, 201)
	var sent struct {
		Message database.Message `json:"message"`
	}
	if err := json.Unmarshal(payload, &sent); err != nil {
		t.Fatal(err)
	}
	if sent.Message.CreatedAt.IsZero() || sent.Message.Forwarded {
		t.Fatal("new message metadata is invalid")
	}
	resp, payload = env.request("POST", fmt.Sprintf("/conversations/%d/messages/%d/forward/%d", source, sent.Message.MessageId, destination), alice, nil)
	env.expectStatus(resp, payload, 201)
	var forwarded struct {
		Message database.Message `json:"message"`
	}
	if err := json.Unmarshal(payload, &forwarded); err != nil {
		t.Fatal(err)
	}
	if !forwarded.Message.Forwarded || forwarded.Message.OriginalMessageID == nil || *forwarded.Message.OriginalMessageID != sent.Message.MessageId {
		t.Fatalf("forward metadata missing: %s", payload)
	}
	resp, payload = env.request("GET", fmt.Sprintf("/conversations/%d/messages", destination), carol, nil)
	env.expectStatus(resp, payload, 200)
	var read struct {
		Messages []database.Message `json:"messages"`
	}
	if err := json.Unmarshal(payload, &read); err != nil {
		t.Fatal(err)
	}
	if len(read.Messages) != 1 || !read.Messages[0].Forwarded {
		t.Fatal("recipient did not receive stored forward flag")
	}
	resp, payload = env.request("DELETE", fmt.Sprintf("/conversations/%d/messages/%d", source, sent.Message.MessageId), alice, nil)
	env.expectStatus(resp, payload, 204)
	resp, payload = env.request("GET", fmt.Sprintf("/conversations/%d/messages", destination), carol, nil)
	env.expectStatus(resp, payload, 200)
	if err := json.Unmarshal(payload, &read); err != nil {
		t.Fatal(err)
	}
	if len(read.Messages) != 1 || !read.Messages[0].Forwarded || read.Messages[0].OriginalMessageID != nil || read.Messages[0].Text != "persisted forward" {
		t.Fatalf("source deletion damaged forwarded message: %s", payload)
	}
}
