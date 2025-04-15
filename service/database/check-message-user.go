package database

func (db *appdbimpl) CheckMessage_User(messageId int, userId int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM messages_table WHERE message_id = ? AND user_id=?)"
	err := db.c.QueryRow(query, messageId, userId).Scan(&exists)
	return exists, err
}
