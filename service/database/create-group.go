package database

import (
	"fmt"
	"strconv"
	"strings"
)

func (db *appdbimpl) CreateGroup(titolo string, members []int, userId int) (int, error) {
	//devo trasformare members da lista a stringa
	members = append(members, userId)
	strMembers := make([]string, len(members))

	for i, num := range members {
		strMembers[i] = strconv.Itoa(num)
	}

	partecipants := strings.Join(strMembers, ",")
	photo := DefaultGroupPhoto
	isGroup := 1

	res, err := db.c.Exec("INSERT INTO conversations_table (titolo,photo,isGroup,members) VALUES (?,?,?,?)", titolo, photo, isGroup, partecipants)
	if err != nil {
		return -1, fmt.Errorf("error during group creation: %w", err)
	}

	groupId, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("error getting last insert ID: %w", err)
	}

	return int(groupId), nil

}
