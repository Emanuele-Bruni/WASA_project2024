package database

import (
	"fmt"
	"strconv"
	"strings"
)

// seleziono la conversazione dalla tabella delle conversazioni
func (db *appdbimpl) GetSpecificConversation(conversationId int) (Conversation, error) {
	var conversation Conversation
	var members string
	err := db.c.QueryRow("SELECT conversation_id,titolo,photo,isGroup,members FROM conversations_table WHERE conversation_id=?", conversationId).Scan(&conversation.ConversationId, &conversation.Title, &conversation.Photo, &conversation.IsGroup, &members)
	if err != nil {
		return Conversation{}, fmt.Errorf("errore nella selezione dei dati della conversazione: %w", err)
	}

	var membersIdList []int
	var finalMembers []User
	// devo creare una lista di utenti
	parti := strings.Split(members, ",")
	for _, id_stringa := range parti {
		id, err := strconv.Atoi(id_stringa)
		if err != nil {
			return Conversation{}, fmt.Errorf("errore nella conversione del memberId: %w", err)
		}
		membersIdList = append(membersIdList, id)
	}

	// ora prendo l'id e ritorno l'utente, aggiungendolo ai membri finali
	for _, id := range membersIdList {
		var utente User
		utente.Id = id
		err = db.c.QueryRow("SELECT name,photo FROM users_table WHERE id=?", id).Scan(&utente.Name, &utente.Photo)
		finalMembers = append(finalMembers, utente)
	}
	conversation.Members = finalMembers
	return conversation, err
}
