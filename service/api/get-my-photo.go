package api

import (
	"WASA_project2024/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) getMyPhoto(w http.ResponseWriter, r *http.Request, p httprouter.Params, ctx reqcontext.RequestContext) {
	rt.respondUser(w, ctx)
}
