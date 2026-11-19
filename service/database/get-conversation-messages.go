package database

import (
	"database/sql"
	"time"
)

type queryer interface {
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
}

const messageSelect = `SELECT m.message_id,m.user_id,u.name,CASE WHEN m.content_type='image' THEN '' ELSE m.text END,m.created_at,m.forwarded,m.original_message_id,
 r.comment_id,r.author_id,r.emoji,m.content_type,i.image_id,i.mime_type,CASE WHEN m.image_name='' THEN 'Immagine' ELSE m.image_name END,i.width,i.height,length(i.data)
 FROM messages_table m JOIN users_table u ON u.id=m.user_id
 LEFT JOIN comments_table r ON r.message_id=m.message_id
 LEFT JOIN images i ON i.image_id=m.image_id`

func readMessages(q queryer, filter string, args ...interface{}) ([]Message, error) {
	rows, err := q.Query(messageSelect+filter+" ORDER BY m.created_at,m.message_id,r.comment_id", args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]Message, 0)
	for rows.Next() {
		var message Message
		var timestamp string
		var original, commentID, authorID sql.NullInt64
		var emoji sql.NullString
		var imageID, width, height, size sql.NullInt64
		var mime, name sql.NullString
		if err = rows.Scan(&message.MessageId, &message.AuthorId, &message.Author, &message.Text, &timestamp, &message.Forwarded, &original, &commentID, &authorID, &emoji, &message.Type, &imageID, &mime, &name, &width, &height, &size); err != nil {
			return nil, err
		}
		if len(result) == 0 || result[len(result)-1].MessageId != message.MessageId {
			message.CreatedAt, err = time.Parse(time.RFC3339Nano, timestamp)
			if err != nil {
				return nil, err
			}
			message.Hours = message.CreatedAt.Format("15:04")
			message.Comments = make([]Comment, 0)
			if imageID.Valid {
				message.Image = &ImageMetadata{ImageID: int(imageID.Int64), URL: imageURL(int(imageID.Int64)), MimeType: mime.String, Name: name.String, Width: int(width.Int64), Height: int(height.Int64), Size: int(size.Int64)}
			}
			if original.Valid {
				id := int(original.Int64)
				message.OriginalMessageID = &id
			}
			result = append(result, message)
		}
		if commentID.Valid {
			current := &result[len(result)-1]
			current.Comments = append(current.Comments, Comment{CommentId: int(commentID.Int64), AuthorId: int(authorID.Int64), Emoji: emoji.String})
		}
	}
	return result, rows.Err()
}

func (db *appdbimpl) GetConversationMessages(conversationID int) ([]Message, error) {
	return readMessages(db.c, " WHERE m.conversation_id=?", conversationID)
}
