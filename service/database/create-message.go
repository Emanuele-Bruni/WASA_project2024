package database

import (
	"WASA_project2024/service/media"
	"database/sql"
	"time"
)

const timestampLayout = "2006-01-02T15:04:05.000000000Z"

func insertMessage(tx *sql.Tx, text string, conversationID, userID int, forwarded bool, original interface{}, kind string, imageID interface{}, imageName string) (Message, int, error) {
	result, err := tx.Exec(`INSERT INTO messages_table(text,conversation_id,user_id,created_at,forwarded,original_message_id,content_type,image_id,image_name)
 SELECT ?,conversation_id,user_id,?,?,?,?,?,? FROM conversation_members WHERE conversation_id=? AND user_id=?`,
		text, time.Now().UTC().Format(timestampLayout), forwarded, original, kind, imageID, imageName, conversationID, userID)
	if err = requireAffected(result, err); err != nil {
		return Message{}, 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return Message{}, 0, err
	}
	message, err := messageByID(tx, int(id))
	return message, int(id), err
}

func (db *appdbimpl) CreateMessage(text string, conversationID, userID int) (Message, int, error) {
	tx, err := db.c.Begin()
	if err != nil {
		return Message{}, 0, err
	}
	defer tx.Rollback()
	message, id, err := insertMessage(tx, text, conversationID, userID, false, nil, "text", nil, "")
	if err != nil {
		return Message{}, 0, err
	}
	if err = tx.Commit(); err != nil {
		return Message{}, 0, err
	}
	return message, id, nil
}

// ForwardMessage stores a snapshot of the content and a nullable link to the source.
// Deleting the source clears that link while preserving the forwarded flag and content.
func (db *appdbimpl) ForwardMessage(messageID, destinationID, userID int) (Message, int, error) {
	tx, err := db.c.Begin()
	if err != nil {
		return Message{}, 0, err
	}
	defer tx.Rollback()
	var text, kind, imageName string
	var imageID sql.NullInt64
	err = tx.QueryRow(`SELECT m.text,m.content_type,m.image_id,m.image_name FROM messages_table m
 JOIN conversation_members cm ON cm.conversation_id=m.conversation_id AND cm.user_id=?
 WHERE m.message_id=?`, userID, messageID).Scan(&text, &kind, &imageID, &imageName)
	if err != nil {
		return Message{}, 0, err
	}
	message, id, err := insertMessage(tx, text, destinationID, userID, true, messageID, kind, imageID, imageName)
	if err != nil {
		return Message{}, 0, err
	}
	if err = tx.Commit(); err != nil {
		return Message{}, 0, err
	}
	return message, id, nil
}

func (db *appdbimpl) CreateImageMessage(image media.Image, conversationID, userID int) (Message, int, error) {
	tx, err := db.c.Begin()
	if err != nil {
		return Message{}, 0, err
	}
	defer tx.Rollback()
	imageID, err := saveImage(tx, image)
	if err != nil {
		return Message{}, 0, err
	}
	message, id, err := insertMessage(tx, "Immagine", conversationID, userID, false, nil, "image", imageID, image.Name)
	if err != nil {
		return Message{}, 0, err
	}
	if err = tx.Commit(); err != nil {
		return Message{}, 0, err
	}
	return message, id, nil
}
