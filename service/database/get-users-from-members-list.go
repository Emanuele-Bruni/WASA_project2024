package database

import "fmt"

func (db *appdbimpl) GetUsersMembers(membri []int) ([]User, error) {
	if len(membri) == 0 {
		return []User{}, nil
	}

	var utenti []User

	for _, id := range membri {
		var u User
		err := db.c.QueryRow("SELECT id, name, photo FROM users_table WHERE id = ?", id).Scan(&u.Id, &u.Name, &u.Photo)
		if err != nil {
			return nil, fmt.Errorf("errore nel recuperare utente con id %d: %w", id, err)
		}
		utenti = append(utenti, u)
	}

	return utenti, nil
}
