package database

// SetName is an example that shows you how to execute insert/update
func (db *appdbimpl) SetName(name string, id int) error {
	_, err := db.c.Exec("UPDATE users_table SET name=? WHERE id=?", name, id)
	return err
}
