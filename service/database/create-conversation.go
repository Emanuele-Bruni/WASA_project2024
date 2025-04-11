package database

import (
	"database/sql"
	"fmt"
	"strconv"
)

func (db *appdbimpl) CreateConversation(userId1 int, userId2 int) (int, error) {
	//prendo le info dell'utente2
	//titolo -> username dell'utente2
	//photo -> photo standard
	//isGroup -> 0
	//members -> "userId1,userId2"
	var username string
	photo := DefaultContactPhoto
	isGroup := 0
	members := strconv.Itoa(userId1) + "," + strconv.Itoa(userId2)

	err := db.c.QueryRow("SELECT name FROM users_table WHERE id=?", userId2).Scan(&username)
	if err != nil {
		return -1, fmt.Errorf("error during user selection: %w", err)
	}

	res, err := db.c.Exec("INSERT INTO conversations_table (titolo,photo,isGroup,members) VALUES (?,?,?,?)", username, photo, isGroup, members)
	if err != nil {
		return -1, fmt.Errorf("error during conversation creation: %w", err)
	}

	conversationId, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("error getting last insert ID: %w", err)
	}

	//devo inserirlo in conversation_ids nella users_table di entrambi gli utenti
	//UTENTE 1
	var currentConversationsString1 sql.NullString
	err = db.c.QueryRow("SELECT conversation_ids FROM users_table WHERE id=?", userId1).Scan(&currentConversationsString1)
	if err != nil {
		return -1, fmt.Errorf("error during conversations selection: %w", err)
	}
	var updatedConversations1 string
	//controllo se è non NULL e se non è vuota
	if currentConversationsString1.Valid && currentConversationsString1.String != "" {
		updatedConversations1 = fmt.Sprintf("%s,%d", currentConversationsString1.String, conversationId)
	} else {
		updatedConversations1 = fmt.Sprintf("%d", conversationId)
	}
	_, err = db.c.Exec("UPDATE users_table SET conversation_ids=? WHERE id=?", updatedConversations1, userId1)
	if err != nil {
		return -1, fmt.Errorf("error insert new conversation into the database: %w", err)
	}

	//UTENTE 2
	var currentConversationsString2 sql.NullString
	err = db.c.QueryRow("SELECT conversation_ids FROM users_table WHERE id=?", userId2).Scan(&currentConversationsString2)
	if err != nil {
		return -1, fmt.Errorf("error during conversations selection: %w", err)
	}
	var updatedConversations2 string
	//controllo se è non NULL e se non è vuota
	if currentConversationsString2.Valid && currentConversationsString2.String != "" {
		updatedConversations2 = fmt.Sprintf("%s,%d", currentConversationsString2.String, conversationId)
	} else {
		updatedConversations2 = fmt.Sprintf("%d", conversationId)
	}
	_, err = db.c.Exec("UPDATE users_table SET conversation_ids=? WHERE id=?", updatedConversations2, userId2)
	if err != nil {
		return -1, fmt.Errorf("error insert new conversation into the database: %w", err)
	}

	return int(conversationId), nil
}
