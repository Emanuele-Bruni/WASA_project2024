package api

import (
	"WASA_project2024/service/api/reqcontext"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, p httprouter.Params, ctx reqcontext.RequestContext) {
	//calcolo id e vedo se autorizzato
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

	//requestBody con groupName = il titolo del gruppo
	var requestBody = struct {
		GroupName string `json:"groupName"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Struttura non compatibile", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Input non compatibile con il RequestBody")
		return
	}
	//controllo che sia >1 e <16
	if len(requestBody.GroupName) < 1 || len(requestBody.GroupName) > 16 {
		http.Error(w, "Invalid Group Name", http.StatusBadRequest)
		ctx.Logger.WithField("groupName", requestBody.GroupName).Error("Invalid length of field 'groupName'")
		return
	}

	//aggiorno con una funzione -> SetGroupName(groupName, cnvId)
	err = rt.db.SetGroupName(requestBody.GroupName, cnvId)
	if err != nil {
		http.Error(w, "Nome gruppo non modificato", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Errore durante la modifica del nome del gruppo")
		return
	}

	//responso con il titolo
	var setGroupNameResponse = struct {
		GroupName string `json:"groupName"`
	}{
		GroupName: requestBody.GroupName,
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(setGroupNameResponse); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Failed to encode response")
		return
	}

	ctx.Logger.WithField("id", id).Infof("New group name setted succesful: %s", setGroupNameResponse)
}
