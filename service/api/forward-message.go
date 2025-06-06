package api

import (
	"WASA_project2024/service/api/reqcontext"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, p httprouter.Params, ctx reqcontext.RequestContext) {
	// calcolo l'id utente e vedo se è autorizzato
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	id, err := strconv.Atoi(token)
	if err != nil {
		ctx.Logger.Error("Token must be a number: ", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	// CONTROLLI CONVERSAZIONE
	cnvId, err := strconv.Atoi(p.ByName("conversationid"))
	if err != nil {
		// errore nel prendere il parametro
		http.Error(w, "Parametro non valido", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("conversationId parametro non valido")
		return
	}
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
	// controllo se l'utente è nella conversazione
	found, err := rt.db.CheckUser_Conversation(cnvId, id)
	if err != nil {
		http.Error(w, "Errore nel controllo del membro", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Errore durante la verifica dell'appartenenza dell'utente alla conversazione")
		return
	}
	if !found {
		http.Error(w, "Utente non appartiene alla conversazione", http.StatusForbidden)
		ctx.Logger.Error("Errore perché l'utente non appartiene alla conversazione")
		return
	}

	// CONTROLLI MESSAGGIO
	messageId, err := strconv.Atoi(p.ByName("messageid"))
	if err != nil {
		// errore nel prendere il parametro
		http.Error(w, "Parametro non valido", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("messageId parametro non valido")
		return
	}
	exists, err = rt.db.MessageExists(messageId)
	if err != nil {
		http.Error(w, "Errore interno nel controllo del messaggio", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Errore durante la verifica dell'esistenza del messaggio")
		return
	}
	if !exists {
		http.Error(w, "Il messaggio non esiste", http.StatusNotFound)
		ctx.Logger.Error("messaggio inesistente")
		return
	}
	// vedo se esiste nella conversazione
	messageFound, err := rt.db.CheckMessage_Conversation(messageId, cnvId)
	if err != nil {
		http.Error(w, "Errore nel controllo del messaggio", http.StatusInternalServerError)
		ctx.Logger.Error("Errore durante la verifica dell'appartenenza del messaggio alla conversazione")
		return
	}
	if !messageFound {
		http.Error(w, "Messaggio non appartiene alla conversazione", http.StatusForbidden)
		ctx.Logger.Error("Errore perché il messaggio non appartiene alla conversazione")
		return
	}

	// prendo forwardId, che sarebbe l'id della conversazione a cui inoltrare il messaggio e controllo se esiste
	forwardId, err := strconv.Atoi(p.ByName("forwardid"))
	if err != nil {
		// errore nel prendere il parametro
		http.Error(w, "Parametro non valido", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("forwardId parametro non valido")
		return
	}
	exists, err = rt.db.ConversationExists(forwardId)
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
	// controllo se l'utente è nella conversazione a cui inoltro il messaggio
	found, err = rt.db.CheckUser_Conversation(forwardId, id)
	if err != nil {
		http.Error(w, "Errore nel controllo del membro", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Errore durante la verifica dell'appartenenza dell'utente alla conversazione")
		return
	}
	if !found {
		http.Error(w, "Utente non appartiene alla conversazione", http.StatusForbidden)
		ctx.Logger.Error("Errore perché l'utente non appartiene alla conversazione")
		return
	}
	// devo creare un nuovo messaggio, quindi con un nuovo id -> le uniche cose che rimangono uguali sono il contenuto e l'autore
	// estrapolo le informazioni che mi servono del messaggio
	message, err := rt.db.GetMessage(messageId)
	if err != nil {
		http.Error(w, "Errore durante l'acquisizione del messaggio", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Errore durante l'acquisizione del messaggio")
		return
	}
	// calcolo il nuovo orario
	currentTime := time.Now().Format("15:04")
	// devo richiamare la funzione che crea il messaggio
	forwardResponse, fwdId, err := rt.db.CreateMessage(message.Author, message.Text, currentTime, forwardId, id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Errore durante la creazione del messaggio")
		return
	}

	// preparo l'header
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	// ritorno il risultato
	if err := json.NewEncoder(w).Encode(forwardResponse); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Failed to encode response")
		return
	}
	// Logga il successo
	ctx.Logger.Infof("Messaggio inviato con successo: %s", fwdId)

}
