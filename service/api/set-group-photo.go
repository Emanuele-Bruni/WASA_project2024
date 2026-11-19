package api

import (
	"WASA_project2024/service/api/reqcontext"
	"WASA_project2024/service/media"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, p httprouter.Params, ctx reqcontext.RequestContext) {
	var request struct {
		Image media.Upload `json:"image"`
	}
	if !readJSON(w, r, &request, false) {
		return
	}
	image, err := media.Decode(request.Image)
	if err != nil {
		writeError(w, 400, err.Error(), "image")
		return
	}
	if err = rt.db.SetGroupPhoto(image, ctx.ConversationID); err != nil {
		respondDatabaseError(w, err, ctx)
		return
	}
	rt.respondConversation(w, 200, ctx.ConversationID, ctx)
}
