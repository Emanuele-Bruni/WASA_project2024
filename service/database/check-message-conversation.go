package database

func (db *appdbimpl) CheckMessage_Conversation(messageId int, cnvId int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM messages_table WHERE message_id = ? AND conversation_id=?)"
	err := db.c.QueryRow(query, messageId, cnvId).Scan(&exists)
	return exists, err
}
