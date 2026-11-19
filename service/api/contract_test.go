package api_test

import (
	"WASA_project2024/service/database"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
	"unicode/utf8"

	"gopkg.in/yaml.v2"
)

type contractSchema struct {
	Ref                  string                     `yaml:"$ref"`
	Type                 string                     `yaml:"type"`
	Nullable             bool                       `yaml:"nullable"`
	Properties           map[string]*contractSchema `yaml:"properties"`
	Required             []string                   `yaml:"required"`
	AdditionalProperties *bool                      `yaml:"additionalProperties"`
	Items                *contractSchema            `yaml:"items"`
	Minimum              *float64                   `yaml:"minimum"`
	MinLength            *int                       `yaml:"minLength"`
	MaxLength            *int                       `yaml:"maxLength"`
	MaxBytes             *int                       `yaml:"x-maxUtf8Bytes"`
	Pattern              string                     `yaml:"pattern"`
	Format               string                     `yaml:"format"`
	Enum                 []interface{}              `yaml:"enum"`
}
type contractResponse struct {
	Ref     string `yaml:"$ref"`
	Content map[string]struct {
		Schema *contractSchema `yaml:"schema"`
	} `yaml:"content"`
}
type contractOperation struct {
	Responses map[string]contractResponse `yaml:"responses"`
}
type contractPath struct {
	Get    *contractOperation `yaml:"get"`
	Post   *contractOperation `yaml:"post"`
	Put    *contractOperation `yaml:"put"`
	Delete *contractOperation `yaml:"delete"`
}
type contractDocument struct {
	Components struct {
		Schemas   map[string]*contractSchema  `yaml:"schemas"`
		Responses map[string]contractResponse `yaml:"responses"`
	} `yaml:"components"`
	Paths map[string]contractPath `yaml:"paths"`
}

var contractOnce sync.Once
var contract contractDocument
var contractError error

func apiContract(t *testing.T) *contractDocument {
	t.Helper()
	contractOnce.Do(func() {
		_, filename, _, _ := runtime.Caller(0)
		data, err := os.ReadFile(filepath.Join(filepath.Dir(filename), "../../doc/api.yaml"))
		if err != nil {
			contractError = err
			return
		}
		contractError = yaml.Unmarshal(data, &contract)
	})
	if contractError != nil {
		t.Fatal(contractError)
	}
	return &contract
}

// This validates the response subset used by our OpenAPI: exact objects, required fields,
// arrays, primitive types, nullability, string bounds, UTF-8 byte bounds and timestamps.
func checkSchema(t *testing.T, doc *contractDocument, schema *contractSchema, value interface{}, path string) {
	t.Helper()
	if schema.Ref != "" {
		resolved := doc.Components.Schemas[strings.TrimPrefix(schema.Ref, "#/components/schemas/")]
		if resolved == nil {
			t.Fatalf("unknown schema %s", schema.Ref)
		}
		checkSchema(t, doc, resolved, value, path)
		return
	}
	if value == nil {
		if !schema.Nullable {
			t.Fatalf("%s unexpectedly null", path)
		}
		return
	}
	switch schema.Type {
	case "object":
		fields, ok := value.(map[string]interface{})
		if !ok {
			t.Fatalf("%s is not an object", path)
		}
		for _, key := range schema.Required {
			if _, ok = fields[key]; !ok {
				t.Fatalf("%s missing %s", path, key)
			}
		}
		for key, entry := range fields {
			property, exists := schema.Properties[key]
			if !exists {
				if schema.AdditionalProperties != nil && !*schema.AdditionalProperties {
					t.Fatalf("%s has unknown field %s", path, key)
				}
				continue
			}
			checkSchema(t, doc, property, entry, path+"."+key)
		}
	case "array":
		entries, ok := value.([]interface{})
		if !ok {
			t.Fatalf("%s is not an array", path)
		}
		for index, entry := range entries {
			checkSchema(t, doc, schema.Items, entry, fmt.Sprintf("%s[%d]", path, index))
		}
	case "string":
		text, ok := value.(string)
		if !ok {
			t.Fatalf("%s is not a string", path)
		}
		length := utf8.RuneCountInString(text)
		if schema.MinLength != nil && length < *schema.MinLength || schema.MaxLength != nil && length > *schema.MaxLength || schema.MaxBytes != nil && len(text) > *schema.MaxBytes {
			t.Fatalf("%s violates string bounds", path)
		}
		if schema.Pattern != "" && !regexp.MustCompile(schema.Pattern).MatchString(text) {
			t.Fatalf("%s violates pattern", path)
		}
		if schema.Format == "date-time" {
			if _, err := time.Parse(time.RFC3339Nano, text); err != nil {
				t.Fatalf("%s invalid date-time: %v", path, err)
			}
		}
	case "integer":
		number, ok := value.(float64)
		if !ok || math.Trunc(number) != number {
			t.Fatalf("%s is not an integer", path)
		}
		if schema.Minimum != nil && number < *schema.Minimum {
			t.Fatalf("%s below minimum", path)
		}
	case "boolean":
		if _, ok := value.(bool); !ok {
			t.Fatalf("%s is not boolean", path)
		}
	default:
		t.Fatalf("unsupported response schema %s at %s", schema.Type, path)
	}
	if len(schema.Enum) > 0 {
		found := false
		for _, entry := range schema.Enum {
			found = found || fmt.Sprint(entry) == fmt.Sprint(value)
		}
		if !found {
			t.Fatalf("%s violates enum", path)
		}
	}
}

