package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"WASA_project2024/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

func TestExpiredAndRestartedSessionsAreRejected(t *testing.T) {
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	rt := &_router{baseLogger: logger, sessions: make(map[string]session)}
	token, expiresAt, err := rt.createSession(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(token) != 43 || time.Until(expiresAt) < 23*time.Hour {
		t.Fatal("invalid token or expiration")
	}
	if _, found := rt.sessions[token]; found {
		t.Fatal("raw token stored on server")
	}
	key := sessionKey(token)
	rt.sessions[key] = session{userID: 1, expiresAt: time.Now().Add(-time.Second)}
	handler := rt.protected(func(http.ResponseWriter, *http.Request, httprouter.Params, reqcontext.RequestContext) {
		t.Fatal("unauthorized request reached handler")
	}, memberAccess)
	for _, test := range []string{"expired", "restarted"} {
		t.Run(test, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/users", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			recorder := httptest.NewRecorder()
			handler(recorder, req, nil)
			if recorder.Code != 401 {
				t.Fatalf("got %d, want 401", recorder.Code)
			}
			if len(rt.sessions) != 0 {
				t.Fatal("expired session not removed")
			}
		})
	}
}

func TestConcurrentSessionCreationAndRevocation(t *testing.T) {
	rt := &_router{sessions: make(map[string]session)}
	results := make(chan string, 32)
	for i := 0; i < 32; i++ {
		go func() {
			token, _, err := rt.createSession(1)
			if err != nil {
				results <- ""
				return
			}
			rt.revokeSession(sessionKey(token))
			results <- token
		}()
	}
	seen := make(map[string]bool)
	for i := 0; i < 32; i++ {
		token := <-results
		if token == "" || seen[token] {
			t.Fatal("empty or repeated session token")
		}
		seen[token] = true
	}
	if len(rt.sessions) != 0 {
		t.Fatal("sessions were not revoked")
	}
}
