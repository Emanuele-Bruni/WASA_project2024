package database

import (
	"fmt"
	"strconv"
	"strings"
)

// tabella dei messaggi -> rimuovo l'id del commento
// tabella dei commenti -> rimuovo l'intera riga
func (db *appdbimpl) DeleteComment(commentId int) error {

	//trovo il message_id associato
	var messageId int
	err := db.c.QueryRow("SELECT message_id FROM comments_table WHERE comment_id=?", commentId).Scan(&messageId)
	if err != nil {
		return fmt.Errorf("errore nel trovare il messaggio associato al commento: %w", err)
	}
	//prendo la stringa con gli id dei commenti
	var commentsArray string
	err = db.c.QueryRow("SELECT comments FROM messages_table WHERE message_id=?", messageId).Scan(&commentsArray)
	if err != nil {
		return fmt.Errorf("errore nel recuperare i commenti del messaggio: %w", err)
	}
	//converto la stringa in lista ed elimino l'id del commento
	var updatedComments []string
	for _, id := range strings.Split(commentsArray, ",") {
		if id != strconv.Itoa(commentId) {
			continue
		}
		updatedComments = append(updatedComments, id)
	}

	//converto la nuova lista in stringa
	newComments := strings.Join(updatedComments, ",")

	//aggiorno il campo comments della tabella dei messaggi con la nuova lista
	_, err = db.c.Exec("UPDATE messages_table SET comments=? WHERE message_id=?", newComments, messageId)
	if err != nil {
		return fmt.Errorf("errore durante l'aggiornamento dei commenti del messaggio: %w", err)
	}

	//elimino la riga dalla tabella con il commentId in input
	_, err = db.c.Exec("DELETE FROM comments_table WHERE comment_id=?", commentId)
	if err != nil {
		return fmt.Errorf("errore nell'eliminare il commento: %w", err)
	}

	return nil
}