func validateContractResponse(t *testing.T, response *http.Response, body []byte) {
	t.Helper()
	doc := apiContract(t)
	var operation *contractOperation
	for template, item := range doc.Paths {
		pattern := "^" + regexp.MustCompile(`\{[^}]+\}`).ReplaceAllString(template, `[^/]+`) + "$"
		if !regexp.MustCompile(pattern).MatchString(response.Request.URL.Path) {
			continue
		}
		switch response.Request.Method {
		case "GET":
			operation = item.Get
		case "POST":
			operation = item.Post
		case "PUT":
			operation = item.Put
		case "DELETE":
			operation = item.Delete
		}
		break
	}
	var declaration contractResponse
	if operation != nil {
		var exists bool
		declaration, exists = operation.Responses[strconv.Itoa(response.StatusCode)]
		if !exists {
			t.Fatalf("OpenAPI does not document %d for %s %s", response.StatusCode, response.Request.Method, response.Request.URL.Path)
		}
	} else {
		if response.StatusCode != 404 && response.StatusCode != 405 {
			t.Fatalf("response for undocumented route: %d", response.StatusCode)
		}
		declaration.Content = map[string]struct {
			Schema *contractSchema `yaml:"schema"`
		}{"application/json": {Schema: &contractSchema{Ref: "#/components/schemas/Error"}}}
	}
	if declaration.Ref != "" {
		declaration = doc.Components.Responses[strings.TrimPrefix(declaration.Ref, "#/components/responses/")]
	}
	if len(declaration.Content) == 0 {
		if len(body) != 0 {
			t.Fatalf("unexpected body for %d", response.StatusCode)
		}
		return
	}
	media, _, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if err != nil {
		t.Fatal(err)
	}
	content, ok := declaration.Content[media]
	if !ok {
		t.Fatalf("undocumented content type %s", media)
	}
	if media == "application/json" {
		var decoded interface{}
		if err = json.Unmarshal(body, &decoded); err != nil {
			t.Fatal(err)
		}
		checkSchema(t, doc, content.Schema, decoded, "response")
	}
	if response.StatusCode == 201 && response.Header.Get("Location") == "" {
		t.Fatal("created resource has no Location")
	}
	if response.StatusCode == 401 && response.Header.Get("WWW-Authenticate") != "Bearer" {
		t.Fatal("missing Bearer challenge")
	}
	if response.StatusCode == 405 && response.Header.Get("Allow") == "" {
		t.Fatal("missing Allow header")
	}
}

