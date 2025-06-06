package database

import (
	"strconv"
	"strings"
)

func (db *appdbimpl) CheckUser_Conversation(cnvId int, userId int) (bool, error) {
	// prendo la conversazioni dalla tabella e prendo i membri
	var members string
	err := db.c.QueryRow("SELECT members FROM conversations_table WHERE conversation_id=?", cnvId).Scan(&members)
	if err != nil {
		return false, err
	}
	// controllo se l'id è tra i membri
	ids := strings.Split(members, ",")
	for _, val := range ids {
		val1, err := strconv.Atoi(val)
		if err != nil {
			return false, err
		}
		if userId == val1 {
			return true, nil
		}
	}
	// non è tra i membri
	return false, nil
}
