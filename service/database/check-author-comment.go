package database

import "fmt"

func (db *appdbimpl) CheckAuthorComment(commentId int, username string) (bool, error) {
	//prendo l'autore del commento
	var author string
	err := db.c.QueryRow("SELECT author FROM comments_table WHERE comment_id=?", commentId).Scan(&author)
	if err != nil {
		return false, err
	}
	//lo paragono all'username
	if author != username {
		return false, fmt.Errorf("autore e utente non coincidono")
	}
	return true, nil
}