func (env *testAPI) rawBody(method, path, token, body, media string, chunked bool) (*http.Response, []byte) {
	env.t.Helper()
	req, err := http.NewRequest(method, env.server.URL+path, strings.NewReader(body))
	if err != nil {
		env.t.Fatal(err)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if media != "" {
		req.Header.Set("Content-Type", media)
	}
	if chunked {
		req.ContentLength = -1
	}
	response, err := env.server.Client().Do(req)
	if err != nil {
		env.t.Fatal(err)
	}
	defer response.Body.Close()
	payload, err := io.ReadAll(response.Body)
	if err != nil {
		env.t.Fatal(err)
	}
	return response, payload
}

func TestStrictJSONAndRequestLimits(t *testing.T) {
	env := newTestAPI(t)
	tests := []struct {
		name, body, media string
		status            int
	}{
		{"unknown field", `{"name":"alice","extra":1}`, "application/json", 400},
		{"wrong case", `{"Name":"alice"}`, "application/json", 400},
		{"duplicate field", `{"name":"alice","name":"bobby"}`, "application/json", 400},
		{"trailing object", `{"name":"alice"}{}`, "application/json", 400},
		{"trailing data", `{"name":"alice"}junk`, "application/json", 400},
		{"null field", `{"name":null}`, "application/json", 400},
		{"null object", `null`, "application/json", 400},
		{"array", `[]`, "application/json", 400},
		{"wrong type", `{"name":123}`, "application/json", 400},
		{"empty object", `{}`, "application/json", 400},
		{"bad syntax", `{"name":`, "application/json", 400},
		{"invalid UTF8", "{\"name\":\"ali" + string([]byte{0xff}) + "\"}", "application/json", 400},
		{"wrong media", `{"name":"alice"}`, "text/plain", 415},
		{"missing media", `{"name":"alice"}`, "", 415},
		{"wrong charset", `{"name":"alice"}`, "application/json; charset=latin1", 415},
		{"valid charset", `{"name":"alice"}`, "application/json; charset=UTF-8", 200},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, body := env.rawBody("POST", "/login", "", test.body, test.media, false)
			env.expectStatus(response, body, test.status)
		})
	}
	const limit = 2 * 1024 * 1024
	valid := `{"name":"alice"}`
	response, body := env.rawBody("POST", "/login", "", valid+strings.Repeat(" ", limit-len(valid)), "application/json", false)
	env.expectStatus(response, body, 200)
	for _, chunked := range []bool{false, true} {
		response, body = env.rawBody("POST", "/login", "", valid+strings.Repeat(" ", limit-len(valid)+1), "application/json", chunked)
		env.expectStatus(response, body, 413)
	}
}

func TestValidationConflictsAndUpdateStatuses(t *testing.T) {
	env := newTestAPI(t)
	alice, bob, carol, dave := env.login("alice"), env.login("bobby"), env.login("carol"), env.login("david")
	chat := env.createConversation(alice, bob)
	groupBody := map[string]interface{}{"title": " Team ", "userIds": []int{bob, carol}}
	response, body := env.request("POST", "/conversations", alice, groupBody)
	env.expectStatus(response, body, 201)
	var group struct {
		Conversation database.Conversation `json:"conversation"`
	}
	if err := json.Unmarshal(body, &group); err != nil {
		t.Fatal(err)
	}
	path := fmt.Sprintf("/conversations/%d", group.Conversation.ConversationId)
	cases := []struct {
		method, path string
		body         interface{}
		status       int
	}{
		{"POST", fmt.Sprintf("/users/%d", alice), nil, 400},
		{"POST", fmt.Sprintf("/users/%d", bob), nil, 409},
		{"POST", "/users/999999", nil, 404},
		{"PUT", "/settings/username", map[string]string{"name": "bobby"}, 409},
		{"PUT", "/settings/username", map[string]string{"name": "alice"}, 200},
		{"PUT", "/settings/username", map[string]string{"name": "   "}, 400},
		{"PUT", "/settings/username", map[string]string{"name": "ali\nce"}, 400},
		{"PUT", "/settings/username", map[string]string{"name": strings.Repeat("é", 16)}, 200},
		{"PUT", "/settings/username", map[string]string{"name": strings.Repeat("é", 17)}, 400},
		{"PUT", path + "/name", map[string]string{"title": "Updated"}, 200},
		{"PUT", path + "/name", map[string]string{"title": "\t"}, 400},
		{"PUT", path + "/photo", map[string]interface{}{"image": testImageUpload()}, 200},
		{"PUT", "/settings/photo", map[string]interface{}{"image": testImageUpload()}, 200},
		{"GET", fmt.Sprintf("/settings/photo/%d", alice), nil, 200},
		{"PUT", path + "/members", map[string]int{"userId": 999999}, 404},
		{"PUT", path + "/members", map[string]int{"userId": -1}, 400},
		{"PUT", path + "/members", map[string]int{"userId": bob}, 409},
		{"PUT", path + "/members", map[string]int{"userId": dave}, 200},
		{"POST", fmt.Sprintf("/conversations/%d/messages", chat), map[string]string{"type": "text", "text": " \n\t "}, 400},
		{"POST", fmt.Sprintf("/conversations/%d/messages", chat), map[string]string{"type": "text", "text": "hello\x00"}, 400},
		{"POST", fmt.Sprintf("/conversations/%d/messages", chat), map[string]string{"type": "text", "text": strings.Repeat("x", 2000001)}, 400},
		{"PUT", "/settings/photo", map[string]string{"photo": strings.Repeat("x", 2000001)}, 400},
		{"POST", "/conversations", map[string]interface{}{"title": "Team", "userIds": []int{bob, bob}}, 400},
		{"POST", "/conversations", map[string]interface{}{"title": "Team", "userIds": []int{alice, bob}}, 400},
		{"POST", "/conversations", map[string]interface{}{"title": "Team", "userIds": []int{bob, 999999}}, 404},
	}
	for _, test := range cases {
		response, body = env.request(test.method, test.path, alice, test.body)
		env.expectStatus(response, body, test.status)
	}
	response, body = env.request("POST", fmt.Sprintf("/conversations/%d/messages", chat), alice, map[string]string{"type": "text", "text": "multiline\nmessage"})
	env.expectStatus(response, body, 201)
	var sent struct {
		Message database.Message `json:"message"`
	}
	if err := json.Unmarshal(body, &sent); err != nil {
		t.Fatal(err)
	}
	reaction := fmt.Sprintf("/conversations/%d/messages/%d/comments", chat, sent.Message.MessageId)
	for i, emoji := range []string{"👍🏽", "👨‍👩‍👧‍👦", "🇮🇹", "1️⃣", "💙"} {
		response, body = env.request("POST", reaction, alice, map[string]string{"emoji": emoji})
		want := 200
		if i == 0 {
			want = 201
		}
		env.expectStatus(response, body, want)
	}
	for _, emoji := range []string{"", "word", "👍oops", "\u200d", strings.Repeat("👍", 17)} {
		response, body = env.request("POST", reaction, alice, map[string]string{"emoji": emoji})
		env.expectStatus(response, body, 400)
	}
}

