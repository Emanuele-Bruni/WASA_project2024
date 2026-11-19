package database

// LeaveGroup removes the membership and, for the last participant, the empty group.
// Foreign keys also remove the group's messages and reactions in the same transaction.
func (db *appdbimpl) LeaveGroup(userID, conversationID int) error {
	tx, err := db.c.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err = requireAffected(tx.Exec(`DELETE FROM conversation_members
 WHERE conversation_id=? AND user_id=? AND EXISTS(
 SELECT 1 FROM conversations_table WHERE conversation_id=? AND is_group=1)`, conversationID, userID, conversationID)); err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM conversations_table WHERE conversation_id=? AND is_group=1
 AND NOT EXISTS(SELECT 1 FROM conversation_members WHERE conversation_id=?)`, conversationID, conversationID); err != nil {
		return err
	}
	return tx.Commit()
}
