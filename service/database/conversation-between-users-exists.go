package database

import (
	"database/sql"
	"fmt"
)

func (db *appdbimpl) Conversation2MembersExists(userId1 int, userId2 int) (bool, int, error) {

	if userId1 == userId2 {
		return false, 0, fmt.Errorf("conversazione con se stessi non ammessa")
	}

	var membri string
	if userId1 < userId2 {
		membri = fmt.Sprintf("%d,%d", userId1, userId2) //"10,20"
	} else {
		membri = fmt.Sprintf("%d,%d", userId2, userId1)
	}

	var conversationId int
	//seleziono la conversazione
	query := "SELECT conversation_id FROM conversations_table WHERE members = ? AND isGroup = 0"
	err := db.c.QueryRow(query, membri).Scan(&conversationId)
	if err != nil {
		if err == sql.ErrNoRows {
			// Non esiste conversazione
			return false, 0, nil
		}
		// Errore DB
		return false, 0, err
	}
	//conversazione trovata
	return true, conversationId, err
}
