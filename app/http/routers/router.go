package routers

import (
	api "goskeleton/app/http/routers/api"
	"goskeleton/app/http/routers/captcha"
	"goskeleton/app/http/routers/web"
	"goskeleton/app/http/routers/websocket"
)

type RouterGroup struct {
	HomeRouter      api.HomeRouter
	UploadRouter    web.UploadRouter
	UserRouter      web.UserRouter
	CaptchaRouter   captcha.CaptchaRouter
	WebSocketRouter websocket.WebSocketRouter
	MenuRouter      web.MenuRouter
}

var RouterGroupApp = new(RouterGroup)
