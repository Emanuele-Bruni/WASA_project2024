package database

// SetName is an example that shows you how to execute insert/update
func (db *appdbimpl) SetName(name string, id int) error {
	_, err := db.c.Exec("UPDATE users_table SET name=? WHERE id=?", name, id)
	if err != nil {
		return err
	}

	// aggiorno il nome nella tabella dei messaggi
	_, err = db.c.Exec("UPDATE messages_table SET author=? WHERE user_id=?", name, id)
	if err != nil {
		return err
	}

	return nil
}
