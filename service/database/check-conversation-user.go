package database

func (db *appdbimpl) CheckUser_Conversation(conversationID, userID int) (bool, error) {
	var member bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM conversation_members WHERE conversation_id=? AND user_id=?)", conversationID, userID).Scan(&member)
	return member, err
}
