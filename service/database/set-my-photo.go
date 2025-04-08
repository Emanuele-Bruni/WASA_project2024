package database

func (db *appdbimpl) SetMyPhoto(photo string, id int) error {
	_, err := db.c.Exec("UPDATE users_table SET photo=? WHERE id=?", photo, id)
	return err
}
