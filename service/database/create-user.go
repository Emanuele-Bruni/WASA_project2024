package database

import "fmt"

// CreateUser crea un nuovo utente con il nome specificato e restituisce il suo ID.
func (db *appdbimpl) CreateUser(name string) (int, error) {
	res, err := db.c.Exec("INSERT INTO users_table (name) VALUES (?)", name)
	if err != nil {
		return 0, fmt.Errorf("error creating user: %w", err)
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert ID: %w", err)
	}
	return int(lastID), nil
}
