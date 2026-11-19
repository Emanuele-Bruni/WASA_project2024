package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// EndPoints
	rt.router.POST("/login", rt.wrap(rt.doLogin))
	rt.router.POST("/logout", rt.protected(rt.logout, memberAccess))
	rt.router.PUT("/settings/username", rt.protected(rt.setMyUserName, memberAccess))
	rt.router.GET("/conversations", rt.protected(rt.getMyConversations, memberAccess))
	rt.router.GET("/conversations/:conversationid", rt.protected(rt.getConversation, memberAccess))
	rt.router.POST("/conversations/:conversationid/messages", rt.protected(rt.sendMessage, memberAccess))
	rt.router.POST("/conversations/:conversationid/messages/:messageid/forward/:forwardid", rt.protected(rt.forwardMessage, memberAccess))
	rt.router.POST("/conversations/:conversationid/messages/:messageid/comments", rt.protected(rt.commentMessage, memberAccess))
	rt.router.DELETE("/conversations/:conversationid/messages/:messageid/comments/:commentid", rt.protected(rt.uncommentMessage, commentAuthorAccess))
	rt.router.DELETE("/conversations/:conversationid/messages/:messageid", rt.protected(rt.deleteMessage, messageAuthorAccess))
	rt.router.PUT("/conversations/:conversationid/members", rt.protected(rt.addToGroup, groupAccess))
	rt.router.DELETE("/conversations/:conversationid/members", rt.protected(rt.leaveGroup, groupAccess))
	rt.router.PUT("/conversations/:conversationid/name", rt.protected(rt.setGroupName, groupAccess))
	rt.router.PUT("/settings/photo", rt.protected(rt.setMyPhoto, memberAccess))
	rt.router.PUT("/conversations/:conversationid/photo", rt.protected(rt.setGroupPhoto, groupAccess))
	rt.router.GET("/users", rt.protected(rt.getUsers, memberAccess))
	rt.router.POST("/users/:userid", rt.protected(rt.newConversation, memberAccess))
	rt.router.POST("/conversations", rt.protected(rt.newGroup, memberAccess))
	rt.router.GET("/conversations/:conversationid/messages", rt.protected(rt.getMessages, memberAccess))
	rt.router.GET("/settings/photo/:userid", rt.protected(rt.getMyPhoto, ownProfileAccess))

	rt.router.GET("/media/:mediaid", rt.protected(rt.getMedia, memberAccess))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	rt.router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { writeError(w, 404, "Rotta non trovata") })
	rt.router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { writeError(w, 405, "Metodo non consentito") })
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength > maxBodyBytes {
			writeError(w, 413, "Il corpo della richiesta supera 2 MiB")
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
		rt.router.ServeHTTP(w, r)
	})
}
