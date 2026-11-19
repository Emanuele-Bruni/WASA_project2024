package database

func (db *appdbimpl) SetGroupName(groupName string, cnvId int) error {
	return requireAffected(db.c.Exec("UPDATE conversations_table SET title=? WHERE conversation_id=?", groupName, cnvId))
}