func TestRoutingErrorsAndInternalErrorsHaveCommonJSON(t *testing.T) {
	env := newTestAPI(t)
	alice := env.login("alice")
	response, body := env.request("GET", "/not-a-route", alice, nil)
	env.expectStatus(response, body, 404)
	response, body = env.request("PATCH", "/users", alice, nil)
	env.expectStatus(response, body, 405)
	if err := env.dbconn.Close(); err != nil {
		t.Fatal(err)
	}
	response, body = env.request("GET", "/users", alice, nil)
	env.expectStatus(response, body, 500)
	if strings.Contains(string(body), "database") || strings.Contains(string(body), "sql") {
		t.Fatal("internal details leaked")
	}
}

func TestOpenAPIDescribesEveryRegisteredRoute(t *testing.T) {
	doc := apiContract(t)
	_, filename, _, _ := runtime.Caller(0)
	source, err := os.ReadFile(filepath.Join(filepath.Dir(filename), "api-handler.go"))
	if err != nil {
		t.Fatal(err)
	}
	routes := regexp.MustCompile(`rt\.router\.(GET|POST|PUT|DELETE)\("([^"]+)"`).FindAllStringSubmatch(string(source), -1)
	count := 0
	for _, route := range routes {
		template := regexp.MustCompile(`:([a-z]+)`).ReplaceAllString(route[2], `{$1}`)
		item, ok := doc.Paths[template]
		if !ok {
			t.Fatalf("undocumented route %s", template)
		}
		var operation *contractOperation
		switch route[1] {
		case "GET":
			operation = item.Get
		case "POST":
			operation = item.Post
		case "PUT":
			operation = item.Put
		case "DELETE":
			operation = item.Delete
		}
		if operation == nil {
			t.Fatalf("undocumented operation %s %s", route[1], template)
		}
		count++
	}
	documented := 0
	for _, item := range doc.Paths {
		for _, operation := range []*contractOperation{item.Get, item.Post, item.Put, item.Delete} {
			if operation != nil {
				documented++
			}
		}
	}
	if count != documented {
		t.Fatalf("OpenAPI has %d operations, router has %d", documented, count)
	}
}
