package database

import (
	"database/sql"
	"errors"
)

// SetReaction atomically creates or replaces the single reaction of an author on a message.
func (db *appdbimpl) SetReaction(author int, emoji string, messageID int) (Comment, bool, error) {
	tx, err := db.c.Begin()
	if err != nil {
		return Comment{}, false, err
	}
	defer tx.Rollback()
	var existing int
	lookup := tx.QueryRow("SELECT comment_id FROM comments_table WHERE message_id=? AND author_id=?", messageID, author).Scan(&existing)
	created := errors.Is(lookup, sql.ErrNoRows)
	if lookup != nil && !created {
		return Comment{}, false, lookup
	}
	if err = requireAffected(tx.Exec(`INSERT INTO comments_table(message_id,author_id,emoji)
 SELECT m.message_id,?,? FROM messages_table m
 JOIN conversation_members cm ON cm.conversation_id=m.conversation_id AND cm.user_id=?
 WHERE m.message_id=?
 ON CONFLICT(message_id,author_id) DO UPDATE SET emoji=excluded.emoji`, author, emoji, author, messageID)); err != nil {
		return Comment{}, false, err
	}
	var comment Comment
	if err = tx.QueryRow("SELECT comment_id,author_id,emoji FROM comments_table WHERE message_id=? AND author_id=?", messageID, author).Scan(&comment.CommentId, &comment.AuthorId, &comment.Emoji); err != nil {
		return Comment{}, false, err
	}
	if err = tx.Commit(); err != nil {
		return Comment{}, false, err
	}
	return comment, created, nil
}
