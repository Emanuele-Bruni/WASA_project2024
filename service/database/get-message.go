package database

import (
	"database/sql"
	"strconv"
	"strings"
)

func (db *appdbimpl) GetMessage(messageId int) (Message, error) {
	// prendo il messaggio dalla tabella
	var message Message
	var commentsString sql.NullString
	err := db.c.QueryRow("SELECT author,text,hours,comments FROM messages_table WHERE message_id=?", messageId).Scan(&message.Author, &message.Text, &message.Hours, &commentsString)
	if err != nil {
		return Message{}, err
	}

	// calcolo i commenti
	if commentsString.Valid && commentsString.String != "" {
		commentIDs := strings.Split(commentsString.String, ",")

		for _, idStringa := range commentIDs {
			commentID, err := strconv.Atoi(strings.TrimSpace(idStringa))
			if err != nil {
				return Message{}, err
			}

			var comment Comment
			comment.CommentId = commentID
			err = db.c.QueryRow("SELECT author, emoji FROM comments_table WHERE comment_id=?", commentID).Scan(&comment.Author, &comment.Emoji)
			if err != nil {
				return Message{}, err
			} else {
				message.Comments = append(message.Comments, comment)
			}
		}

	}
	return message, nil
}
