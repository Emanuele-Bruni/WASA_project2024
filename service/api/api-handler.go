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
	rt.router.POST("/login", rt.wrap(rt.doLogin))                    //OK
	rt.router.PATCH("/settings/username", rt.wrap(rt.setMyUserName)) //OK
	rt.router.GET("/conversations", rt.wrap(rt.getMyConversations))
	rt.router.GET("/conversations/:conversationid", rt.wrap(rt.getConversation))
	rt.router.POST("/conversations/:conversationid/message", rt.wrap(rt.sendMessage))
	rt.router.POST("/conversations/:conversationid/messages/:messageid/forward/:forwardid", rt.wrap(rt.forwardMessage))
	rt.router.POST("/conversations/:conversationid/messages/:messageid/comments", rt.wrap(rt.commentMessage))
	rt.router.DELETE("/conversations/:conversationid/messages/:messageid/comments/:commentid", rt.wrap(rt.uncommentMessage))
	//rt.router.DELETE("/conversations/{conversationid}/messages/{messageid}", rt.wrap(rt.deleteMessage))
	//rt.router.PATCH("/conversations/{conversationid}/members", rt.wrap(rt.addToGroup))
	//rt.router.DELETE("/conversations/{conversationid}/members/{userid}", rt.wrap(rt.leaveGroup))
	//rt.router.PATCH("/conversations/{conversationid}/name", rt.wrap(rt.setGroupName))
	//rt.router.PATCH("/settings/photo", rt.wrap(rt.setMyPhoto))
	//rt.router.PATCH("/conversations/{conversationid}/photo", rt.wrap(rt.setGroupPhoto))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
