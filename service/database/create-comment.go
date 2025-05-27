package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// crea un commento, lo aggiunge alla tabella dei commenti e alla lista dei commenti del messaggio
func (db *appdbimpl) CreateComment(author string, emoji string, messageId int) (Comment, int, error) {
	// inserisco il commento nella tabella dei commenti
	res, err := db.c.Exec("INSERT INTO comments_table (author,emoji, message_id) VALUES (?,?,?)", author, emoji, messageId)
	if err != nil {
		return Comment{}, -1, fmt.Errorf("errore nel creare il commento: %w", err)
	}
	// calcolo l'id del commento
	commentId, err := res.LastInsertId()
	if err != nil {
		return Comment{}, -1, fmt.Errorf("errore nel calcolare l'ID del commento: %w", err)
	}

	// determino il commento
	finalComment := Comment{
		CommentId: int(commentId),
		Author:    author,
		Emoji:     emoji,
	}

	// recupero l'array con gli id dei commenti
	var currentComments sql.NullString //stringa "20,30,40"
	err = db.c.QueryRow("SELECT comments FROM messages_table WHERE message_id=?", messageId).Scan(&currentComments)
	if err != nil {
		return Comment{}, -1, fmt.Errorf("errore nel recuperare i commenti attuali: %w", err)
	}

	if currentComments.Valid && currentComments.String != "" {
		var commentsidList []int //[20,30,40]
		for _, id_stringa := range strings.Split(currentComments.String, ",") {
			// converto in numero
			num, err := strconv.Atoi(id_stringa)
			if err != nil {
				logrus.WithField("comment_id", id_stringa).WithError(err).Error("id non valido")
			} else {
				commentsidList = append(commentsidList, num)
			}
		}
		// aggiungo il commento creato
		commentsidList = append(commentsidList, finalComment.CommentId)
		// converto la lista in stringa
		updatedComments := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(commentsidList)), ","), "[]")
		// aggiorno la colonna "comments" del database
		_, err = db.c.Exec("UPDATE messages_table SET comments=? WHERE message_id=?", updatedComments, messageId)
		if err != nil {
			return Comment{}, -1, fmt.Errorf("errore nell'aggiornare la lista dei commenti: %w", err)
		}
	}

	return finalComment, int(commentId), nil
}
