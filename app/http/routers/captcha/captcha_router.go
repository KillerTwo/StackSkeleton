package captcha

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/http/controller/captcha"
)

type CaptchaRouter struct {
}

func (s *CaptchaRouter) InitCaptchaRouter(router *gin.RouterGroup) {
	captchaRouter := router.Group("/")
	{
		captchaRouter.GET("/", (&captcha.Captcha{}).GenerateId)                          // 获取图片验证码ID
		captchaRouter.GET("/:captcha_id", (&captcha.Captcha{}).GetImg)                   // 获取验证码图片
		captchaRouter.GET("/:captcha_id/:captcha_value", (&captcha.Captcha{}).CheckCode) // 校验验证码
	}
}
