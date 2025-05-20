package database

func (db *appdbimpl) MessageExists(messageId int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM messages_table WHERE message_id=?)"
	err := db.c.QueryRow(query, messageId).Scan(&exists)
	return exists, err
}
