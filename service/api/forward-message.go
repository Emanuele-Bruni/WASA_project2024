package api

import (
	"WASA_project2024/service/api/reqcontext"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, p httprouter.Params, ctx reqcontext.RequestContext) {
	if !noRequestBody(w, r) {
		return
	}
	message, id, err := rt.db.ForwardMessage(ctx.MessageID, ctx.ForwardID, ctx.UserID)
	if err != nil {
		respondDatabaseError(w, err, ctx)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/conversations/%d/messages/%d", ctx.ForwardID, id))
	writeJSON(w, 201, messageResponse{message}, ctx)
}
