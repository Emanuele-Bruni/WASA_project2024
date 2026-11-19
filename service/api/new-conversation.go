package api

import (
	"WASA_project2024/service/api/reqcontext"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) newConversation(w http.ResponseWriter, r *http.Request, p httprouter.Params, ctx reqcontext.RequestContext) {
	if !noRequestBody(w, r) {
		return
	}
	if ctx.UserID == ctx.TargetUserID {
		writeError(w, 400, "Non puoi creare una chat privata con te stesso", "userid")
		return
	}
	exists, _, err := rt.db.Conversation2MembersExists(ctx.UserID, ctx.TargetUserID)
	if err != nil {
		respondDatabaseError(w, err, ctx)
		return
	}
	if exists {
		writeError(w, 409, "La conversazione esiste già")
		return
	}
	id, err := rt.db.CreateConversation(ctx.UserID, ctx.TargetUserID)
	if err != nil {
		respondDatabaseError(w, err, ctx)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/conversations/%d", id))
	rt.respondConversation(w, 201, id, ctx)
}
