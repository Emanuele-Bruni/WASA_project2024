package database

func (db *appdbimpl) SetGroupName(groupName string, cnvId int) error {
	_, err := db.c.Exec("UPDATE conversations_table SET titolo=? WHERE conversation_id=?", groupName, cnvId)
	return err
}
