package database

import "database/sql"

func messageByID(q queryer, messageID int) (Message, error) {
	messages, err := readMessages(q, " WHERE m.message_id=?", messageID)
	if err != nil {
		return Message{}, err
	}
	if len(messages) == 0 {
		return Message{}, sql.ErrNoRows
	}
	return messages[0], nil
}

func (db *appdbimpl) GetMessage(messageID int) (Message, error) {
	return messageByID(db.c, messageID)
}
