package api

import (
	"WASA_project2024/service/api/reqcontext"
	"WASA_project2024/service/database"
	"WASA_project2024/service/media"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, p httprouter.Params, ctx reqcontext.RequestContext) {
	var request struct {
		Type  string        `json:"type"`
		Text  string        `json:"text"`
		Image *media.Upload `json:"image"`
	}
	if !readJSON(w, r, &request, false) {
		return
	}
	var message database.Message
	var id int
	var err error
	switch request.Type {
	case "text":
		if request.Image != nil {
			writeError(w, 400, "Un messaggio testuale non può includere un’immagine", "image")
			return
		}
		if !validateMessage(w, request.Text) {
			return
		}
		message, id, err = rt.db.CreateMessage(request.Text, ctx.ConversationID, ctx.UserID)
	case "image":
		if request.Image == nil || request.Text != "" {
			writeError(w, 400, "Un messaggio immagine richiede image senza testo", "image")
			return
		}
		image, decodeError := media.Decode(*request.Image)
		if decodeError != nil {
			writeError(w, 400, decodeError.Error(), "image")
			return
		}
		message, id, err = rt.db.CreateImageMessage(image, ctx.ConversationID, ctx.UserID)
	default:
		writeError(w, 400, "type deve essere text o image", "type")
		return
	}
	if err != nil {
		respondDatabaseError(w, err, ctx)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/conversations/%d/messages/%d", ctx.ConversationID, id))
	writeJSON(w, 201, messageResponse{message}, ctx)
}
