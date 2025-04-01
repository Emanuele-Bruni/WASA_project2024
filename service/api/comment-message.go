package api

import (
	"WASA_project2024/service/api/reqcontext"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, p httprouter.Params, ctx reqcontext.RequestContext) {

	//--------- CONTROLLI ---------
	//calcolo l'id utente e vedo se è autorizzato
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "bearerAuth ")
	id, err := strconv.Atoi(token)
	if err != nil {
		ctx.Logger.Error("Token must be a number: ", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	//prendo il conversation id dal path e vedo se esiste
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

	//prendo il message id dal path e vedo se esiste nella conversazione
	messageId, err := strconv.Atoi(p.ByName("messageid"))
	if err != nil {
		//errore nel prendere il parametro
		http.Error(w, "Parametro non valido", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("messageId parametro non valido")
		return
	}
	messageFound, err := rt.db.CheckMessage_Conversation(messageId, cnvId)
	if err != nil {
		http.Error(w, "Errore nel controllo del messaggio", http.StatusInternalServerError)
		ctx.Logger.Error("Errore durante la verifica dell'appartenenza del messaggio alla conversazione")
		return
	}
	if !messageFound {
		http.Error(w, "Messaggio non appartiene alla conversazione", http.StatusForbidden)
		ctx.Logger.Error("Errore perché il messaggio non appartiene alla conversazione")
		return
	}

	//nel request body devo mettere autore e commento
	var requestBody = struct {
		Author string `json:"author"`
		Emoji  string `json:"emoji"`
	}{}

	//decodifico la richiesta secondo il request body
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Struttura non compatibile", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Input non compatibile con il RequestBody")
		return
	}

	//controllo la lunghezza dei campi
	if len(requestBody.Author) < 1 || len(requestBody.Author) > 16 || len(requestBody.Emoji) < 1 || len(requestBody.Emoji) > 2 {
		http.Error(w, "Invalid lenght of fields", http.StatusBadRequest)
		ctx.Logger.Error("Errore nella lunghezza dei campi")
		return
	}

	//funzione che crea commento e lo aggiunge ai messaggi (append ai commenti)
	commentResponse, err := rt.db.CreateComment(requestBody.Author, requestBody.Emoji, messageId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Errore durante la creazione del commento")
		return
	}

	//preparo l'header
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	//ritorno il risultato
	if err := json.NewEncoder(w).Encode(commentResponse); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Failed to encode response")
		return
	}
	//Logga il successo
	ctx.Logger.Infof("Commento inviato con successo: %s", requestBody.Emoji)
}
