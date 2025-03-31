package database

func (db *appdbimpl) GetMessage(messageId int) (Message, error) {
	//prendo il messaggio dalla tabella
	var message Message
	err := db.c.QueryRow("SELECT * FROM messages_table WHERE message_id=?", messageId).Scan(&message)
	if err != nil {
		return Message{}, err
	}
	return message, nil
}
