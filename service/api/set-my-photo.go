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

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	//calcolo userId e vedo se è autorizzato
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	id, err := strconv.Atoi(token)
	if err != nil {
		ctx.Logger.Error("Token must be a number: ", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	//controllo che esista
	idFound, err := rt.db.CheckUserId(id)
	if err != nil {
		ctx.Logger.Error("Errore nel controllare l'id utente: ", err)
		http.Error(w, "Error while checking userId", http.StatusInternalServerError)
		return
	}
	if !idFound {
		http.Error(w, "L'utente non esiste", http.StatusUnauthorized)
		ctx.Logger.Error("Tentativo di modificare la foto di un utente che non esiste")
		return
	}

	//requestBody con la foto nuova
	var requestBody = struct {
		Photo string `json:"photo"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Struttura non compatibile", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Input non compatibile con il RequestBody")
		return
	}
	//controllo se va bene (se è lunga <=1000000, può essere anche vuota)
	if len(requestBody.Photo) <= 1000000 {
		http.Error(w, "Invalid photo", http.StatusBadRequest)
		ctx.Logger.WithField("photo", requestBody.Photo).Error("Invalid length of field 'photo'")
		return
	} else {
		if len(requestBody.Photo) == 0 {
			requestBody.Photo = database.DefaultContactPhoto
		}
	}

	//aggiorno la foto
	err = rt.db.SetMyPhoto(requestBody.Photo, id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Errore durante l'aggiornamento della foto")
		return
	}

	//responso con la nuova foto
	var setMyPhotoResponse = struct {
		Photo string `json:"photo"`
	}{
		Photo: requestBody.Photo,
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(setMyPhotoResponse); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Failed to encode response")
		return
	}

	ctx.Logger.Infof("New photo set successfully: %s", setMyPhotoResponse)
}
