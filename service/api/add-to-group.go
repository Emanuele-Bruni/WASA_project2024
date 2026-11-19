package api

import (
	"WASA_project2024/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, p httprouter.Params, ctx reqcontext.RequestContext) {
	var request struct {
		UserID int `json:"userId"`
	}
	if !readJSON(w, r, &request, false) {
		return
	}
	if request.UserID < 1 {
		writeError(w, 400, "ID utente non valido", "userId")
		return
	}
	exists, err := rt.db.CheckUserId(request.UserID)
	if err != nil {
		respondDatabaseError(w, err, ctx)
		return
	}
	if !exists {
		writeError(w, 404, "L’utente non esiste", "userId")
		return
	}
	member, err := rt.db.CheckUser_Conversation(ctx.ConversationID, request.UserID)
	if err != nil {
		respondDatabaseError(w, err, ctx)
		return
	}
	if member {
		writeError(w, 409, "L’utente appartiene già al gruppo", "userId")
		return
	}
	if err = rt.db.AddUserToGroup(request.UserID, ctx.ConversationID); err != nil {
		respondDatabaseError(w, err, ctx)
		return
	}
	rt.respondConversation(w, 200, ctx.ConversationID, ctx)
}
