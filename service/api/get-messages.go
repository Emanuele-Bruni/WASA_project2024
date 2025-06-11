package api

import (
	"WASA_project2024/service/api/reqcontext"
	"WASA_project2024/service/database"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getMessages(w http.ResponseWriter, r *http.Request, p httprouter.Params, ctx reqcontext.RequestContext) {
	// calcolo l'id utente e vedo se è autorizzato
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	id, err := strconv.Atoi(token)
	if err != nil {
		ctx.Logger.Error("Token must be a number: ", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	// prendo il conversationid
	cnvId, err := strconv.Atoi(p.ByName("conversationid"))
	if err != nil {
		// errore nel prendere il parametro
		http.Error(w, "Parametro non valido", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("conversationId invalid parameter")
		return
	}
	// controllo se l'utente fa parte della conversazione
	// funzione che controlla se l'utente è nella conversazione
	found, err := rt.db.CheckUser_Conversation(cnvId, id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Errore durante la verifica dell'appartenenza dell'utente alla conversazione")
		return
	}
	if !found {
		http.Error(w, "Utente non appartiene alla conversazione", http.StatusForbidden)
		ctx.Logger.Error("Errore perché l'utente non appartiene alla conversazione")
		return
	}
	// Verifica se la conversazione esiste
	exists, err := rt.db.ConversationExists(cnvId)
	if err != nil {
		http.Error(w, "Errore interno nel controllo della conversazione", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Errore durante la verifica dell'esistenza della conversazione")
		return
	}
	if !exists {
		http.Error(w, "La conversazione non esiste", http.StatusNotFound)
		ctx.Logger.Error("conversazione inesistente")
		return
	}

	// funzione che ritorna tutti gli ids dei messaggi della conversazione
	messages, err := rt.db.GetConversationMessages(cnvId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Errore durante la lettura dei messaggi")
		return
	}

	// responso
	var messagesResponse = struct {
		Messages []database.Message
	}{
		Messages: messages,
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(messagesResponse); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Failed to encode response")
		return
	}

	ctx.Logger.Infof("Lista di messaggi ricevuta")
}
