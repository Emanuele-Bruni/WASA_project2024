package api

import (
	"WASA_project2024/service/api/reqcontext"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, p httprouter.Params, ctx reqcontext.RequestContext) {
	// calcolo l'id utente e vedo se è autorizzato
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	id, err := strconv.Atoi(token)
	if err != nil {
		ctx.Logger.Error("Token must be a number: ", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	// controllo se la conversazione esiste
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

	// controllo se il messaggio esiste
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

	// controllo se il messaggio è stato scritto dall'utente
	message_user, err := rt.db.CheckMessage_User(messageId, id)
	if err != nil {
		http.Error(w, "Errore nel controllo dell'autore del messaggio", http.StatusInternalServerError)
		ctx.Logger.Error("Errore nel controllo dell'autore del messaggio")
		return
	}
	if !message_user {
		http.Error(w, "Il messaggio non è dell'utente", http.StatusForbidden)
		ctx.Logger.Error("Un messaggio può essere rimosso solo dall'autore")
		return
	}

	// lo rimuovo dal database tramite una funzione (elimino anche tutti i commenti)
	err = rt.db.RemoveMessage(messageId)
	if err != nil {
		http.Error(w, "Errore durante l'eliminazione del messaggio", http.StatusInternalServerError)
		ctx.Logger.Error("Messaggio non rimosso")
		return
	}

	w.WriteHeader(http.StatusNoContent)
	ctx.Logger.Infof("Messaggio eliminato con successo")
}
