package database

//seleziono la conversazione dalla tabella delle conversazioni
func (db *appdbimpl) GetSpecificConversation(conversationId int) (Conversation, error) {
	var conversation Conversation
	err := db.c.QueryRow("SELECT conversation_id,titolo,photo,isGroup,members FROM conversations_table WHERE conversation_id=?", conversationId).Scan(&conversation)
	return conversation, err
}
