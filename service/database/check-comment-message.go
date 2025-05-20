package database

func (db *appdbimpl) CheckComment_Message(commentId int, messageId int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM comments_table WHERE comment_id=? AND message_id=?)"
	err := db.c.QueryRow(query, commentId, messageId).Scan(&exists)
	return exists, err
}
