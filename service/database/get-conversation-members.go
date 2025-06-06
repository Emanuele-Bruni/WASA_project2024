package database

import (
	"fmt"
	"strconv"
	"strings"
)

func (db *appdbimpl) GetConversationMembers(cnvId int) ([]int, error) {
	// prendo l'array dei membri e lo trasformo in lista
	var stringMembers string // "20,30,40"
	err := db.c.QueryRow("SELECT members FROM conversations_table WHERE conversation_id=?", cnvId).Scan(&stringMembers)
	if err != nil {
		return []int{}, fmt.Errorf("errore nel selezionare i membri della conversazione: %w", err)
	}
	// converto in lista
	var listMembers []int
	for _, id_stringa := range strings.Split(stringMembers, ",") {
		// converto in numero
		num, err := strconv.Atoi(id_stringa)
		if err != nil {
			return []int{}, fmt.Errorf("id non convertibile: %w", err)
		} else {
			listMembers = append(listMembers, num)
		}
	}

	return listMembers, nil
}
