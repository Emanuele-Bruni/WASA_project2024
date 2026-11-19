// Package database owns the SQLite schema, migrations and persistence operations.
// Open configures SQLite connections; New applies embedded migrations before serving requests.
package database

import (
	"WASA_project2024/service/media"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type User struct {
	Id    int    `json:"userId"`
	Name  string `json:"name"`
	Photo string `json:"photo"`
}
type Conversation struct {
	ConversationId int             `json:"conversationId"`
	Title          string          `json:"title"`
	Photo          string          `json:"photo"`
	IsGroup        bool            `json:"isGroup"`
	Members        []User          `json:"members"`
	LastMessage    *MessagePreview `json:"lastMessage"`
}

type Message struct {
	MessageId         int            `json:"messageId"`
	Type              string         `json:"type"`
	Image             *ImageMetadata `json:"image"`
	CreatedAt         time.Time      `json:"createdAt"`
	Forwarded         bool           `json:"forwarded"`
	OriginalMessageID *int           `json:"originalMessageId"`
	AuthorId          int            `json:"authorId"`
	Author            string         `json:"author"`
	Text              string         `json:"text"`
	Hours             string         `json:"hours"`
	Comments          []Comment      `json:"comments"`
}

type Comment struct {
	CommentId int    `json:"commentId"`
	AuthorId  int    `json:"authorId"`
	Emoji     string `json:"emoji"`
}

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetOrCreateUser(name string) (int, error)
	SetName(name string, id int) error
	GetConversations(id int) ([]Conversation, error)
	GetSpecificConversation(conversationId int, currentUserId int) (Conversation, error)
	CreateMessage(text string, cnvid int, userid int) (Message, int, error)
	ForwardMessage(messageID int, destinationID int, userID int) (Message, int, error)
	ConversationExists(cnvId int) (bool, error)
	CommentExists(commentId int) (bool, error)
	MessageExists(messageId int) (bool, error)
	CheckUser_Conversation(cnvId int, userId int) (bool, error)
	CheckMessage_Conversation(messageId int, cnvId int) (bool, error)
	GetMessage(messageId int) (Message, error)
	SetReaction(author int, emoji string, messageID int) (Comment, bool, error)
	CheckComment_Message(commentId int, messageId int) (bool, error)
	DeleteComment(commentId int) error
	CheckUserId(id int) (bool, error)
	CheckMessage_User(messageId int, userId int) (bool, error)
	RemoveMessage(messageId int) error
	AddUserToGroup(userId int, cnvId int) error
	Conversation_isGroup(conversationId int) (bool, error)
	GetConversationMembers(cnvId int) ([]User, error)
	LeaveGroup(userID int, conversationID int) error
	SetGroupName(groupName string, cnvId int) error
	SetMyPhoto(photo media.Image, id int) error
	SetGroupPhoto(photo media.Image, cnvId int) error
	GetUsers(userID int) ([]User, error)
	GetUser(userID int) (User, error)
	CreateConversation(userId1 int, userId2 int) (int, error)
	CreateGroup(titolo string, members []int, userId int) (int, error)
	Conversation2MembersExists(userId1 int, userId2 int) (bool, int, error)
	GetNameById(userId int) (string, error)
	CheckAuthorComment(commentId int, userId int) (bool, error)
	GetConversationMessages(cnvId int) ([]Message, error)
	CreateImageMessage(image media.Image, conversationID, userID int) (Message, int, error)
	GetMessagePage(conversationID int, before *MessagePosition, limit int) (MessagePage, error)
	GetMessageChanges(conversationID int, since int64, limit int) (MessagePage, error)
	GetImage(id int) (media.Image, error)
	CanReadImage(id, userID int) (bool, error)
	ImageExists(id int) (bool, error)
	GetUserPhoto(userId int) (string, error)
	CheckUsernameExist(username string) (bool, error)
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// One connection avoids SQLite write contention and keeps in-memory databases coherent.
	// Open additionally configures foreign keys on every connection the pool may replace.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	if _, err := db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		return nil, err
	}
	var enabled int
	if err := db.QueryRow("PRAGMA foreign_keys").Scan(&enabled); err != nil {
		return nil, err
	}
	if enabled != 1 {
		return nil, errors.New("SQLite foreign keys must be enabled")
	}
	if err := migrate(db, migrationFiles); err != nil {
		return nil, fmt.Errorf("initialize V2 schema: %w", err)
	}

	return &appdbimpl{c: db}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
