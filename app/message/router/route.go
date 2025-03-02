package router

import (
	"github.com/gin-gonic/gin"
	"github.com/woshilaixuex/csd_chat_backend/app/message/hander"
	"github.com/woshilaixuex/csd_chat_backend/app/message/ws"
	"github.com/woshilaixuex/csd_chat_backend/app/util/middleware"
)

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description:
 * @Date: 2025-02-27 15:35
 */
func RouterRegister(r *gin.Engine) {

	chatGroup := r.Group("/chat")
	chatGroup.Use(middleware.TokenAuthMiddleware())
	chatGroup.POST("/private/send", hander.PrivateChatHander)
	chatGroup.GET("/private/ws", ws.WebSocketHandler)
}
