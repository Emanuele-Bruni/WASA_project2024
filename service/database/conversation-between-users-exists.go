package database

import (
	"database/sql"
	"fmt"
)

func (db *appdbimpl) Conversation2MembersExists(userId1 int, userId2 int) (bool, int, error) {

	if userId1 == userId2 {
		return false, 0, fmt.Errorf("conversazione con se stessi non ammessa")
	}

	members1 := fmt.Sprintf("%d,%d", userId1, userId2)
	members2 := fmt.Sprintf("%d,%d", userId2, userId1)

	var conversationId int
	// seleziono la conversazione
	query := "SELECT conversation_id FROM conversations_table WHERE (members = ? OR members = ?)AND isGroup = 0"
	err := db.c.QueryRow(query, members1, members2).Scan(&conversationId)
	if err != nil {
		if err == sql.ErrNoRows {
			// Non esiste conversazione
			return false, 0, nil
		}
		// Errore DB
		return false, 0, err
	}
	// conversazione trovata
	return true, conversationId, nil
}
