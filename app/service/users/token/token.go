package token

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/global/variable"
	"goskeleton/app/http/middleware/my_jwt"
	"goskeleton/app/model"
	"goskeleton/app/service/users/token_cache_redis"
	"goskeleton/app/utils/redis_factory"
	"time"
)

// CreateUserFactory 创建 userToken 工厂
func CreateUserFactory() *userToken {
	return &userToken{
		userJwt: my_jwt.CreateMyJWT(variable.ConfigYml.GetString("Token.JwtTokenSignKey")),
	}
}

type userToken struct {
	userJwt *my_jwt.JwtSign
}

//GenerateToken 生成token
func (u *userToken) GenerateToken(user *model.UsersModel, roleKeys []string, expireAt int64) (tokens string, err error) {
	// 根据实际业务自定义token需要包含的参数，生成token，注意：用户密码请勿包含在token
	customClaims := my_jwt.CustomClaims{
		UserId: user.Id,
		Name:   user.UserName,
		Phone:  user.Phone,
		Roles:  roleKeys,
		// 特别注意，针对前文的匿名结构体，初始化的时候必须指定键名，并且不带 jwt. 否则报错：Mixture of field: value and value initializers
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 10,       // 生效开始时间
			ExpiresAt: time.Now().Unix() + expireAt, // 失效截止时间
		},
	}
	return u.userJwt.CreateToken(customClaims)
}

// RecordLoginToken 用户login成功，记录用户token
func (u *userToken) RecordLoginToken(userToken, clientIp string) bool {
	if customClaims, err := u.userJwt.ParseToken(userToken); err == nil {
		userId := customClaims.UserId
		expiresAt := customClaims.ExpiresAt
		return model.CreateUserFactory("").OauthLoginToken(userId, userToken, expiresAt, clientIp)
	} else {
		return false
	}
}

// InvalidLoginToken 清除登录token
func (u *userToken) InvalidLoginToken(token string) bool {
	return u.ClearTokenByKey(token)
}

// ClearTokenByKey 删除指定的token
func (u *userToken) ClearTokenByKey(key string) bool {
	redisClient := redis_factory.GetOneRedisClient()
	if redisClient != nil {
		// 设置token，并设置有效时间
		if _, err := redisClient.Execute("del", key); err == nil {
			return true
		}
	}
	return false
}

//TokenIsMeetRefreshCondition 检查token是否满足刷新条件
func (u *userToken) TokenIsMeetRefreshCondition(token string) bool {
	// token基本信息是否有效：1.过期时间在允许的过期范围内;2.基本格式正确
	customClaims, code := u.isNotExpired(token, variable.ConfigYml.GetInt64("Token.JwtTokenRefreshAllowSec"))
	switch code {
	case consts.JwtTokenOK, consts.JwtTokenExpired:
		//在数据库的存储信息是否也符合过期刷新刷新条件
		if model.CreateUserFactory("").OauthRefreshConditionCheck(customClaims.UserId, token) {
			return true
		}
	}
	return false
}

// RefreshToken 刷新token的有效期（默认+3600秒，参见常量配置项）
func (u *userToken) RefreshToken(oldToken, clientIp string) (newToken string, res bool) {
	var err error
	//如果token是有效的、后者在在过期时间内，那么执行更新，换取新token
	if newToken, err = u.userJwt.RefreshToken(oldToken, variable.ConfigYml.GetInt64("Token.JwtTokenRefreshExpireAt")); err == nil {
		if customClaims, err := u.userJwt.ParseToken(newToken); err == nil {
			userId := customClaims.UserId
			expiresAt := customClaims.ExpiresAt
			if model.CreateUserFactory("").OauthRefreshToken(userId, expiresAt, oldToken, newToken, clientIp) {
				return newToken, true
			}
		}
	}

	return "", false
}

// 判断token本身是否未过期
// 参数解释：
// token： 待处理的token值
// expireAtSec： 过期时间延长的秒数，主要用于用户刷新token时，判断是否在延长的时间范围内，非刷新逻辑默认为0
func (u *userToken) isNotExpired(token string, expireAtSec int64) (*my_jwt.CustomClaims, int) {
	if customClaims, err := u.userJwt.ParseToken(token); err == nil {

		if time.Now().Unix()-(customClaims.ExpiresAt+expireAtSec) < 0 {
			// token有效
			return customClaims, consts.JwtTokenOK
		} else {
			// 过期的token
			return customClaims, consts.JwtTokenExpired
		}
	} else {
		// 无效的token
		return nil, consts.JwtTokenInvalid
	}
}

// IsEffective 判断token是否有效
func (u *userToken) IsEffective(token string) (*my_jwt.CustomClaims, bool) {
	// customClaims, code := u.isNotExpired(token, 0)
	customClaims, err := u.ParseToken(token)
	if err != nil {
		//1.首先在redis检测是否存在某个用户对应的有效token，如果存在就直接返回，不再继续查询mysql，否则最后查询mysql逻辑，确保万无一失
		if variable.ConfigYml.GetInt("Token.IsCacheToRedis") == 1 {
			tokenRedisFact := token_cache_redis.CreateUsersTokenCacheFactory(customClaims.UserId)
			if tokenRedisFact != nil {
				defer tokenRedisFact.ReleaseRedisConn()
				/*if tokenRedisFact.TokenCacheIsExists(token) {
					return true
				}*/
				if tokenRedisFact.UserTokenCacheIsExists(token) {
					return &customClaims, true
				}
			}
			//2.token符合token本身的规则以后，继续在数据库校验是不是符合本系统其他设置，例如：一个用户默认只允许10个账号同时在线（10个token同时有效）
			/*if tokenRedisFact.OauthCheckTokenIsOk(token) {
				return true
			}*/
		}
	}
	return nil, false
}

// ParseToken 将 token 解析为绑定时传递的参数
func (u *userToken) ParseToken(tokenStr string) (CustomClaims my_jwt.CustomClaims, err error) {
	if customClaims, err := u.userJwt.ParseToken(tokenStr); err == nil {
		return *customClaims, nil
	} else {
		return my_jwt.CustomClaims{}, errors.New(my_errors.ErrorsParseTokenFail)
	}
}

// RefreshTokenExpire 刷新token的过期时间
func (u *userToken) RefreshTokenExpire(token string, expireAt int64) bool {
	if redCli := redis_factory.GetOneRedisClient(); redCli != nil {
		if _, err := redCli.Execute("EXPIREAT", token, expireAt); err == nil {
			return true
		}
	}
	return false
}
