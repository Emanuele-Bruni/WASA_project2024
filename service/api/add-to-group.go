package api

import (
	"WASA_project2024/service/api/reqcontext"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, p httprouter.Params, ctx reqcontext.RequestContext) {
	//calcolo l'id utente e vedo se è autorizzato
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "bearerAuth ")
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
	//requestBody con l'userId da aggiungere
	var requestBody = struct {
		UserId int `json:"user"`
	}{}
	//decodifico la richiesta secondo il requestBody
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Struttura non compatibile", http.StatusBadRequest)           //tipo di errore
		ctx.Logger.WithError(err).Error("Input non compatibile con il RequestBody") //descrizione errore
		return
	}
	//controllo se l'userId esiste
	exists, err = rt.db.CheckUserId(requestBody.UserId)
	if err != nil {
		http.Error(w, "Errore nel controllo dell'utente da aggiungere", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Errore durante la verifica che l'utente da aggiungere esista")
		return
	}
	if !exists {
		http.Error(w, "L'utente da aggiungere non esiste", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("L'userId è di un utente che non esiste")
		return
	}
	//aggiungo l'utente al gruppo
	err = rt.db.AddUserToGroup(requestBody.UserId, cnvId)
	if err != nil {
		http.Error(w, "Utente non aggiunto", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Errore durante l'aggiunta dell'utente al gruppo")
		return
	}

	//devo restituire una lista di utenti -> faccio una funzione GetConversationMembers che ritorna una lista di id dei membri
	membri, err := rt.db.GetConversationMembers(cnvId)
	if err != nil {
		http.Error(w, "Membri non calcolati", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Errore durante la lettura dei membri del gruppo")
		return
	}
	// dalla lista di id, passo alla lista di utenti grazie alla funzione GetUsersMembers
	utenti, err := rt.db.GetUsersMembers(membri)
	if err != nil {
		http.Error(w, "Errore nel selezionare gli utenti", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Errore durante la selezione dei membri del gruppo")
		return
	}

	//responso finale
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(utenti); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Failed to encode response")
		return
	}

	//Log del successo
	ctx.Logger.Infof("Membro aggiunto correttamente: %s", requestBody.UserId)

}
