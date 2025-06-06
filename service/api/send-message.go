package api

import (
	"WASA_project2024/service/api/reqcontext"
	"WASA_project2024/service/database"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, p httprouter.Params, ctx reqcontext.RequestContext) {

	// estraggo userId e controllo se è autorizzato
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	id, err := strconv.Atoi(token)
	if err != nil {
		ctx.Logger.Error("Token must be a number: ", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	// funzione che ritorna il nome dell'userId
	// autore = username dell'utente
	name, err := rt.db.GetNameById(id)
	if err != nil {
		http.Error(w, "Autore non identificato", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Errore nel recuperare il nome dell'utente")
		return
	}

	// request body -> messaggio da inviare
	var requestBody = struct {
		Text string `json:"text"`
	}{}
	// decodifico la richiesta secondo il request body
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Struttura non compatibile", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Input non compatibile con il RequestBody")
		return
	}

	// prendo il conversation id dal path
	cnvId, err := strconv.Atoi(p.ByName("conversationid"))
	if err != nil {
		// errore nel prendere il parametro
		http.Error(w, "Parametro non valido", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("conversationId parametro non valido")
		return
	}

	// Verifica se la conversazione esiste prima di inserire il messaggio
	exists, err := rt.db.ConversationExists(cnvId)
	if err != nil {
		http.Error(w, "Errore interno nel controllo della conversazione", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Errore durante la verifica dell'esistenza della conversazione")
		return
	}
	if !exists {
		http.Error(w, "La conversazione non esiste", http.StatusNotFound)
		ctx.Logger.Error("Tentativo di inviare un messaggio a una conversazione inesistente")
		return
	}

	// controllo se l'utente fa parte della conversazione
	// funzione che controlla se l'utente è nella conversazione
	found, err := rt.db.CheckUser_Conversation(cnvId, id)
	if err != nil {
		http.Error(w, "Errore nel controllo del membro", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Errore durante la verifica dell'appartenenza dell'utente alla conversazione")
		return
	}
	if !found {
		http.Error(w, "Utente non appartiene alla conversazione", http.StatusForbidden)
		ctx.Logger.Error("Errore perché l'utente non appartiene alla conversazione")
	}

	// Controllo che i campi non siano vuoti
	if requestBody.Text == "" {
		http.Error(w, "Impossibile inviare un messaggio vuoto", http.StatusBadRequest)
		ctx.Logger.Error("Testo mancante nel request body")
		return
	}

	// calcolo automaticamente l'orario attuale in formato HH:mm
	hours := time.Now().Format("15:04")

	// FUNZIONE "create-message" che ritorna il messaggio (senza commenti) e il suo id
	message, msgId, err := rt.db.CreateMessage(name, requestBody.Text, hours, cnvId, id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Errore durante la creazione del messaggio")
		return
	}

	var messageResponse = struct {
		FinalMessage database.Message `json:"message"`
	}{
		FinalMessage: message,
	}

	// preparo l'header
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// ritorno il risultato
	if err := json.NewEncoder(w).Encode(messageResponse); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Failed to encode response")
		return
	}
	// Logga il successo
	ctx.Logger.Infof("Messaggio inviato con successo: messaggio: %v, conversazione: %v", msgId, cnvId)

}
