package api

import (
	"WASA_project2024/service/api/reqcontext"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type accessPolicy int

const (
	memberAccess accessPolicy = iota
	groupAccess
	messageAuthorAccess
	commentAuthorAccess
	ownProfileAccess
)

// requireCheck keeps database failures separate from missing resources and denied permissions.
func requireCheck(w http.ResponseWriter, ctx *reqcontext.RequestContext, ok bool, err error, status int) bool {
	if err != nil {
		ctx.Logger.WithError(err).Error("Authorization check failed")
		respondDatabaseError(w, err, *ctx)
		return false
	}
	if !ok {
		writeError(w, status, http.StatusText(status))
		return false
	}
	return true
}

func (rt *_router) requireConversation(w http.ResponseWriter, ctx *reqcontext.RequestContext, id int) bool {
	exists, err := rt.db.ConversationExists(id)
	if !requireCheck(w, ctx, exists, err, http.StatusNotFound) {
		return false
	}
	member, err := rt.db.CheckUser_Conversation(id, ctx.UserID)
	return requireCheck(w, ctx, member, err, http.StatusForbidden)
}

func (rt *_router) authorize(w http.ResponseWriter, ps httprouter.Params, ctx *reqcontext.RequestContext, policy accessPolicy) bool {
	for _, param := range []struct {
		name string
		id   *int
	}{
		{"conversationid", &ctx.ConversationID},
		{"messageid", &ctx.MessageID},
		{"commentid", &ctx.CommentID},
		{"forwardid", &ctx.ForwardID},
		{"userid", &ctx.TargetUserID},
		{"mediaid", &ctx.MediaID},
	} {
		if raw := ps.ByName(param.name); raw != "" {
			id, err := strconv.Atoi(raw)
			if err != nil || id < 1 {
				writeError(w, http.StatusBadRequest, "Parametro non valido", param.name)
				return false
			}
			*param.id = id
		}
	}
	if ctx.ConversationID != 0 && !rt.requireConversation(w, ctx, ctx.ConversationID) {
		return false
	}
	if policy == groupAccess {
		group, err := rt.db.Conversation_isGroup(ctx.ConversationID)
		if !requireCheck(w, ctx, group, err, http.StatusForbidden) {
			return false
		}
	}
	if ctx.MessageID != 0 {
		exists, err := rt.db.MessageExists(ctx.MessageID)
		if !requireCheck(w, ctx, exists, err, http.StatusNotFound) {
			return false
		}
		belongs, err := rt.db.CheckMessage_Conversation(ctx.MessageID, ctx.ConversationID)
		if !requireCheck(w, ctx, belongs, err, http.StatusForbidden) {
			return false
		}
	}
	if policy == messageAuthorAccess {
		author, err := rt.db.CheckMessage_User(ctx.MessageID, ctx.UserID)
		if !requireCheck(w, ctx, author, err, http.StatusForbidden) {
			return false
		}
	}
	if ctx.CommentID != 0 {
		exists, err := rt.db.CommentExists(ctx.CommentID)
		if !requireCheck(w, ctx, exists, err, http.StatusNotFound) {
			return false
		}
		belongs, err := rt.db.CheckComment_Message(ctx.CommentID, ctx.MessageID)
		if !requireCheck(w, ctx, belongs, err, http.StatusForbidden) {
			return false
		}
	}
	if policy == commentAuthorAccess {
		author, err := rt.db.CheckAuthorComment(ctx.CommentID, ctx.UserID)
		if !requireCheck(w, ctx, author, err, http.StatusForbidden) {
			return false
		}
	}
	if ctx.ForwardID != 0 && !rt.requireConversation(w, ctx, ctx.ForwardID) {
		return false
	}
	if policy == ownProfileAccess && ctx.TargetUserID != ctx.UserID {
		writeError(w, http.StatusForbidden, "Utente non autorizzato")
		return false
	}
	if ctx.TargetUserID != 0 {
		exists, err := rt.db.CheckUserId(ctx.TargetUserID)
		if !requireCheck(w, ctx, exists, err, http.StatusNotFound) {
			return false
		}
	}
	if ctx.MediaID != 0 {
		exists, err := rt.db.ImageExists(ctx.MediaID)
		if !requireCheck(w, ctx, exists, err, http.StatusNotFound) {
			return false
		}
		allowed, err := rt.db.CanReadImage(ctx.MediaID, ctx.UserID)
		if !requireCheck(w, ctx, allowed, err, http.StatusForbidden) {
			return false
		}
	}

	return true
}
