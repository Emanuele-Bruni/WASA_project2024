package api

import (
	"WASA_project2024/service/api/reqcontext"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, p httprouter.Params, ctx reqcontext.RequestContext) {
	//request body -> messaggio da inviare
	var requestBody = struct {
		Author string `json:"author"`
		Text   string `json:"text"`
		Hours  string `json:"hours"`
	}{}

	//estraggo userId e controllo se è autorizzato
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "bearerAuth ")
	id, err := strconv.Atoi(token)
	if err != nil {
		ctx.Logger.Error("Token must be a number: ", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	//decodifico la richiesta secondo il request body
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Struttura non compatibile", http.StatusBadRequest)           //tipo di errore
		ctx.Logger.WithError(err).Error("Input non compatibile con il RequestBody") //descrizione errore
		return
	}

	// Controllo che i campi non siano vuoti
	if requestBody.Author == "" || requestBody.Text == "" || requestBody.Hours == "" {
		http.Error(w, "I campi author, text e hours sono obbligatori", http.StatusBadRequest)
		ctx.Logger.Error("Campi mancanti nel request body")
		return
	}

	//prendo il conversation id dal path
	cnvId, err := strconv.Atoi(p.ByName("conversationId"))
	if err != nil {
		//errore nel prendere il parametro
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

	//controllo se l'utente fa parte della conversazione
	//funzione che controlla se l'utente è nella conversazione
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

	// FUNZIONE "create-message" che ritorna il messaggio (senza commenti)
	messageResponse, err := rt.db.CreateMessage(requestBody.Author, requestBody.Text, requestBody.Hours, cnvId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Errore durante la creazione del messaggio")
		return
	}

	//preparo l'header
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	//ritorno il risultato
	if err := json.NewEncoder(w).Encode(messageResponse); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Failed to encode response")
		return
	}
	//Logga il successo
	ctx.Logger.Infof("Messaggio inviato con successo: %s", requestBody.Text)

}
