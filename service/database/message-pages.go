package database

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

type MessagePosition struct {
	ConversationID int       `json:"conversationId"`
	CreatedAt      time.Time `json:"createdAt"`
	MessageID      int       `json:"messageId"`
}
type MessagePage struct {
	Messages       []Message
	DeletedIDs     []int
	Before         *MessagePosition
	HasMore        bool
	Cursor         int64
	HasMoreChanges bool
}

var ErrInvalidCursor = errors.New("invalid message cursor")

func newPage() MessagePage { return MessagePage{Messages: []Message{}, DeletedIDs: []int{}} }
func revision(tx *sql.Tx, conversationID int) (int64, error) {
	var cursor int64
	err := tx.QueryRow("SELECT COALESCE(MAX(revision),0) FROM message_changes WHERE conversation_id=?", conversationID).Scan(&cursor)
	return cursor, err
}
func (db *appdbimpl) GetMessagePage(conversationID int, before *MessagePosition, limit int) (MessagePage, error) {
	if limit < 1 || limit > 100 {
		return MessagePage{}, ErrInvalidCursor
	}
	tx, err := db.c.Begin()
	if err != nil {
		return MessagePage{}, err
	}
	defer tx.Rollback()
	page := newPage()
	page.Cursor, err = revision(tx, conversationID)
	if err != nil {
		return page, err
	}
	query := " WHERE m.message_id IN(SELECT message_id FROM messages_table WHERE conversation_id=?"
	args := []interface{}{conversationID}
	if before != nil {
		if before.ConversationID != conversationID || before.MessageID < 1 || before.CreatedAt.IsZero() {
			return page, ErrInvalidCursor
		}
		query += " AND (created_at<? OR (created_at=? AND message_id<?))"
		timestamp := before.CreatedAt.UTC().Format(timestampLayout)
		args = append(args, timestamp, timestamp, before.MessageID)
	}
	query += " ORDER BY created_at DESC,message_id DESC LIMIT ?)"
	args = append(args, limit+1)
	page.Messages, err = readMessages(tx, query, args...)
	if err != nil {
		return page, err
	}
	if len(page.Messages) > limit {
		page.HasMore = true
		page.Messages = page.Messages[1:]
	}
	if page.HasMore {
		oldest := page.Messages[0]
		page.Before = &MessagePosition{ConversationID: conversationID, CreatedAt: oldest.CreatedAt, MessageID: oldest.MessageId}
	}
	if err = tx.Commit(); err != nil {
		return page, err
	}
	return page, nil
}
func (db *appdbimpl) GetMessageChanges(conversationID int, since int64, limit int) (MessagePage, error) {
	if since < 0 || limit < 1 || limit > 100 {
		return MessagePage{}, ErrInvalidCursor
	}
	tx, err := db.c.Begin()
	if err != nil {
		return MessagePage{}, err
	}
	defer tx.Rollback()
	page := newPage()
	latest, err := revision(tx, conversationID)
	if err != nil {
		return page, err
	}
	if since > latest {
		return page, ErrInvalidCursor
	}
	page.Cursor = since
	rows, err := tx.Query("SELECT revision,message_id FROM message_changes WHERE conversation_id=? AND revision>? ORDER BY revision LIMIT ?", conversationID, since, limit)
	if err != nil {
		return page, err
	}
	ids := make(map[int]bool)
	for rows.Next() {
		var id int
		if err = rows.Scan(&page.Cursor, &id); err != nil {
			rows.Close()
			return page, err
		}
		ids[id] = true
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return page, err
	}
	if len(ids) > 0 {
		placeholders := make([]string, 0, len(ids))
		args := make([]interface{}, 0, len(ids))
		for id := range ids {
			placeholders = append(placeholders, "?")
			args = append(args, id)
		}
		page.Messages, err = readMessages(tx, " WHERE m.message_id IN("+strings.Join(placeholders, ",")+")", args...)
		if err != nil {
			return page, err
		}
		for _, message := range page.Messages {
			delete(ids, message.MessageId)
		}
		for id := range ids {
			page.DeletedIDs = append(page.DeletedIDs, id)
		}
		sort.Ints(page.DeletedIDs)
	}
	page.HasMoreChanges = page.Cursor < latest
	if err = tx.Commit(); err != nil {
		return page, fmt.Errorf("read changes: %w", err)
	}
	return page, nil
}
