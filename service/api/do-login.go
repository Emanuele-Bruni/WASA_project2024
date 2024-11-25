package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type User struct {
	name   string
	photo  Photo
	bio    string
	number string
}

type Photo struct {
	photoId      int
	photoContent string
	width        int
	height       int
}

type Conversation struct {
	conversatioId int
	title         string
	photo         Photo
	isGroup       bool
}

type Message struct {
	messageId int
	author    string
	text      string
	hours     string
	comments  []Comment
}

type Comment struct {
	commentId int
	author    string
	text      string
}

type Group struct {
	groupId    int
	groupPhoto Photo
	groupName  string
	members    []User
}

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}
