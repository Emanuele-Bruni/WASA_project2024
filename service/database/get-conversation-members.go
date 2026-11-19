package database

func (db *appdbimpl) GetConversationMembers(conversationID int) ([]User, error) {
	rows, err := db.c.Query(`SELECT u.id,u.name,CASE WHEN u.photo_image_id IS NULL THEN '' ELSE '/media/'||u.photo_image_id END FROM conversation_members cm
 JOIN users_table u ON u.id=cm.user_id WHERE cm.conversation_id=? ORDER BY u.id`, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]User, 0)
	for rows.Next() {
		var user User
		if err = rows.Scan(&user.Id, &user.Name, &user.Photo); err != nil {
			return nil, err
		}
		result = append(result, user)
	}
	return result, rows.Err()
}
