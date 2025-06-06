package database

import (
	"fmt"
	"strconv"
	"strings"
)

func (db *appdbimpl) RemoveGroupFromMember(userId int, groupId int) error {
	// prendo la stringa delle conversazioni di cui fa parte l'utente
	var conversationsString string
	err := db.c.QueryRow("SELECT conversation_ids FROM users_table WHERE id=?", userId).Scan(&conversationsString)
	if err != nil {
		return fmt.Errorf("errore durante la selezione delle conversazioni dell'utente: %w", err)
	}

	// ora la rendo lista e non aggiungo se la conversazione è il gruppo da cui dobbiamo togliere l'utente
	var conversations []int
	for _, id := range strings.Split(conversationsString, ",") {
		num, err := strconv.Atoi(id)
		if err != nil {
			return fmt.Errorf("id non convertibile: %w", err)
		}
		if num != groupId {
			conversations = append(conversations, num)
		}
	}

	// ritrasformo in stringa
	var parts []string
	for _, id := range conversations {
		parts = append(parts, strconv.Itoa(id))
	}
	updatedConversations := strings.Join(parts, ",")

	// aggiorno nel database
	_, err = db.c.Exec("UPDATE users_table SET conversation_ids = ? WHERE id=?", updatedConversations, userId)
	if err != nil {
		return fmt.Errorf("errore durante l'aggiornamento delle conversazioni utente: %w", err)
	}

	return nil
}
