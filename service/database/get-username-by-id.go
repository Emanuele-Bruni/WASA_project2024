package database

func (db *appdbimpl) GetNameById(userId int) (string, error) {
	var name string
	err := db.c.QueryRow("SELECT name FROM users_table WHERE id=?", userId).Scan(&name)
	return name, err
}
