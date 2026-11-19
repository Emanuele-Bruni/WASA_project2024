package database

func (db *appdbimpl) Conversation_isGroup(conversationId int) (bool, error) {
	var isGroup bool
	query := "SELECT EXISTS(SELECT 1 FROM conversations_table WHERE conversation_id=? AND is_group=?)"
	err := db.c.QueryRow(query, conversationId, 1).Scan(&isGroup)
	return isGroup, err
}
