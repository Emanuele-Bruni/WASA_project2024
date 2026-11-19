package database

func (db *appdbimpl) GetUserPhoto(userId int) (string, error) {
	var photo string
	err := db.c.QueryRow("SELECT CASE WHEN photo_image_id IS NULL THEN '' ELSE '/media/'||photo_image_id END FROM users_table WHERE id=?", userId).Scan(&photo)
	return photo, err
}
