package api

import (
	"WASA_project2024/service/api/reqcontext"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, p httprouter.Params, ctx reqcontext.RequestContext) {
	var request struct {
		Emoji string `json:"emoji"`
	}
	if !readJSON(w, r, &request, false) || !validateEmoji(w, request.Emoji) {
		return
	}
	comment, created, err := rt.db.SetReaction(ctx.UserID, request.Emoji, ctx.MessageID)
	if err != nil {
		respondDatabaseError(w, err, ctx)
		return
	}
	status := 200
	if created {
		status = 201
		w.Header().Set("Location", fmt.Sprintf("/conversations/%d/messages/%d/comments/%d", ctx.ConversationID, ctx.MessageID, comment.CommentId))
	}
	writeJSON(w, status, commentResponse{comment}, ctx)
}
