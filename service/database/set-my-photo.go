package database

import "WASA_project2024/service/media"

func (db *appdbimpl) SetMyPhoto(photo media.Image, id int) error {
	return db.setPhoto(photo, id, false)
}
