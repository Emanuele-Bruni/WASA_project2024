package database

import "fmt"

// CreateUser crea un nuovo utente con il nome specificato e restituisce il suo ID.
func (db *appdbimpl) CreateUser(name string) (int, error) {
	photo := DefaultContactPhoto
	res, err := db.c.Exec("INSERT INTO users_table (name,photo,bio,number) VALUES (?,?,?,?)", name, photo, "", "")
	if err != nil {
		return -1, fmt.Errorf("error creating user: %w", err)
	}
	userID, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("error getting last insert ID: %w", err)
	}
	return int(userID), nil
}
