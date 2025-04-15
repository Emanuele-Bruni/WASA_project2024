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
	rt.router.POST("/login", rt.wrap(rt.doLogin))                                                                            //OK
	rt.router.PUT("/settings/username", rt.wrap(rt.setMyUserName))                                                           //OK
	rt.router.GET("/conversations", rt.wrap(rt.getMyConversations))                                                          //?? con i gruppi
	rt.router.GET("/conversations/:conversationid", rt.wrap(rt.getConversation))                                             //?? con i gruppi
	rt.router.POST("/conversations/:conversationid/message", rt.wrap(rt.sendMessage))                                        //OK
	rt.router.POST("/conversations/:conversationid/messages/:messageid/forward/:forwardid", rt.wrap(rt.forwardMessage))      //OK
	rt.router.POST("/conversations/:conversationid/messages/:messageid/comments", rt.wrap(rt.commentMessage))                //OK
	rt.router.DELETE("/conversations/:conversationid/messages/:messageid/comments/:commentid", rt.wrap(rt.uncommentMessage)) //OK
	rt.router.DELETE("/conversations/:conversationid/messages/:messageid", rt.wrap(rt.deleteMessage))                        //OK
	rt.router.PUT("/conversations/:conversationid/members", rt.wrap(rt.addToGroup))
	rt.router.DELETE("/conversations/:conversationid/members", rt.wrap(rt.leaveGroup))
	rt.router.PUT("/conversations/:conversationid/name", rt.wrap(rt.setGroupName))
	rt.router.PUT("/settings/photo", rt.wrap(rt.setMyPhoto))
	rt.router.PUT("/conversations/:conversationid/photo", rt.wrap(rt.setGroupPhoto))
	rt.router.GET("/users", rt.wrap(rt.getUsers))                 //OK
	rt.router.POST("/users/:userid", rt.wrap(rt.newConversation)) //OK
	rt.router.POST("/conversations", rt.wrap(rt.newGroup))        //??

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
