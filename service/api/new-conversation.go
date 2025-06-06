package api

import (
	"WASA_project2024/service/api/reqcontext"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

func (rt *_router) newConversation(w http.ResponseWriter, r *http.Request, p httprouter.Params, ctx reqcontext.RequestContext) {
	// calcolo l'id utente e vedo se è autorizzato
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	userId1, err := strconv.Atoi(token)
	if err != nil {
		ctx.Logger.Error("Token must be a number: ", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	// prendo userId dal path
	userId2, err := strconv.Atoi(p.ByName("userid"))
	if err != nil {
		// errore nel prendere il parametro
		http.Error(w, "Parametro non valido", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("userId parametro non valido")
		return
	}
	// controllo se l'userId esiste
	exists, err := rt.db.CheckUserId(userId2)
	if err != nil {
		http.Error(w, "Errore nel controllo dell'utente con cui creare una conversazione", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Errore durante la verifica che l'utente destinatario esista")
		return
	}
	if !exists {
		http.Error(w, "L'utente non esiste", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("L'userId è di un utente che non esiste")
		return
	}
	// controllo se la conversazione esiste già
	// funzione che controlla se una conversazione tra i due utenti esiste già
	exists, cnvId, err := rt.db.Conversation2MembersExists(userId1, userId2)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Errore durante il controllo dell'esistenza di una conversazione tra i due utenti")
		return
	}
	if exists {
		http.Error(w, "Conversazione già esistente", http.StatusBadRequest)
		ctx.Logger.WithField("conversationId", cnvId).Error("Conversazione già esistente")
		return
	} else {
		// creo la conversazione
		cnvId, err = rt.db.CreateConversation(userId1, userId2)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("Errore durante la creazione della conversazione")
			return
		}
	}

	// ritorno il conversationId
	var newConversationResponse = struct {
		ConversationId int `json:"conversationId"`
	}{
		ConversationId: cnvId,
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(newConversationResponse); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Failed to encode response")
		return
	}

	ctx.Logger.WithFields(logrus.Fields{"utente1": userId1, "utente2": userId2}).Infof("Conversazione creata con ID: %v", cnvId)
}
