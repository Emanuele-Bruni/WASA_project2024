package database

import (
	"database/sql"
	"fmt"
	"strconv"
)

func (db *appdbimpl) AddUserToGroup(userId int, cnvId int) error {
	// devo prendere la lista dei membri della conversazione
	var currentStringMembers string // "20,30,40"
	err := db.c.QueryRow("SELECT members FROM conversations_table WHERE conversation_id=?", cnvId).Scan(&currentStringMembers)
	if err != nil {
		return fmt.Errorf("errore nel selezionare i membri della conversazione: %w", err)
	}
	var updatedMembers string
	if currentStringMembers != "" {
		updatedMembers = currentStringMembers + "," + strconv.Itoa(userId)
	} else {
		updatedMembers = strconv.Itoa(userId)
	}
	// e lo aggiorno nella tabella conversations_table
	_, err = db.c.Exec("UPDATE conversations_table SET members=? WHERE conversation_id=?", updatedMembers, cnvId)
	if err != nil {
		return fmt.Errorf("errore nell'aggiornare i membri della conversazione: %w", err)
	}

	// devo aggiungere l'id della conversazione a conversation_ids della tabella users_table
	// prendo l'array delle conversazioni
	var currentStringConversations sql.NullString
	err = db.c.QueryRow("SELECT conversation_ids FROM users_table WHERE id=?", userId).Scan(&currentStringConversations)
	if err != nil {
		return fmt.Errorf("errore nel selezionare la lista di conversazioni dell'utente: %w", err)
	}
	var updatedConversations string
	if currentStringConversations.Valid && currentStringConversations.String != "" {
		updatedConversations = currentStringConversations.String + "," + strconv.Itoa(cnvId)
	} else {
		updatedConversations = strconv.Itoa(cnvId)
	}
	// aggiorno users_table
	_, err = db.c.Exec("UPDATE users_table SET conversation_ids=? WHERE id=?", updatedConversations, userId)
	if err != nil {
		return fmt.Errorf("errore nell'aggiornare le conversazioni dell'utente: %w", err)
	}

	return nil
}
