package database

import (
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

func (db *appdbimpl) GetConversations(userId int) ([]Conversation, error) {

	var userConversations []Conversation //lista finale delle conversazioni
	var stringaConversazioni string      //stringa "20,30,40"
	var lista_id []int                   // [20,30,40]

	//controllo l'id nella tabella utenti e prendo la stringa di conversazioni
	err := db.c.QueryRow("SELECT conversation_ids FROM users_table WHERE id=?", userId).Scan(&stringaConversazioni)
	if err != nil {
		return []Conversation{}, err
	}
	if stringaConversazioni == "" {
		return []Conversation{}, ErrorNoConversation
	}

	//prendo gli id delle conversazioni dell'utente e la aggiungo alla lista degli id
	for _, id_stringa := range strings.Split(stringaConversazioni, ",") {
		//converto il valore in numero
		num, err := strconv.Atoi(id_stringa)
		if err != nil {
			logrus.WithField("conversation_id", id_stringa).WithError(err).Error("id non valido")
		} else {
			lista_id = append(lista_id, num)
		}
	}

	for _, conversation_id := range lista_id {
		var conversation Conversation
		var members string
		var isGroup int

		//recupero le info della conversazione
		err := db.c.QueryRow("SELECT titolo,photo,isGroup,members FROM conversations_table WHERE conversation_id=?", conversation_id).Scan(&conversation.Title, &conversation.Photo, &isGroup, &members)
		if err != nil {
			return []Conversation{}, err
		}
		if isGroup == 0 {
			conversation.IsGroup = false
		} else {
			conversation.IsGroup = true
		}

		//controllo i membri
		for _, memberId_stringa := range strings.Split(members, ",") {
			//converto membro in numero
			memberId, err := strconv.Atoi(memberId_stringa)
			var utente User
			if err != nil {
				logrus.WithField("memberId not valid", memberId).WithError(err).Error("id non valido")
			} else {
				//ottenuto id membro, ora devo prendere le sue info
				err := db.c.QueryRow("SELECT name,photo,number FROM users_table WHERE id=?", memberId).Scan(&utente.Name, &utente.Photo, &utente.Number)
				if err != nil {
					return []Conversation{}, err
				}
				//aggiungo l'utente alla lista dei membri della conversazione
				conversation.Members = append(conversation.Members, utente)
			}
			// 0 -> no gruppo
			if isGroup == 0 && userId != memberId {
				conversation.Photo = utente.Photo //se non è un gruppo, allora la foto e il nome sono quelle dell'altro utente
				conversation.Title = utente.Name
			}
		}

		userConversations = append(userConversations, conversation)
	}
	return userConversations, nil
}
