package database

import "fmt"

func (db *appdbimpl) CreateGroup(title string, members []int, userID int) (int, error) {
	if len(members) < 2 || len(members) > 16 {
		return 0, fmt.Errorf("invalid member count")
	}
	all := make([]int, 0, len(members)+1)
	all = append(all, userID)
	all = append(all, members...)
	return db.createConversation(title, all, true)
}
