package database

import (
	"fmt"
	"strconv"
	"strings"
)

func (db *appdbimpl) AddUserToGroup(userId int, cnvId int) error {
	//devo prendere la lista dei membri della conversazione
	var currentStringMembers string // "20,30,40"
	err := db.c.QueryRow("SELECT members FROM conversations_table WHERE conversation_id=?", cnvId).Scan(&currentStringMembers)
	if err != nil {
		return fmt.Errorf("errore nel selezionare i membri della conversazione: %w", err)
	}
	//converto in lista
	var currentListMembers []int
	for _, id_stringa := range strings.Split(currentStringMembers, ",") {
		//converto in numero
		num, err := strconv.Atoi(id_stringa)
		if err != nil {
			return fmt.Errorf("id non convertibile: %w", err)
		} else {
			currentListMembers = append(currentListMembers, num)
		}
	}
	//aggiungo l'userId
	currentListMembers = append(currentListMembers, userId)
	// la rendo nuovamente un array
	updatedMembers := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(currentListMembers)), ","), "[]")
	// e lo aggiorno nella tabella conversations_table
	_, err = db.c.Exec("UPDATE conversations_table SET members=? WHERE conversation_id=?", updatedMembers, cnvId)
	if err != nil {
		return fmt.Errorf("errore nell'aggiornare i membri della conversazione: %w", err)
	}

	//devo aggiungere l'id della conversazione a conversation_ids della tabella users_table
	//prendo l'array delle conversazioni
	var currentStringConversations string
	err = db.c.QueryRow("SELECT conversation_ids FROM users_table WHERE user_id=?", userId).Scan(&currentStringConversations)
	if err != nil {
		return fmt.Errorf("errore nel selezionare la lista di conversazioni dell'utente: %w", err)
	}
	//lo rendo lista
	var currentsListConversations []int
	for _, cnvId_stringa := range strings.Split(currentStringConversations, ",") {
		//converto in numero
		num, err := strconv.Atoi(cnvId_stringa)
		if err != nil {
			return fmt.Errorf("errore nella conversione del conversationId della tabella in numero: %w", err)
		} else {
			currentsListConversations = append(currentsListConversations, num)
		}
	}
	//aggiungo il cnvId alle conversazioni dell'utente
	currentsListConversations = append(currentsListConversations, cnvId)
	// lo converto in array
	updatedConversations := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(currentsListConversations)), ","), "[]")
	//aggiorno users_table
	_, err = db.c.Exec("UPDATE users_table SET conversation_ids=? WHERE user_id=?", updatedConversations, userId)
	if err != nil {
		return fmt.Errorf("errore nell'aggiornare le conversazioni dell'utente: %w", err)
	}
	return nil
}
