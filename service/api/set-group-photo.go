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

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, p httprouter.Params, ctx reqcontext.RequestContext) {
	// calcolo userId e vedo se è autorizzato
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
	// controllo se la conversazione è un gruppo
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
	// requestBody con groupName = il titolo del gruppo
	var requestBody = struct {
		Photo string `json:"photo"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Struttura non compatibile", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Input non compatibile con il RequestBody")
		return
	}
	// controllo se va bene
	if len(requestBody.Photo) == 0 {
		http.Error(w, "Invalid photo", http.StatusBadRequest)
		ctx.Logger.WithField("photo", requestBody.Photo).Error("Invalid length of field 'photo'")
		return
	} else if len(requestBody.Photo) == 0 {
		requestBody.Photo = database.DefaultGroupPhoto
	}

	// aggiorno la foto della conversazione
	err = rt.db.SetGroupPhoto(requestBody.Photo, cnvId)
	if err != nil {
		http.Error(w, "Errore nella modifica della foto", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Errore durante la modifica della foto del gruppo")
		return
	}

	// responso con la nuova foto
	var setGroupPhotoResponse = struct {
		Photo string `json:"photo"`
	}{
		Photo: requestBody.Photo,
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(setGroupPhotoResponse); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Failed to encode response")
		return
	}

	ctx.Logger.Infof("New photo set successfully: %s", setGroupPhotoResponse)
}
