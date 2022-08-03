package api

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/http/controller/api"
)

type HomeRouter struct {
}

func (s *HomeRouter) InitHomeRouter(Router *gin.RouterGroup) {
	home := Router.Group("/")
	{
		home.GET("news", (&api.Home{}).News)
	}
}
