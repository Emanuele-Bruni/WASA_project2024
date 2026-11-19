package database

func (db *appdbimpl) RemoveMessage(messageID int) error {
	return requireAffected(db.c.Exec("DELETE FROM messages_table WHERE message_id=?", messageID))
}
