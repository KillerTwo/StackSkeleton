package web

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/http/controller/web"
)

type MenuRouter struct {
}

// InitMenuRouter 初始化路由
func (s *MenuRouter) InitMenuRouter(router *gin.RouterGroup) {

	menuRouter := router.Group("/")
	{
		// menuRouter.GET("routes", (&web.Menu{}).GetRoutes)
		menuRouter.GET("routes", (&web.Menu{}).GetCurrentRoutes)
		// menuRouter.GET("currentRoutes", (&web.Menu{}).GetCurrentRoutes)
	}

}
