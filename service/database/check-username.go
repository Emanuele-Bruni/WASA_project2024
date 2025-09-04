package database

func (db *appdbimpl) CheckUsernameExist(username string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users_table WHERE name=?)"
	err := db.c.QueryRow(query, username).Scan(&exists)
	return exists, err
}
