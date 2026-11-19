package api

import (
	"WASA_project2024/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, p httprouter.Params, ctx reqcontext.RequestContext) {
	if !noRequestBody(w, r) {
		return
	}
	if err := rt.db.RemoveMessage(ctx.MessageID); err != nil {
		respondDatabaseError(w, err, ctx)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
