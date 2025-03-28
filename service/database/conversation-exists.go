package database

func (db *appdbimpl) ConversationExists(cnvId int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM conversations_table WHERE conversation_id = ?)"
	err := db.c.QueryRow(query, cnvId).Scan(&exists)
	return exists, err
}
