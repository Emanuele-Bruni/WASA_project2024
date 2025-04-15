package database

import "fmt"

// CreateMessage crea un nuovo messaggio e lo restituisce (ovviamente senza commenti, perché è stato appena creato).
func (db *appdbimpl) CreateMessage(author string, text string, hours string, cnvid int, userid int) (Message, int, error) {
	//inserisco i dati in input nella tabella
	res, err := db.c.Exec("INSERT INTO messages_table (author,text,hours,conversation_id,user_id) VALUES (?,?,?,?,?)", author, text, hours, cnvid, userid)
	if err != nil {
		return Message{}, -1, fmt.Errorf("error creating message: %w", err)
	}

	//calcolo l'ID
	messageId, err := res.LastInsertId()
	if err != nil {
		return Message{}, -1, fmt.Errorf("error getting last insert ID: %w", err)
	}

	//determino il messaggio da restituire (con commenti vuoti)
	finalMessage := Message{
		MessageId: int(messageId),
		Author:    author,
		Text:      text,
		Hours:     hours,
		Comments:  []Comment{},
	}
	return finalMessage, int(messageId), nil
}
