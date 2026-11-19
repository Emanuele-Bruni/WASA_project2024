package api

import (
	"WASA_project2024/service/api/reqcontext"
	"WASA_project2024/service/database"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) getUsers(w http.ResponseWriter, r *http.Request, p httprouter.Params, ctx reqcontext.RequestContext) {
	users, err := rt.db.GetUsers(ctx.UserID)
	if err != nil {
		respondDatabaseError(w, err, ctx)
		return
	}
	if users == nil {
		users = []database.User{}
	}
	writeJSON(w, 200, struct {
		Users []database.User `json:"users"`
	}{users}, ctx)
}
