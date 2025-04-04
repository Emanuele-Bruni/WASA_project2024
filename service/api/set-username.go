package api

import (
	"WASA_project2024/service/api/reqcontext"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// Set your username
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {

	var requestBody = struct {
		Username string `json:"username"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Struttura non compatibile", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Input non compatibile con il RequestBody")
		return
	}

	//Calcolo l'id del nome utente da modificare nella tabella
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	id, err := strconv.Atoi(token)
	if err != nil {
		ctx.Logger.Error("Token must be a number: ", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	//controllo se l'id esiste
	idFound, err := rt.db.CheckUserId(id)
	if err != nil {
		ctx.Logger.Error("Errore nel controllare l'id utente: ", err)
		http.Error(w, "Error while checking userId", http.StatusInternalServerError)
		return
	}
	if !idFound {
		http.Error(w, "L'utente non esiste", http.StatusUnauthorized)
		ctx.Logger.Error("Tentativo di modificare l'username di un utente che non esiste")
		return
	}

	//controllo se l'username è lungo minimo 3 e massimo 16 (perché deve rispettare il login)
	if len(requestBody.Username) < 3 || len(requestBody.Username) > 16 {
		http.Error(w, "Invalid Username", http.StatusBadRequest)
		ctx.Logger.WithField("username", requestBody.Username).Error("Invalid length of field 'username'")
		return
	}
	//username valido, lo aggiorno tramite la funzione
	err = rt.db.SetName(requestBody.Username, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.Logger.WithError(err).Error("User not found")
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		ctx.Logger.Error("Failed to set username: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//Responso finale
	ctx.Logger.WithField("id", id).Infof("New username setted succesful: %s", requestBody.Username)

	w.WriteHeader(http.StatusNoContent) // Imposta lo stato HTTP 204
}
