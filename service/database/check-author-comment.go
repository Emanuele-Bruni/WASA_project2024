package database

import "fmt"

func (db *appdbimpl) CheckAuthorComment(commentId int, userId int) (bool, error) {
	// prendo l'autore del commento
	var authorId int
	err := db.c.QueryRow("SELECT author_id FROM comments_table WHERE comment_id=?", commentId).Scan(&authorId)
	if err != nil {
		return false, err
	}
	// lo paragono all'username
	if authorId != userId {
		return false, fmt.Errorf("autore e utente non coincidono")
	}
	return true, nil
}
