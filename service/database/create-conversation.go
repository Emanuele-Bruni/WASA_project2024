package database

import (
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

	return int(conversationId), nil
}
