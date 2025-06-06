package database

import (
	"fmt"
	"strconv"
	"strings"
)

func (db *appdbimpl) RemoveMemberFromGroup(userId int, groupId int) error {
	// prendo la lista di membri della conversazione
	var currentMembersString string
	err := db.c.QueryRow("SELECT members FROM conversations_table WHERE conversation_id=?", groupId).Scan(&currentMembersString)
	if err != nil {
		return fmt.Errorf("errore durante la selezione dei membri del gruppo: %w", err)
	}

	// creo la lista e salto l'userId
	var members []int
	for _, id := range strings.Split(currentMembersString, ",") {
		num, err := strconv.Atoi(id)
		if err != nil {
			return fmt.Errorf("id non convertibile: %w", err)
		}
		if num != userId {
			members = append(members, num)
		}
	}

	// ritrasformo in stringa
	var parts []string
	for _, id := range members {
		parts = append(parts, strconv.Itoa(id))
	}
	updatedMembers := strings.Join(parts, ",")

	// aggiorno nel database
	_, err = db.c.Exec("UPDATE conversations_table SET members=? WHERE conversation_id=?", updatedMembers, groupId)
	if err != nil {
		return fmt.Errorf("errore durante l'aggiornamente dei membri: %w", err)
	}

	return nil
}
