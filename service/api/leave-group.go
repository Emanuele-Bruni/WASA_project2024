package api

import (
	"WASA_project2024/service/api/reqcontext"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, p httprouter.Params, ctx reqcontext.RequestContext) {
	//calcolo l'utente e vedo se è autorizzato
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	id, err := strconv.Atoi(token)
	if err != nil {
		ctx.Logger.Error("Token must be a number: ", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	//controllo se la conversazione esiste
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
	//controllo se la conversazione è un gruppo
	isGroup, err := rt.db.Conversation_isGroup(cnvId)
	if err != nil {
		http.Error(w, "Errore nel controllo del gruppo", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Errore durante la verifica che la conversazione sia un gruppo")
		return
	}
	if !isGroup {
		http.Error(w, "Conversazione non è un gruppo", http.StatusForbidden)
		ctx.Logger.Error("La conversazione indicata non è un gruppo")
		return
	}

	// rimuovere la conversazione dalla tabella utenti
	err = rt.db.RemoveGroupFromMember(id, cnvId)
	if err != nil {
		http.Error(w, "Utente non rimosso", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Errore durante la rimozione del gruppo dalle conversazioni dell'utente")
		return
	}
	//rimuovere l'utente dai membri della conversazione
	err = rt.db.RemoveMemberFromGroup(id, cnvId)
	if err != nil {
		http.Error(w, "Utente non rimosso", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Errore durante la rimozione dell'utente dal gruppo")
		return
	}

	w.WriteHeader(http.StatusNoContent)
	ctx.Logger.Infof("Gruppo lasciato con successo")
}
