package database

// Message author names are derived from the user relation; no duplicated names need updating.
func (db *appdbimpl) SetName(name string, id int) error {
	return requireAffected(db.c.Exec("UPDATE users_table SET name=? WHERE id=?", name, id))
}
