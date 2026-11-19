package api

import (
	"WASA_project2024/service/api/reqcontext"
	"WASA_project2024/service/database"
	"encoding/json"
	"net/http"
)

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

func writeError(w http.ResponseWriter, status int, message string, field ...string) {
	codes := map[int]string{400: "invalid_request", 401: "unauthorized", 403: "forbidden", 404: "not_found", 405: "method_not_allowed", 409: "conflict", 413: "payload_too_large", 415: "unsupported_media_type", 500: "internal_error"}
	code := codes[status]
	if code == "" {
		code = "internal_error"
	}
	if message == "" {
		message = http.StatusText(status)
	}
	failure := apiError{Code: code, Message: message}
	if len(field) > 0 {
		failure.Field = field[0]
	}
	body, _ := json.Marshal(struct {
		Error apiError `json:"error"`
	}{failure})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_, _ = w.Write(append(body, '\n'))
}

func writeJSON(w http.ResponseWriter, status int, value interface{}, ctx reqcontext.RequestContext) {
	body, err := json.Marshal(value)
	if err != nil {
		ctx.Logger.WithError(err).Error("Serialize API response")
		writeError(w, 500, "Errore interno del server")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	if _, err = w.Write(append(body, '\n')); err != nil {
		ctx.Logger.WithError(err).Error("Write API response")
	}
}

func respondDatabaseError(w http.ResponseWriter, err error, ctx reqcontext.RequestContext) {
	switch database.ClassifyError(err) {
	case database.MissingResource:
		writeError(w, 404, "Risorsa non trovata")
	case database.ConflictingResource:
		writeError(w, 409, "La risorsa esiste già o il valore è già utilizzato")
	case database.InvalidData:
		writeError(w, 400, "Dati non validi")
	default:
		ctx.Logger.WithError(err).Error("Database operation failed")
		writeError(w, 500, "Errore interno del server")
	}
}

type userResponse struct {
	User database.User `json:"user"`
}
type conversationResponse struct {
	Conversation database.Conversation `json:"conversation"`
}
type messageResponse struct {
	Message database.Message `json:"message"`
}
type commentResponse struct {
	Comment database.Comment `json:"comment"`
}

func (rt *_router) respondUser(w http.ResponseWriter, ctx reqcontext.RequestContext) {
	user, err := rt.db.GetUser(ctx.UserID)
	if err != nil {
		respondDatabaseError(w, err, ctx)
		return
	}
	writeJSON(w, 200, userResponse{user}, ctx)
}
func (rt *_router) respondConversation(w http.ResponseWriter, status, id int, ctx reqcontext.RequestContext) {
	conversation, err := rt.db.GetSpecificConversation(id, ctx.UserID)
	if err != nil {
		respondDatabaseError(w, err, ctx)
		return
	}
	writeJSON(w, status, conversationResponse{conversation}, ctx)
}
