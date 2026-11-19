package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) Conversation2MembersExists(first, second int) (bool, int, error) {
	if first == second {
		return false, 0, fmt.Errorf("conversation with self is not allowed")
	}
	if first > second {
		first, second = second, first
	}
	var id int
	err := db.c.QueryRow("SELECT conversation_id FROM conversations_table WHERE is_group=0 AND direct_user_low=? AND direct_user_high=?", first, second).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return false, 0, nil
	}
	return err == nil, id, err
}
