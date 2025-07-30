package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

func (db *appdbimpl) GetConversationMessages(cnvId int) ([]Message, error) {
	righe, err := db.c.Query("SELECT message_id, author, text, hours, comments, user_id FROM messages_table WHERE conversation_id = ?", cnvId)
	if err != nil {
		return nil, fmt.Errorf("errore nella selezione dei messaggi: %w", err)
	}
	defer righe.Close()

	var result []Message

	for righe.Next() {
		var msgID, authorId int
		var author, text, hours string
		var comments sql.NullString

		if err := righe.Scan(&msgID, &author, &text, &hours, &comments, &authorId); err != nil {
			return nil, fmt.Errorf("errore nella scansione del messaggio: %w", err)
		}

		// Usa una stringa vuota se comments è NULL
		commentsStr := ""
		if comments.Valid {
			commentsStr = comments.String
		}

		var tempMessage Message
		tempMessage.MessageId = msgID
		tempMessage.AuthorId = authorId
		tempMessage.Author = author
		tempMessage.Text = text
		tempMessage.Hours = hours

		var allComments []Comment
		var commentIds []int

		if strings.TrimSpace(commentsStr) != "" {
			parts := strings.Split(commentsStr, ",")
			for _, p := range parts {
				n, err := strconv.Atoi(strings.TrimSpace(p))
				if err != nil {
					return nil, fmt.Errorf("errore nella conversione dell'id del commento: %w", err)
				}
				commentIds = append(commentIds, n)
			}
		}

		for _, n := range commentIds {
			var tempComment Comment
			err := db.c.QueryRow("SELECT comment_id, author, emoji FROM comments_table WHERE comment_id = ?", n).
				Scan(&tempComment.CommentId, &tempComment.Author, &tempComment.Emoji)
			if err != nil {
				return nil, fmt.Errorf("errore nella lettura del commento con id %d: %w", n, err)
			}
			allComments = append(allComments, tempComment)
		}

		tempMessage.Comments = allComments
		result = append(result, tempMessage)
	}

	if err := righe.Err(); err != nil {
		return nil, fmt.Errorf("errore nella lettura delle righe: %w", err)
	}

	return result, nil
}
