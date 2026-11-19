package api

import (
	"WASA_project2024/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, p httprouter.Params, ctx reqcontext.RequestContext) {
	var request struct {
		Name string `json:"name"`
	}
	if !readJSON(w, r, &request, false) || !validateName(w, &request.Name, "name", 3, 16) {
		return
	}
	if err := rt.db.SetName(request.Name, ctx.UserID); err != nil {
		respondDatabaseError(w, err, ctx)
		return
	}
	rt.respondUser(w, ctx)
}
