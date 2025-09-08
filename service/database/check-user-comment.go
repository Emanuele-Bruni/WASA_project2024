package database

func (db *appdbimpl) CheckUserComment(userId int, messageId int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM comments_table WHERE message_id = ? AND author_id=?)"
	err := db.c.QueryRow(query, messageId, userId).Scan(&exists)
	return exists, err
}
