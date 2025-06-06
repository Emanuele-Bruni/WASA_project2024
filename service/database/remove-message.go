package database

import "fmt"

func (db *appdbimpl) RemoveMessage(messageId int) error {
	_, err := db.c.Exec("DELETE FROM messages_table WHERE message_id=?", messageId)
	if err != nil {
		return fmt.Errorf("errore nell'eliminare il messaggio: %w", err)
	}
	return nil
}
