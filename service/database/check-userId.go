package database

//controlla se l'id esiste
func (db *appdbimpl) CheckUserId(id int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users_table WHERE id=?)"
	err := db.c.QueryRow(query, id).Scan(&exists)
	return exists, err
}
