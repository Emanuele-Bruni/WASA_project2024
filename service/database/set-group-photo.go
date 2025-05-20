package database

func (db *appdbimpl) SetGroupPhoto(photo string, cnvId int) error {
	_, err := db.c.Exec("UPDATE conversations_table SET photo=? WHERE conversation_id=?", photo, cnvId)
	return err
}
