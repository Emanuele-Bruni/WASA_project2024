package api

import (
	"WASA_project2024/service/api/reqcontext"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

const sessionLifetime = 24 * time.Hour

type session struct {
	userID    int
	expiresAt time.Time
}

func sessionKey(token string) string {
	digest := sha256.Sum256([]byte(token))
	return hex.EncodeToString(digest[:])
}

func (rt *_router) createSession(userID int) (string, time.Time, error) {
	var secret [32]byte
	if _, err := rand.Read(secret[:]); err != nil {
		return "", time.Time{}, err
	}
	token := base64.RawURLEncoding.EncodeToString(secret[:])
	now := time.Now()
	expiresAt := now.Add(sessionLifetime)
	rt.sessionsMu.Lock()
	defer rt.sessionsMu.Unlock()
	for key, value := range rt.sessions {
		if !now.Before(value.expiresAt) {
			delete(rt.sessions, key)
		}
	}
	rt.sessions[sessionKey(token)] = session{userID: userID, expiresAt: expiresAt}
	return token, expiresAt, nil
}

func (rt *_router) revokeSession(key string) {
	rt.sessionsMu.Lock()
	defer rt.sessionsMu.Unlock()
	delete(rt.sessions, key)
}

func rejectAuthentication(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("WWW-Authenticate", "Bearer")
	writeError(w, http.StatusUnauthorized, "Sessione assente, scaduta o non valida")
}

// protected authenticates every application route before decoding bodies or checking permissions.
func (rt *_router) protected(fn httpRouterHandler, permission accessPolicy) httprouter.Handle {
	return rt.wrap(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
		parts := strings.Split(r.Header.Get("Authorization"), " ")
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || len(parts[1]) != 43 {
			rejectAuthentication(w)
			return
		}
		key := sessionKey(parts[1])
		rt.sessionsMu.Lock()
		identity, found := rt.sessions[key]
		if found && !time.Now().Before(identity.expiresAt) {
			delete(rt.sessions, key)
			found = false
		}
		rt.sessionsMu.Unlock()
		if !found {
			rejectAuthentication(w)
			return
		}
		name, err := rt.db.GetNameById(identity.userID)
		if errors.Is(err, sql.ErrNoRows) {
			rt.revokeSession(key)
			rejectAuthentication(w)
			return
		}
		if err != nil {
			ctx.Logger.WithError(err).Error("Failed to load authenticated user")
			writeError(w, http.StatusInternalServerError, "Errore interno del server")
			return
		}
		ctx.UserID, ctx.UserName, ctx.SessionKey = identity.userID, name, key
		w.Header().Set("Cache-Control", "no-store")
		if !rt.authorize(w, ps, &ctx, permission) {
			return
		}
		fn(w, r, ps, ctx)
	})
}

func (rt *_router) logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	if !noRequestBody(w, r) {
		return
	}
	rt.revokeSession(ctx.SessionKey)
	w.WriteHeader(http.StatusNoContent)
}
