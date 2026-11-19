package api

import (
	"WASA_project2024/service/api/reqcontext"
	"WASA_project2024/service/database"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, p httprouter.Params, ctx reqcontext.RequestContext) {
	conversations, err := rt.db.GetConversations(ctx.UserID)
	if err != nil {
		respondDatabaseError(w, err, ctx)
		return
	}
	if conversations == nil {
		conversations = []database.Conversation{}
	}
	writeJSON(w, 200, struct {
		Conversations []database.Conversation `json:"conversations"`
	}{conversations}, ctx)
}
