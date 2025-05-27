package api

import (
	"WASA_project2024/service/api/reqcontext"
	"WASA_project2024/service/database"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	//non ho nessun requestBody
	//Calcolo l'id del nome utente da modificare nella tabella
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	id, err := strconv.Atoi(token)
	if err != nil {
		ctx.Logger.Error("Token must be a number: ", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	allMyConversations, err := rt.db.GetConversations(id)
	if err != nil {
		//controllo se l'errore è dovuto alla mancanza di conversazioni
		if errors.Is(err, database.ErrNoConversation) {
			http.Error(w, "L'utente non ha conversazioni", http.StatusNotFound)
			ctx.Logger.WithError(err).Error("Error finding conversations")
			return
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("Error finding id from name")
			return
		}
	}

	var responseConversations = struct {
		AllMyConversations []database.Conversation `json:"conversations"`
	}{
		AllMyConversations: allMyConversations,
	}

	//ritorno le conversazioni dall'id utente tramite la funzione che creo nel database
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	//json.NewEncoder() serve solo quando devo restituire un contenuto alla risposta (es con metodo GET)
	if err := json.NewEncoder(w).Encode(responseConversations); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Failed to encode response")
		return
	}

	//successo
	ctx.Logger.Infof("Conversazioni ricevute")
}
