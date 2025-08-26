/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

var ErrNoConversation = errors.New("non ha conversazioni")

var DefaultGroupPhoto = "iVBORw0KGgoAAAANSUhEUgAABP8AAAUACAYAAAA4LfcjAAEAAElEQVR42uzdeXzV1Z3/8c855/u9CUurpdVRW621trXWpbIIam21tasz03V+WqtiWEMI4NbO1HamtTOd1gWUPYRERJQdVFqlVVuX4oIg2qptLXZ0RAghITd3377nnN8f5DK4g2zJzev5ePBAIQTyud/7vZ/7zuecowQAAKB3Ud57+elPf6o+9alPKRGR559/XomIfOpTn..."
var DefaultContactPhoto = "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"

// Faccio i diversi schemas
type User struct {
	Id    int    `json:"userId"`
	Name  string `json:"name"`
	Photo string `json:"photo"`
}
type Conversation struct {
	ConversationId int    `json:"conversationId"`
	Title          string `json:"title"`
	Photo          string `json:"photo"`
	IsGroup        bool   `json:"isGroup"`
	Members        []User `json:"members"`
}

type Message struct {
	MessageId int       `json:"messageId"`
	AuthorId  int       `json:"authorId"`
	Author    string    `json:"author"`
	Text      string    `json:"text"`
	Hours     string    `json:"hours"`
	Comments  []Comment `json:"comments"`
}

type Comment struct {
	CommentId int    `json:"commentId"`
	AuthorId  int    `json:"authorId"`
	Emoji     string `json:"emoji"`
}

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetIdByName(name string) (int, error)
	CreateUser(name string) (int, error)
	SetName(name string, id int) error
	GetConversations(id int) ([]Conversation, error)
	GetSpecificConversation(conversationId int, currentUserId int) (Conversation, error)
	CreateMessage(author string, text string, hours string, cnvid int, userid int) (Message, int, error)
	ConversationExists(cnvId int) (bool, error)
	CommentExists(commentId int) (bool, error)
	MessageExists(messageId int) (bool, error)
	CheckUser_Conversation(cnvId int, userId int) (bool, error)
	CheckMessage_Conversation(messageId int, cnvId int) (bool, error)
	GetMessage(messageId int) (Message, error)
	CreateComment(author int, emoji string, messageId int) (Comment, int, error)
	CheckComment_Message(commentId int, messageId int) (bool, error)
	DeleteComment(commentId int) error
	CheckUserId(id int) (bool, error)
	CheckMessage_User(messageId int, userId int) (bool, error)
	RemoveMessage(messageId int) error
	AddUserToGroup(userId int, cnvId int) error
	Conversation_isGroup(conversationId int) (bool, error)
	GetConversationMembers(cnvId int) ([]int, error)
	GetUsersMembers(membri []int) ([]User, error)
	RemoveGroupFromMember(userId int, groupId int) error
	RemoveMemberFromGroup(userId int, groupId int) error
	SetGroupName(groupName string, cnvId int) error
	SetMyPhoto(photo string, id int) error
	SetGroupPhoto(photo string, cnvId int) error
	GetAllUsers(userId int) ([]int, []string, error)
	CreateConversation(userId1 int, userId2 int) (int, error)
	CreateGroup(titolo string, members []int, userId int) (int, error)
	Conversation2MembersExists(userId1 int, userId2 int) (bool, int, error)
	GetNameById(userId int) (string, error)
	CheckAuthorComment(commentId int, userId int) (bool, error)
	GetConversationMessages(cnvId int) ([]Message, error)
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

	var tableName string
	// TABELLA UTENTI
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE users_table (	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
												name TEXT,
												photo TEXT,
												conversation_ids TEXT);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}
	// TABELLA CONVERSAZIONI
	tableName = ""
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='conversations_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE conversations_table (	conversation_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
															titolo TEXT,
															photo TEXT, 
															isGroup INTEGER NOT NULL CHECK(isGroup IN (0,1)),
															members TEXT);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}
	// TABELLA MESSAGGIxCONVERSAZIONE
	tableName = ""
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='messages_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE messages_table (	message_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
														author TEXT,
														text TEXT, 
														hours TEXT,
														comments TEXT,
														conversation_id INTEGER NOT NULL,
														user_id INTEGER NOT NULL,
														FOREIGN KEY (conversation_id) REFERENCES conversations_table(conversation_id) ON DELETE CASCADE,
														FOREIGN KEY (user_id) REFERENCES users_table(id));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}
	// TABELLA COMMENTIxMESSAGGI
	tableName = ""
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='comments_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE comments_table (	comment_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
														author_id INTEGER NOT NULL,
														emoji TEXT, 
														message_id INTEGER NOT NULL,
														FOREIGN KEY (author_id) REFERENCES users_table(id),
														FOREIGN KEY (message_id) REFERENCES messages_table(message_id) ON DELETE CASCADE);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}
	return &appdbimpl{c: db}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

// --------------------------------------------------------------------------------------------------
