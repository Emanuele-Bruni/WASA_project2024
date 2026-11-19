package api

import (
	"WASA_project2024/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"strings"
)

func (rt *_router) getMedia(w http.ResponseWriter, r *http.Request, p httprouter.Params, ctx reqcontext.RequestContext) {
	image, err := rt.db.GetImage(ctx.MediaID)
	if err != nil {
		respondDatabaseError(w, err, ctx)
		return
	}
	etag := `"` + image.Digest + `"`
	w.Header().Set("Content-Type", image.MimeType)
	w.Header().Set("Cache-Control", "private, max-age=31536000, immutable")
	w.Header().Set("Vary", "Authorization")
	w.Header().Set("ETag", etag)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	for _, candidate := range strings.Split(r.Header.Get("If-None-Match"), ",") {
		if candidate = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(candidate), "W/")); candidate == etag || candidate == "*" {
			w.WriteHeader(304)
			return
		}
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(image.Data)))
	w.WriteHeader(200)
	if _, err = w.Write(image.Data); err != nil {
		ctx.Logger.WithError(err).Error("Write image")
	}
}
