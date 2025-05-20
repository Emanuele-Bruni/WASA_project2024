package api

import (
	"WASA_project2024/service/api/reqcontext"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	//calcolo l'id utente e vedo se è autorizzato
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	id, err := strconv.Atoi(token)
	if err != nil {
		ctx.Logger.Error("Token must be a number: ", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	//creo funzione che ritorna tutti gli utenti della tabella -> GetUsers
	ids, names, err := rt.db.GetAllUsers(id)
	if err != nil {
		ctx.Logger.Error("Errore durante la selezione degli utenti: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//devo restituire ids,names
	var userResponse = struct {
		Ids   []int
		Names []string
	}{
		Ids:   ids,
		Names: names,
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(userResponse); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Failed to encode response")
		return
	}

	ctx.Logger.Infof("Lista di utenti ricevuta: IDs: %v, Names: %v", ids, names)
}
