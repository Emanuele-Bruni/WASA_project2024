package database

import "fmt"

// devo ritornare [lista ids],[lista nomi],errore
func (db *appdbimpl) GetAllUsers(userId int) ([]int, []string, error) {

	righe, err := db.c.Query("SELECT id,name FROM users_table WHERE id!=?", userId)
	if err != nil {
		return nil, nil, fmt.Errorf("errore nella selezione degli utenti: %w", err)
	}
	defer righe.Close()

	var ids []int
	var names []string

	//prendo gli id degli utenti
	for righe.Next() {
		var id int
		var name string
		//prendo l'id e il nome dalla tupla
		err := righe.Scan(&id, &name)
		if err != nil {
			return nil, nil, fmt.Errorf("errore durante la scansione: %w", err)
		}
		//li aggiungo
		ids = append(ids, id)
		names = append(names, name)
	}

	return ids, names, nil
}
