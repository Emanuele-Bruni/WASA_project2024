package database

import (
	"database/sql"
	"time"
)

type MessagePreview struct {
	MessageID int       `json:"messageId"`
	AuthorID  int       `json:"authorId"`
	Author    string    `json:"author"`
	Type      string    `json:"type"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
	Hours     string    `json:"hours"`
	Forwarded bool      `json:"forwarded"`
}

const conversationSelect = `SELECT c.conversation_id,c.title,
 CASE WHEN c.photo_image_id IS NULL THEN '' ELSE '/media/'||c.photo_image_id END,c.is_group,
 u.id,u.name,CASE WHEN u.photo_image_id IS NULL THEN '' ELSE '/media/'||u.photo_image_id END,
 lm.message_id,lm.user_id,la.name,CASE WHEN lm.content_type='image' THEN '' ELSE substr(lm.text,1,200) END,lm.content_type,lm.created_at,lm.forwarded
 FROM conversations_table c
 JOIN conversation_members own ON own.conversation_id=c.conversation_id AND own.user_id=?
 JOIN conversation_members cm ON cm.conversation_id=c.conversation_id
 JOIN users_table u ON u.id=cm.user_id
 LEFT JOIN messages_table lm ON lm.message_id=(SELECT m.message_id FROM messages_table m WHERE m.conversation_id=c.conversation_id ORDER BY m.created_at DESC,m.message_id DESC LIMIT 1)
 LEFT JOIN users_table la ON la.id=lm.user_id`

func (db *appdbimpl) GetConversations(userID int) ([]Conversation, error) {
	return db.readConversations(conversationSelect+" ORDER BY lm.created_at DESC,lm.message_id DESC,c.conversation_id DESC,u.id", userID, userID)
}
func (db *appdbimpl) readConversations(query string, currentUserID int, args ...interface{}) ([]Conversation, error) {
	rows, err := db.c.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]Conversation, 0)
	for rows.Next() {
		var conversation Conversation
		var member User
		var mid, authorID sql.NullInt64
		var author, text, kind, created sql.NullString
		var forwarded sql.NullBool
		if err = rows.Scan(&conversation.ConversationId, &conversation.Title, &conversation.Photo, &conversation.IsGroup, &member.Id, &member.Name, &member.Photo, &mid, &authorID, &author, &text, &kind, &created, &forwarded); err != nil {
			return nil, err
		}
		if len(result) == 0 || result[len(result)-1].ConversationId != conversation.ConversationId {
			conversation.Members = make([]User, 0)
			if mid.Valid {
				timestamp, err := time.Parse(time.RFC3339Nano, created.String)
				if err != nil {
					return nil, err
				}
				conversation.LastMessage = &MessagePreview{MessageID: int(mid.Int64), AuthorID: int(authorID.Int64), Author: author.String, Type: kind.String, Text: text.String, CreatedAt: timestamp, Hours: timestamp.Format("15:04"), Forwarded: forwarded.Bool}
			}
			result = append(result, conversation)
		}
		current := &result[len(result)-1]
		current.Members = append(current.Members, member)
		if !current.IsGroup && member.Id != currentUserID {
			current.Title, current.Photo = member.Name, member.Photo
		}
	}
	return result, rows.Err()
}
