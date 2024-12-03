package api

import (
	"WASA_project2024/service/api/reqcontext"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/*type User struct {
	name   string
	photo  Photo
	bio    string
	number string
}

type Photo struct {
	photoId      int
	photoContent string
	width        int
	height       int
}

type Conversation struct {
	conversatioId int
	title         string
	photo         Photo
	isGroup       bool
}

type Message struct {
	messageId int
	author    string
	text      string
	hours     string
	comments  []Comment
}

type Comment struct {
	commentId int
	author    string
	text      string
}

type Group struct {
	groupId    int
	groupPhoto Photo
	groupName  string
	members    []User
} */

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var requestBody = struct {
		Name string `json:"name"` //"name" nel RequestBody dell'API
	}{}

	//decodifico la richiesta secondo il requestBody
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Struttura non compatibile", http.StatusBadRequest)           //tipo di errore
		ctx.Logger.WithError(err).Error("Input non compatibile con il RequestBody") //descrizione errore
		return
	}

	//controllo se il nome è lungo minimo 3 e massimo 16
	if len(requestBody.Name) < 3 || len(requestBody.Name) > 16 || len(requestBody.Name) == 0 {
		http.Error(w, "Invalid Name", http.StatusBadRequest)
		ctx.Logger.WithField("name", requestBody.Name).Error("Invalid length of field 'name'")
		return
	}

	var loginResponse = struct {
		Identifier int `json:"identifier"`
	}{}

	//controllo se l'utente è nella tabella
	var err error
	loginResponse.Identifier, err = rt.db.GetIdByName(requestBody.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			loginResponse.Identifier, err = rt.db.CreateUser(requestBody.Name)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				ctx.Logger.WithError(err).Error("Error creating new user")
				return
			}
			ctx.Logger.Infof("Created new user: %s", requestBody.Name)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("Error finding id from name")
			return
		}
	}

	//la risposta sarà in formato JSON, impostando il content-type della risposta
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Authentication", fmt.Sprintf("bearerAuth: %d", loginResponse.Identifier))
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(loginResponse); err != nil {
		ctx.Logger.Error("Failed to encode response: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//Logga il successo
	ctx.Logger.Infof("Login succesful: %s", requestBody.Name)
}
