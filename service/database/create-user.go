package database

// GetOrCreateUser handles concurrent logins without creating duplicate usernames.
func (db *appdbimpl) GetOrCreateUser(name string) (int, error) {
	tx, err := db.c.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	avatar, err := defaultAvatar(false)
	if err != nil {
		return 0, err
	}
	imageID, err := saveImage(tx, avatar)
	if err != nil {
		return 0, err
	}
	if _, err = tx.Exec("INSERT INTO users_table(name,photo,photo_image_id) VALUES(?,'',?) ON CONFLICT(name) DO NOTHING", name, imageID); err != nil {
		return 0, err
	}
	var id int
	if err = tx.QueryRow("SELECT id FROM users_table WHERE name=?", name).Scan(&id); err != nil {
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return id, nil
}
