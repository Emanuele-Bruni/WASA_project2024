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

func (rt *_router) newGroup(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	//calcolo l'id utente e vedo se è autorizzato
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	id, err := strconv.Atoi(token)
	if err != nil {
		ctx.Logger.Error("Token must be a number: ", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	//requestBody con il titolo e un array di interi (sarebbero gli id degli utenti da aggiungere)
	var requestBody = struct {
		Title string `json:"title"`
		Users []int  `json:"users"`
	}{}
	//decodifico la richiesta secondo il requestBody
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Struttura non compatibile", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Input non compatibile con il RequestBody")
		return
	}
	//controllo che i valori passati siano validi
	if len(requestBody.Title) < 1 || len(requestBody.Title) > 16 || len(requestBody.Users) < 2 || len(requestBody.Users) > 16 {
		http.Error(w, "Valori non accettabili", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Lunghezza dei valori errata")
		return
	}

	//creo il gruppo e ritorno l'id
	groupId, err := rt.db.CreateGroup(requestBody.Title, requestBody.Users, id)
	if err != nil {
		ctx.Logger.Error("Errore durante la creazione del gruppo: ", err)
		http.Error(w, "Impossible to add yourself into a group", http.StatusInternalServerError)
		return
	}
	//devo ritornare il gruppo -> richiamo la funzione GetConversation
	gruppo, err := rt.db.GetSpecificConversation(groupId)
	if err != nil {
		ctx.Logger.Error("Errore durante la selezione del nuovo gruppo: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	var newGroupResponse = struct {
		GroupId int             `json:"groupId"`
		Title   string          `json:"title"`
		Photo   string          `json:"photo"`
		IsGroup bool            `json:"isGroup"`
		Members []database.User `json:"members"`
	}{
		GroupId: gruppo.ConversationId,
		Title:   gruppo.Title,
		Photo:   gruppo.Photo,
		IsGroup: gruppo.IsGroup,
		Members: gruppo.Members,
	}
	if err := json.NewEncoder(w).Encode(newGroupResponse); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Failed to encode response")
		return
	}

	ctx.Logger.Infof("Gruppo creato: %v", groupId)
}
