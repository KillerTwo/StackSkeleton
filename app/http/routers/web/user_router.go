package web

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/http/controller/web"
	"goskeleton/app/http/middleware/authorization"
)

type UserRouter struct {
}

// InitUserRouter 需要登录
func (s *UserRouter) InitUserRouter(router *gin.RouterGroup) {

	userRouter := router.Group("/")
	{
		userRouter.POST("index", (&web.Users{}).Show)
		userRouter.POST("create", (&web.Users{}).Store)
		userRouter.POST("edit", (&web.Users{}).Update)
		userRouter.DELETE("delete/:userId", (&web.Users{}).Destroy)
		userRouter.GET("getUserInfo", (&web.Users{}).CurrentUser)
		userRouter.GET("logout", (&web.Users{}).Logout)
	}

}

// InitNoAuthUserRouter 不需要登录
func (s *UserRouter) InitNoAuthUserRouter(router *gin.RouterGroup) {
	userRouter := router.Group("/")
	{
		userRouter.POST("register", (&web.Users{}).Register)
		userRouter.POST("login", (&web.Users{}).Login)
		userRouter.Use(authorization.RefreshTokenConditionCheck()).POST("refreshtoken", (&web.Users{}).RefreshToken)
	}
}
