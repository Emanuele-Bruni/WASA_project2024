package database

// GetName is an example that shows you how to query data
func (db *appdbimpl) GetIdByName(name string) (int, error) {
	var id int
	err := db.c.QueryRow("SELECT id FROM users_table WHERE name=?", name).Scan(&id)
	return id, err
}
