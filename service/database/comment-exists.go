package database

func (db *appdbimpl) CommentExists(commentId int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM comments_table WHERE comment_id=?)"
	err := db.c.QueryRow(query, commentId).Scan(&exists)
	return exists, err
}
