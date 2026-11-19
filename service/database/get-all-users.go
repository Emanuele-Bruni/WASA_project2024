package database

func (db *appdbimpl) GetUsers(userID int) ([]User, error) {
	rows, err := db.c.Query("SELECT id,name,CASE WHEN photo_image_id IS NULL THEN '' ELSE '/media/'||photo_image_id END FROM users_table WHERE id!=? ORDER BY id", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]User, 0)
	for rows.Next() {
		var user User
		if err = rows.Scan(&user.Id, &user.Name, &user.Photo); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}
func (db *appdbimpl) GetUser(userID int) (User, error) {
	var user User
	err := db.c.QueryRow("SELECT id,name,CASE WHEN photo_image_id IS NULL THEN '' ELSE '/media/'||photo_image_id END FROM users_table WHERE id=?", userID).Scan(&user.Id, &user.Name, &user.Photo)
	return user, err
}
