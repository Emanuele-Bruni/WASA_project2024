package api

import (
	"WASA_project2024/service/api/reqcontext"
	"WASA_project2024/service/database"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, p httprouter.Params, ctx reqcontext.RequestContext) {

	// --------- CONTROLLI ---------
	// calcolo l'id utente e vedo se è autorizzato
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	id, err := strconv.Atoi(token)
	if err != nil {
		ctx.Logger.Error("Token must be a number: ", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// --------------------------------------------------------------------------CONTROLLI CONVERSAZIONE
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

	// --------------------------------------------------------------------------CONTROLLI MESSAGGIO
	messageId, err := strconv.Atoi(p.ByName("messageid"))
	if err != nil {
		// errore nel prendere il parametro
		http.Error(w, "Parametro non valido", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("messageId parametro non valido")
		return
	}
	exists, err = rt.db.MessageExists(messageId)
	if err != nil {
		http.Error(w, "Errore interno nel controllo del messaggio", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Errore durante la verifica dell'esistenza del messaggio")
		return
	}
	if !exists {
		http.Error(w, "Il messaggio non esiste", http.StatusNotFound)
		ctx.Logger.Error("messaggio inesistente")
		return
	}
	// vedo se esiste nella conversazione
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

	// nel request body devo mettere autore e commento
	var requestBody = struct {
		Author string `json:"author"`
		Emoji  string `json:"emoji"`
	}{}

	// decodifico la richiesta secondo il request body
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Struttura non compatibile", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Input non compatibile con il RequestBody")
		return
	}
	// controllo che il commento o l'autore non siano vuoti
	if strings.TrimSpace(requestBody.Author) == "" || strings.TrimSpace(requestBody.Emoji) == "" {
		http.Error(w, "Autore ed emoji non possono essere vuoti", http.StatusBadRequest)
		ctx.Logger.Error("Autore ed emoji non possono essere vuoti")
		return
	}
	// controllo che l'autore sia quello giusto -> prende il nome dell'utente dall'id e lo compara
	username, err := rt.db.GetNameById(id)
	if err != nil {
		http.Error(w, "Internal Serve Error", http.StatusInternalServerError)
		ctx.Logger.Error("Errore durante la selezione del nome utente corrente")
		return
	}
	if username != requestBody.Author {
		http.Error(w, "Autore non corretto", http.StatusBadRequest)
		ctx.Logger.Error("L'utente corrente deve corrispondere all'autore")
		return
	}
	// controllo pattern emoji
	emojiRegex := regexp.MustCompile(`^\p{So}+$`)
	matched := emojiRegex.MatchString(requestBody.Emoji)
	if !matched {
		http.Error(w, "Emoji non accettabile", http.StatusBadRequest)
		ctx.Logger.Error("Emoji contiene caratteri non validi")
		return
	}

	// funzione che crea commento e lo aggiunge ai messaggi (append ai commenti)
	comment, commentId, err := rt.db.CreateComment(requestBody.Author, requestBody.Emoji, messageId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Errore durante la creazione del commento")
		return
	}

	var commentResponse = struct {
		Comment database.Comment `json:"comment"`
	}{
		Comment: comment,
	}

	// preparo l'header
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// ritorno il risultato
	if err := json.NewEncoder(w).Encode(commentResponse); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Failed to encode response")
		return
	}
	// Logga il successo
	ctx.Logger.Infof("Commento inviato con successo: %s, ID: %v, conversazione: %v, messaggio: %v", requestBody.Emoji, commentId, cnvId, messageId)
}
