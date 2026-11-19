package api

import (
	"WASA_project2024/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, p httprouter.Params, ctx reqcontext.RequestContext) {
	var request struct {
		Title string `json:"title"`
	}
	if !readJSON(w, r, &request, false) || !validateName(w, &request.Title, "title", 1, 25) {
		return
	}
	if err := rt.db.SetGroupName(request.Title, ctx.ConversationID); err != nil {
		respondDatabaseError(w, err, ctx)
		return
	}
	rt.respondConversation(w, 200, ctx.ConversationID, ctx)
}
