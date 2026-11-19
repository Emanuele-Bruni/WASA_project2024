package database

import "database/sql"

func (db *appdbimpl) GetSpecificConversation(conversationID, currentUserID int) (Conversation, error) {
	result, err := db.readConversations(conversationSelect+" WHERE c.conversation_id=? ORDER BY c.conversation_id,u.id", currentUserID, currentUserID, conversationID)
	if err != nil {
		return Conversation{}, err
	}
	if len(result) == 0 {
		return Conversation{}, sql.ErrNoRows
	}
	return result[0], nil
}
