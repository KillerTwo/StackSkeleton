package websocket

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/response"
)

type Connect struct {
	Token string `form:"token" binding:"required,min=10"`
}

// 验证器语法，参见 Register.go文件，有详细说明

func (c Connect) CheckParams(context *gin.Context) {

	// 1. 首先检查是否开启websocket服务配置（在配置项中开启）
	if variable.ConfigYml.GetInt("Websocket.Start") != 1 {
		response.Fail(context, consts.WsServerNotStartCode, consts.WsServerNotStartMsg, "")
		return
	}
	//2.基本的验证规则没有通过
	if err := context.ShouldBind(&c); err != nil {
		errs := gin.H{
			"tips": "请在get参数中提交token信息,demo格式：ws://127.0.0.1:2020?token=this is a series token...",
			"err":  err.Error(),
		}
		response.ErrorParam(context, errs)
		return
	}
}
