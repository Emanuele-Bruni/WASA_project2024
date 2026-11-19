package api_test

import (
	"net/http"
	"testing"
)

func TestLivenessChecksDatabaseWithoutSession(t *testing.T) {
	env := newTestAPI(t)
	response, body := env.request(http.MethodGet, "/liveness", 0, nil)
	if response.StatusCode != http.StatusOK || len(body) != 0 || response.Header.Get("Cache-Control") != "no-store" {
		t.Fatalf("healthy database: %d %q", response.StatusCode, body)
	}
	if err := env.dbconn.Close(); err != nil {
		t.Fatal(err)
	}
	response, body = env.request(http.MethodGet, "/liveness", 0, nil)
	if response.StatusCode != http.StatusServiceUnavailable || len(body) != 0 {
		t.Fatalf("unavailable database: %d %q", response.StatusCode, body)
	}
}
