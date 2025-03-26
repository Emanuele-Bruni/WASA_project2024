package api

import (
	"WASA_project2024/service/api/reqcontext"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, p httprouter.Params, ctx reqcontext.RequestContext) {
	//non ho nessun request body
	//prendo il conversationid
	cnvId, err := strconv.Atoi(p.ByName("conversationId"))
	if err != nil {
		//errore nel prendere il parametro
		http.Error(w, "Parametro non valido", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("conversationId invalid parameter")
		return
	}
	// prendo la conversazione dalla tabella con una funzione
	conversation, err := rt.db.GetSpecificConversation(cnvId)
	if err != nil {
		//il conversationId non esiste
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Not found", http.StatusNotFound)
			ctx.Logger.WithError(err).Error("Conversation not found")
			return
		}
		//errore interno al sistema
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Failed to find conversation")
		return
	}

	// la restituisco
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	//json.NewEncoder() serve solo quando devo restituire un contenuto alla risposta (es con metodo GET)
	if err := json.NewEncoder(w).Encode(conversation); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Failed to encode response")
		return
	}
	ctx.Logger.Infof("Conversazione ricevuta")

}
