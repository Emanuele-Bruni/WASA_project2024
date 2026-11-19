package database

import (
	"database/sql"
	"fmt"
)

func (db *appdbimpl) CreateConversation(userID1, userID2 int) (int, error) {
	if userID1 == userID2 {
		return 0, fmt.Errorf("conversation with self is not allowed")
	}
	if userID1 > userID2 {
		userID1, userID2 = userID2, userID1
	}
	return db.createConversation("", []int{userID1, userID2}, false)
}

// Creation of the conversation and all memberships is a single atomic write.
func (db *appdbimpl) createConversation(title string, members []int, isGroup bool) (int, error) {
	tx, err := db.c.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	var low, high interface{}
	if !isGroup {
		low, high = members[0], members[1]
	}
	avatar, err := defaultAvatar(isGroup)
	if err != nil {
		return 0, err
	}
	imageID, err := saveImage(tx, avatar)
	if err != nil {
		return 0, err
	}
	result, err := tx.Exec("INSERT INTO conversations_table(title,photo,is_group,direct_user_low,direct_user_high,photo_image_id) VALUES(?,'',?,?,?,?)", title, isGroup, low, high, imageID)
	if err != nil {
		return 0, fmt.Errorf("create conversation: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	for _, userID := range members {
		if _, err = tx.Exec("INSERT INTO conversation_members(conversation_id,user_id) VALUES(?,?)", id, userID); err != nil {
			return 0, fmt.Errorf("add member: %w", err)
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return int(id), nil
}

func requireAffected(result sql.Result, err error) error {
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}
