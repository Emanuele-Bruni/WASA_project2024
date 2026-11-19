package database

func (db *appdbimpl) DeleteComment(commentID int) error {
	return requireAffected(db.c.Exec("DELETE FROM comments_table WHERE comment_id=?", commentID))
}
