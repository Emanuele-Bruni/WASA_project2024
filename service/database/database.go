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

var ErrorNoConversation = errors.New("utente non ha conversazioni")

// Faccio i diversi schemas
type User struct {
	Name   string `json:"name"`
	Photo  string `json:"photo"`
	Bio    string `json:"bio"`
	Number string `json:"number"`
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
	Author    string    `json:"author"`
	Text      string    `json:"text"`
	Hours     string    `json:"hours"`
	Comments  []Comment `json:"comments"`
}

type Comment struct {
	CommentId int    `json:"commentId"`
	Author    string `json:"author"`
	Emoji     string `json:"emoji"`
}

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetIdByName(name string) (int, error)
	CreateUser(name string) (int, error)
	SetName(name string, id int) error
	GetConversations(id int) ([]Conversation, error)
	GetSpecificConversation(conversationId int) (Conversation, error)
	CreateMessage(author string, text string, hours string, cnvid int, userid int) (Message, error)
	ConversationExists(cnvId int) (bool, error)
	CommentExists(commentId int) (bool, error)
	MessageExists(messageId int) (bool, error)
	CheckUser_Conversation(cnvId int, userId int) (bool, error)
	CheckMessage_Conversation(messageId int, cnvId int) (bool, error)
	GetMessage(messageId int) (Message, error)
	CreateComment(author string, emoji string, messageId int) (Comment, error)
	CheckComment_Message(commentId int, messageId int) (bool, error)
	DeleteComment(commentId int) error
	CheckUserId(id int) (bool, error)
	CheckMessage_User(messageId int, userId int) (bool, error)
	RemoveMessage(messageId int) error
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
	//TABELLA UTENTI
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE users_table (	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
												name TEXT,
												photo TEXT, 
												number TEXT,
												conversation_ids TEXT);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}
	//TABELLA CONVERSAZIONI
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
	//TABELLA MESSAGGIxCONVERSAZIONE
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
	//TABELLA COMMENTIxMESSAGGI
	tableName = ""
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='comments_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE comments_table (	comment_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
														author TEXT,
														emoji TEXT, 
														message_id INTEGER NOT NULL,
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

//--------------------------------------------------------------------------------------------------
