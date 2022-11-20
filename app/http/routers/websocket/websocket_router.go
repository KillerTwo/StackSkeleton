package websocket

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/http/request/common/websocket"
)

type WebSocketRouter struct {
}

func (s *WebSocketRouter) InitWebSocketRouter(router *gin.RouterGroup) {
	router.GET("ws", (&websocket.Connect{}).CheckParams)
}
