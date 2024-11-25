package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	//EndPoints
	rt.router.POST("/login", rt.doLogin)

	//rt.router.PATCH("/settings/username", rt.wrap(rt.setMyUserName))
	//rt.router.GET("/conversations", rt.wrap(rt.getMyConversations))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
