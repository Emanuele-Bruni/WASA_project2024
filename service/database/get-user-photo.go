package database

func (db *appdbimpl) GetUserPhoto(userId int) (string, error) {
	var photo string
	err := db.c.QueryRow("SELECT photo FROM users_table WHERE id=?", userId).Scan(&photo)
	return photo, err
}
