package api

import (
	"WASA_project2024/service/api/reqcontext"
	"WASA_project2024/service/database"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, p httprouter.Params, ctx reqcontext.RequestContext) {
	var request struct {
		Name string `json:"name"`
	}
	if !readJSON(w, r, &request, false) || !validateName(w, &request.Name, "name", 3, 16) {
		return
	}
	id, err := rt.db.GetOrCreateUser(request.Name)
	if err != nil {
		respondDatabaseError(w, err, ctx)
		return
	}
	user, err := rt.db.GetUser(id)
	if err != nil {
		respondDatabaseError(w, err, ctx)
		return
	}
	token, expires, err := rt.createSession(id)
	if err != nil {
		ctx.Logger.WithError(err).Error("Create session")
		writeError(w, 500, "Errore interno del server")
		return
	}
	writeJSON(w, 200, struct {
		User      database.User `json:"user"`
		Token     string        `json:"token"`
		ExpiresAt time.Time     `json:"expiresAt"`
	}{user, token, expires}, ctx)
}
