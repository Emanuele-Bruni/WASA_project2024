package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

func (db *appdbimpl) CreateGroup(titolo string, members []int, userId int) (int, error) {
	//titolo -> titolo
	//photo -> default
	//isGroup -> 1
	//membri "1,2,3,4"

	//controllo di non star aggiungendo me stesso
	for _, num := range members {
		if num == userId {
			return -1, fmt.Errorf("impossible to add yourself into a group")
		}
	}

	//aggiungo il creatore del gruppo ai membri
	totalmembers := append(members, userId)

	//devo trasformare totalmembers da lista di interi a stringa
	strMembers := make([]string, len(totalmembers))
	for i, num := range totalmembers {
		strMembers[i] = strconv.Itoa(num)
	}
	partecipants := strings.Join(strMembers, ",")

	groupTitle := titolo
	photo := DefaultGroupPhoto
	isGroup := 1

	res, err := db.c.Exec("INSERT INTO conversations_table (titolo,photo,isGroup,members) VALUES (?,?,?,?)", groupTitle, photo, isGroup, partecipants)
	if err != nil {
		return -1, fmt.Errorf("error during group creation: %w", err)
	}

	groupId, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("error getting last insert ID: %w", err)
	}

	//devo inserire il groupId alla lista delle conversazioni dell'utente per ogni membro
	for _, memberId := range totalmembers {
		var conversazioniUtente sql.NullString
		err = db.c.QueryRow("SELECT conversation_ids FROM users_table WHERE id=?", memberId).Scan(&conversazioniUtente)
		if err != nil {
			return -1, fmt.Errorf("error during select users conversation_ids: %w", err)
		}
		var updatedConversazioniUtente string
		//controllo se è non NULL e se non è vuota
		if conversazioniUtente.Valid && conversazioniUtente.String != "" {
			updatedConversazioniUtente = fmt.Sprintf("%s,%d", conversazioniUtente.String, groupId)
		} else {
			updatedConversazioniUtente = fmt.Sprintf("%d", groupId)
		}

		_, err = db.c.Exec("UPDATE users_table SET conversation_ids=? WHERE id=?", updatedConversazioniUtente, memberId)
		if err != nil {
			return -1, fmt.Errorf("error during update users conversation_ids: %w", err)
		}
	}

	return int(groupId), nil
}
