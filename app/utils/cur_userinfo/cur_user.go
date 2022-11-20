package cur_userinfo

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/variable"
	"goskeleton/app/http/middleware/my_jwt"
)

// GetCurrentUserId 获取当前用户的id
// @context 请求上下文
func GetCurrentUserId(context *gin.Context) (int64, bool) {
	tokenKey := variable.ConfigYml.GetString("Token.BindContextKeyName")
	currentUser, exist := context.MustGet(tokenKey).(*my_jwt.CustomClaims)
	return currentUser.UserId, exist
}

// GetCurrentTokenClaims 获取当前token Claims
func GetCurrentTokenClaims(context *gin.Context) (*my_jwt.CustomClaims, bool) {
	tokenKey := variable.ConfigYml.GetString("Token.BindContextKeyName")
	currentUser, exist := context.MustGet(tokenKey).(*my_jwt.CustomClaims)
	return currentUser, exist
}

// GetCurrentToken 获取当前用户的token
func GetCurrentToken(ctx *gin.Context) (string, bool) {
	currentToken, ok := ctx.MustGet("currentToken").(string)
	return currentToken, ok
}
