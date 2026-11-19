package database

import (
	"WASA_project2024/service/media"
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type ImageMetadata struct {
	ImageID  int    `json:"imageId"`
	URL      string `json:"url"`
	MimeType string `json:"mimeType"`
	Name     string `json:"name"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Size     int    `json:"size"`
}

var avatarsOnce sync.Once
var contactAvatar, groupAvatar media.Image
var avatarError error

func defaultAvatar(group bool) (media.Image, error) {
	avatarsOnce.Do(func() {
		contactAvatar, avatarError = media.Decode(media.Upload{MimeType: "image/png", Data: DefaultContactPhoto, Name: "Profilo.png"})
		if avatarError == nil {
			groupAvatar, avatarError = media.Decode(media.Upload{MimeType: "image/png", Data: DefaultGroupPhoto, Name: "Gruppo.png"})
		}
	})
	if group {
		return groupAvatar, avatarError
	}
	return contactAvatar, avatarError
}
func saveImage(tx *sql.Tx, img media.Image) (int, error) {
	_, err := tx.Exec("INSERT INTO images(mime_type,name,width,height,digest,data) VALUES(?,?,?,?,?,?) ON CONFLICT(digest) DO NOTHING", img.MimeType, img.Name, img.Width, img.Height, img.Digest, img.Data)
	if err != nil {
		return 0, err
	}
	var id int
	err = tx.QueryRow("SELECT image_id FROM images WHERE digest=?", img.Digest).Scan(&id)
	return id, err
}
func (db *appdbimpl) GetImage(id int) (media.Image, error) {
	var image media.Image
	err := db.c.QueryRow("SELECT mime_type,name,width,height,digest,data FROM images WHERE image_id=?", id).Scan(&image.MimeType, &image.Name, &image.Width, &image.Height, &image.Digest, &image.Data)
	return image, err
}
func (db *appdbimpl) CanReadImage(id, userID int) (bool, error) {
	var allowed bool
	err := db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM users_table WHERE photo_image_id=?)
 OR EXISTS(SELECT 1 FROM conversations_table c JOIN conversation_members cm ON cm.conversation_id=c.conversation_id WHERE c.photo_image_id=? AND cm.user_id=?)
 OR EXISTS(SELECT 1 FROM messages_table m JOIN conversation_members cm ON cm.conversation_id=m.conversation_id WHERE m.image_id=? AND cm.user_id=?)`, id, id, userID, id, userID).Scan(&allowed)
	return allowed, err
}
func (db *appdbimpl) setPhoto(img media.Image, id int, group bool) error {
	tx, err := db.c.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	imageID, err := saveImage(tx, img)
	if err != nil {
		return err
	}
	query := "UPDATE users_table SET photo='',photo_image_id=? WHERE id=?"
	if group {
		query = "UPDATE conversations_table SET photo='',photo_image_id=? WHERE conversation_id=? AND is_group=1"
	}
	if err = requireAffected(tx.Exec(query, imageID, id)); err != nil {
		return err
	}
	return tx.Commit()
}

// Version 2 data conversion is part of the migration transaction: IDs and history survive.
func migrateLegacyMedia(tx *sql.Tx) error {
	for _, group := range []bool{false, true} {
		table, id := "users_table", "id"
		if group {
			table, id = "conversations_table", "conversation_id"
		}
		rows, err := tx.Query("SELECT " + id + ",photo FROM " + table)
		if err != nil {
			return err
		}
		type legacy struct {
			id   int
			data string
		}
		entries := make([]legacy, 0)
		for rows.Next() {
			var row legacy
			if err = rows.Scan(&row.id, &row.data); err != nil {
				rows.Close()
				return err
			}
			entries = append(entries, row)
		}
		err = rows.Err()
		rows.Close()
		if err != nil {
			return err
		}
		for _, row := range entries {
			decoded, err := decodeLegacyImage(row.data)
			preservedLegacy := ""
			if err != nil {
				preservedLegacy = row.data
				decoded, err = defaultAvatar(group)
				if err != nil {
					return err
				}
			}
			imageID, err := saveImage(tx, decoded)
			if err != nil {
				return err
			}
			if _, err = tx.Exec("UPDATE "+table+" SET photo=?,photo_image_id=? WHERE "+id+"=?", preservedLegacy, imageID, row.id); err != nil {
				return err
			}
		}
	}
	rows, err := tx.Query("SELECT message_id,text FROM messages_table WHERE substr(text,1,6)='[foto]'")
	if err != nil {
		return err
	}
	type legacyMessage struct {
		id   int
		text string
	}
	messages := make([]legacyMessage, 0)
	for rows.Next() {
		var row legacyMessage
		if err = rows.Scan(&row.id, &row.text); err != nil {
			rows.Close()
			return err
		}
		messages = append(messages, row)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return err
	}
	for _, message := range messages {
		img, err := decodeLegacyImage(strings.TrimPrefix(message.text, "[foto]"))
		if err != nil {
			continue
		} // Preserve invalid legacy payloads as literal text.
		id, err := saveImage(tx, img)
		if err != nil {
			return err
		}
		if _, err = tx.Exec("UPDATE messages_table SET text='Immagine',content_type='image',image_id=?,image_name='Immagine' WHERE message_id=?", id, message.id); err != nil {
			return err
		}
	}
	return nil
}
func decodeLegacyImage(encoded string) (media.Image, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return media.Image{}, err
	}
	return media.Decode(media.Upload{MimeType: http.DetectContentType(data), Data: encoded, Name: "Immagine"})
}
func imageURL(id int) string { return fmt.Sprintf("/media/%d", id) }

func (db *appdbimpl) ImageExists(id int) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM images WHERE image_id=?)", id).Scan(&exists)
	return exists, err
}
