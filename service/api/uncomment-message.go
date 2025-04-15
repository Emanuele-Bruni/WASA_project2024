package api

import (
	"WASA_project2024/service/api/reqcontext"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, p httprouter.Params, ctx reqcontext.RequestContext) {
	//--------- CONTROLLI ---------
	//calcolo l'id utente e vedo se è autorizzato
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	id, err := strconv.Atoi(token)
	if err != nil {
		ctx.Logger.Error("Token must be a number: ", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	//CONTROLLI CONVERSAZIONE
	cnvId, err := strconv.Atoi(p.ByName("conversationid"))
	if err != nil {
		//errore nel prendere il parametro
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
	//controllo se l'utente è nella conversazione
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

	//CONTROLLI MESSAGGIO
	messageId, err := strconv.Atoi(p.ByName("messageid"))
	if err != nil {
		//errore nel prendere il parametro
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
	//vedo se esiste nella conversazione
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

	//CONTROLLI COMMENTO
	commentId, err := strconv.Atoi(p.ByName("commentid"))
	if err != nil {
		//errore nel prendere il parametro
		http.Error(w, "Parametro non valido", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("commentId parametro non valido")
		return
	}
	//vedo se esiste
	exists, err = rt.db.CommentExists(commentId)
	if err != nil {
		http.Error(w, "Errore interno nel controllo del commento", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Errore durante la verifica dell'esistenza del commento")
		return
	}
	if !exists {
		http.Error(w, "Il commento non esiste", http.StatusNotFound)
		ctx.Logger.Error("commento inesistente")
		return
	}
	//vedo se è di quel messaggio
	commentFound, err := rt.db.CheckComment_Message(commentId, messageId)
	if err != nil {
		http.Error(w, "Errore nel controllo del commento", http.StatusInternalServerError)
		ctx.Logger.Error("Errore durante la verifica dell'appartenenza del commento al messaggio")
		return
	}
	if !commentFound {
		http.Error(w, "Commento non appartiene al messaggio", http.StatusForbidden)
		ctx.Logger.Error("Errore perché il commento non appartiene al messaggio")
		return
	}

	//rimuovo il commento dal messaggio tramite la funzione
	err = rt.db.DeleteComment(commentId)
	if err != nil {
		http.Error(w, "Errore durante la rimozione del commento", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Commento non rimosso")
		return
	}

	w.WriteHeader(http.StatusNoContent)
	ctx.Logger.Infof("Commento eliminato con successo")
}
