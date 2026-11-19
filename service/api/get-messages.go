package api

import (
	"WASA_project2024/service/api/reqcontext"
	"WASA_project2024/service/database"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"time"
)

func (rt *_router) getMessages(w http.ResponseWriter, r *http.Request, p httprouter.Params, ctx reqcontext.RequestContext) {
	query := r.URL.Query()
	for key, values := range query {
		if (key != "limit" && key != "before" && key != "since") || len(values) != 1 {
			writeError(w, 400, "Parametri della cronologia non validi")
			return
		}
	}
	limit := 50
	if raw, ok := query["limit"]; ok {
		parsed, err := strconv.Atoi(raw[0])
		if err != nil || parsed < 1 || parsed > 100 {
			writeError(w, 400, "limit deve essere compreso tra 1 e 100", "limit")
			return
		}
		limit = parsed
	}
	var page database.MessagePage
	var err error
	if raw, ok := query["since"]; ok {
		if query.Has("before") {
			writeError(w, 400, "before e since non possono essere combinati")
			return
		}
		cursor, parseError := strconv.ParseInt(raw[0], 10, 64)
		if parseError != nil || cursor < 0 {
			writeError(w, 400, "Cursore di aggiornamento non valido", "since")
			return
		}
		page, err = rt.db.GetMessageChanges(ctx.ConversationID, cursor, limit)
	} else {
		var before *database.MessagePosition
		if raw, ok := query["before"]; ok {
			if len(raw[0]) > 256 {
				writeError(w, 400, "Cursore pagina non valido", "before")
				return
			}
			bytes, decodeError := base64.RawURLEncoding.Strict().DecodeString(raw[0])
			var position database.MessagePosition
			if decodeError != nil || json.Unmarshal(bytes, &position) != nil || position.ConversationID != ctx.ConversationID || position.MessageID < 1 || position.CreatedAt.IsZero() || position.CreatedAt.Location() != time.UTC {
				writeError(w, 400, "Cursore pagina non valido", "before")
				return
			}
			before = &position
		}
		page, err = rt.db.GetMessagePage(ctx.ConversationID, before, limit)
	}
	if errors.Is(err, database.ErrInvalidCursor) {
		writeError(w, 400, "Cursore non valido")
		return
	}
	if err != nil {
		respondDatabaseError(w, err, ctx)
		return
	}
	var before *string
	if page.Before != nil {
		encoded, _ := json.Marshal(page.Before)
		value := base64.RawURLEncoding.EncodeToString(encoded)
		before = &value
	}
	writeJSON(w, 200, struct {
		Messages       []database.Message `json:"messages"`
		DeletedIDs     []int              `json:"deletedIds"`
		Cursor         string             `json:"cursor"`
		Before         *string            `json:"before"`
		HasMore        bool               `json:"hasMore"`
		HasMoreChanges bool               `json:"hasMoreChanges"`
	}{page.Messages, page.DeletedIDs, strconv.FormatInt(page.Cursor, 10), before, page.HasMore, page.HasMoreChanges}, ctx)
}
