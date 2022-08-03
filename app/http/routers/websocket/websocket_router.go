package websocket

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/http/validator/common/websocket"
)

type WebSocketRouter struct {
}

func (s *WebSocketRouter) InitWebSocketRouter(router *gin.RouterGroup) {
	router.GET("ws", (&websocket.Connect{}).CheckParams)
}
