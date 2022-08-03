package web

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/service/menus/service"
	"goskeleton/app/utils/cur_userinfo"
	"goskeleton/app/utils/response"
)

type Menu struct {
}

// GetRoutes 获取路由
func (m *Menu) GetRoutes(context *gin.Context) {
	list, err := service.MenuServiceFactory().SelectTreeRoutes()
	if err != nil {
		response.Fail(context, consts.InternalServerErrorCode, consts.InternalServerErrorMsg, nil)
		return
	}
	response.Success(context, consts.CurdStatusOkMsg, list)
}

func (m *Menu) GetCurrentRoutes(ctx *gin.Context) {
	if currentUser, exist := cur_userinfo.GetCurrentTokenClaims(ctx); exist {
		roles := currentUser.Roles
		menus, err := service.MenuServiceFactory().SelectTreeRoutesByRole(roles)
		if err != nil {
			response.Fail(ctx, consts.NotAuthorize, consts.UserIdNotExist, nil)
		} else {
			response.Success(ctx, consts.CurdStatusOkMsg, menus)
		}
	} else {
		response.Fail(ctx, consts.NotAuthorize, consts.UserIdNotExist, nil)
	}
}
