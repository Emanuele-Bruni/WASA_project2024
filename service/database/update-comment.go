package database

import "fmt"

func (db *appdbimpl) UpdateComment(messageId int, author int, emoji string) (Comment, int, error) {
	// aggiorno il commento nella tabella dei commenti (non cambio id)
	_, err := db.c.Exec("UPDATE comments_table SET emoji=? WHERE author_id=? AND message_id=?", emoji, author, messageId)
	if err != nil {
		return Comment{}, -1, fmt.Errorf("errore nell'aggiornare il contenuto del commento: %w", err)
	}
	// recupero id del commento
	var commentId int
	err = db.c.QueryRow("SELECT comment_id FROM comments_table WHERE author_id=? AND message_id=?", author, messageId).Scan(&commentId)
	if err != nil {
		return Comment{}, -1, fmt.Errorf("errore nel selezionare l'id del nuovo commento: %w", err)
	}
	// restituisco il commento aggiornato
	newComment := Comment{
		CommentId: int(commentId),
		AuthorId:  author,
		Emoji:     emoji,
	}
	return newComment, commentId, nil
}
