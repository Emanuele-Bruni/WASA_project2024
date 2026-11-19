package api

import (
	"WASA_project2024/service/api/reqcontext"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) newGroup(w http.ResponseWriter, r *http.Request, p httprouter.Params, ctx reqcontext.RequestContext) {
	var request struct {
		Title   string `json:"title"`
		UserIDs []int  `json:"userIds"`
	}
	if !readJSON(w, r, &request, false) || !validateName(w, &request.Title, "title", 1, 25) {
		return
	}
	if len(request.UserIDs) < 2 || len(request.UserIDs) > 16 {
		writeError(w, 400, "Seleziona da 2 a 16 altri partecipanti", "userIds")
		return
	}
	seen := map[int]bool{ctx.UserID: true}
	for _, id := range request.UserIDs {
		if id < 1 || seen[id] {
			writeError(w, 400, "Partecipanti non validi, ripetuti o comprendenti l’utente corrente", "userIds")
			return
		}
		seen[id] = true
		exists, err := rt.db.CheckUserId(id)
		if err != nil {
			respondDatabaseError(w, err, ctx)
			return
		}
		if !exists {
			writeError(w, 404, "Un partecipante non esiste", "userIds")
			return
		}
	}
	id, err := rt.db.CreateGroup(request.Title, request.UserIDs, ctx.UserID)
	if err != nil {
		respondDatabaseError(w, err, ctx)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/conversations/%d", id))
	rt.respondConversation(w, 201, id, ctx)
}
