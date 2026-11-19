package database

func (db *appdbimpl) AddUserToGroup(userID, conversationID int) error {
	return requireAffected(db.c.Exec(`INSERT INTO conversation_members(conversation_id,user_id)
 SELECT conversation_id,? FROM conversations_table WHERE conversation_id=? AND is_group=1`, userID, conversationID))
}
