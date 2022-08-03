package token

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/global/variable"
	"goskeleton/app/http/middleware/my_jwt"
	"goskeleton/app/model"
	"goskeleton/app/utils/redis_factory"
	"strconv"
	"strings"
	"time"
)

// CreateUserTokenFactory 创建 userToken 工厂
func CreateUserTokenFactory() *userToken {
	redCli := redis_factory.GetOneRedisClient()
	if redCli == nil {
		return nil
	}
	return &userToken{
		userJwt:      my_jwt.CreateMyJWT(variable.ConfigYml.GetString("Token.JwtTokenSignKey")),
		redisClient:  redCli,
		userTokenKey: "",
	}
}

// UserTokenFactory
func UserTokenFactory(userId int64) *userToken {
	redCli := redis_factory.GetOneRedisClient()
	if redCli == nil {
		return nil
	}
	return &userToken{
		userJwt:      my_jwt.CreateMyJWT(variable.ConfigYml.GetString("Token.JwtTokenSignKey")),
		redisClient:  redCli,
		userTokenKey: "token_userid_" + strconv.FormatInt(userId, 10),
	}
}

type userToken struct {
	userJwt      *my_jwt.JwtSign
	redisClient  *redis_factory.RedisClient
	userTokenKey string
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
		return model.UserModelFactory("").OauthLoginToken(userId, userToken, expiresAt, clientIp)
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
			defer u.ReleaseRedisConn()
			if u.UserTokenCacheIsExists(token) {
				return &customClaims, true
			}
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

// StoreToken 将token 对于的用户信息存储到redis
func (u *userToken) StoreToken(userId, expiresAt int64, token, value string) bool {
	defer u.ReleaseRedisConn()
	if success := u.SetTokenCache(expiresAt, token); success {
		expireSec := variable.ConfigYml.GetInt64("Token.redisTokenExpire")
		success = u.SetToken(token, string(value), expireSec*60)
		// 缓存结束之后删除超过系统设置最大在线数量的token
		u.DelOverMaxOnlineCache()
		return success
	}
	return false
}

// DelTokenCacheFromRedis 用户密码修改后，删除redis所有的token
func (u *userToken) DelTokenCacheFromRedis(userId int64) bool {
	u.ClearUserToken()
	u.ReleaseRedisConn()
	return true
}

// SetTokenCache 设置缓存
func (u *userToken) SetTokenCache(tokenExpire int64, token string) bool {
	if _, err := u.redisClient.Int(u.redisClient.Execute("zAdd", u.userTokenKey, tokenExpire, token)); err == nil {
		return true
	} else {
		variable.ZapLog.Error("缓存用户token到redis出错", zap.Error(err))
	}
	return false
}

// DelOverMaxOnlineCache 删除缓存,删除超过系统允许最大在线数量之外的用户
func (u *userToken) DelOverMaxOnlineCache() bool {
	// 首先先删除过期的token
	_, _ = u.redisClient.Execute("zRemRangeByScore", u.userTokenKey, 0, time.Now().Unix()-1)

	onlineUsers := variable.ConfigYml.GetInt("Token.JwtTokenOnlineUsers")
	alreadyCacheNum, err := u.redisClient.Int(u.redisClient.Execute("zCard", u.userTokenKey))
	if err == nil && alreadyCacheNum > onlineUsers {
		if invalidTokens, err := u.redisClient.Execute("ZRANGE", u.userTokenKey, 0, alreadyCacheNum-onlineUsers-1); err == nil {
			if tokens, ok := invalidTokens.([]interface{}); ok {
				_, err := u.redisClient.Execute("del", tokens...)
				if err != nil {
					variable.ZapLog.Error("删除超过系统允许之外的token出错：", zap.Error(err))
					return false
				}
			}
		}
		// 删除超过最大在线数量之外的token
		if alreadyCacheNum, err = u.redisClient.Int(u.redisClient.Execute("zRemRangeByRank", u.userTokenKey, 0, alreadyCacheNum-onlineUsers-1)); err == nil {
			return true
		} else {
			variable.ZapLog.Error("删除超过系统允许之外的token出错：", zap.Error(err))
		}
	}
	return false
}

// OauthCheckTokenIsOk 用户是否符合同时在线要求
func (u *userToken) OauthCheckTokenIsOk() bool {
	onlineUsers := variable.ConfigYml.GetInt("Token.JwtTokenOnlineUsers")
	if count, err := u.redisClient.Int64(u.redisClient.Execute("zCard", u.userTokenKey)); err == nil {
		return int64(onlineUsers) >= count
	}
	return false
}

// TokenCacheIsExists 查询token是否在redis存在
func (u *userToken) TokenCacheIsExists(token string) (exists bool) {
	// token = md5_encrypt.MD5(token)
	curTimestamp := time.Now().Unix()
	onlineUsers := variable.ConfigYml.GetInt("Token.JwtTokenOnlineUsers")
	if strSlice, err := u.redisClient.Strings(u.redisClient.Execute("zRevRange", u.userTokenKey, 0, onlineUsers-1)); err == nil {
		for _, val := range strSlice {
			if score, err := u.redisClient.Int64(u.redisClient.Execute("zScore", u.userTokenKey, token)); err == nil {
				if score > curTimestamp {
					if strings.Compare(val, token) == 0 {
						exists = true
						break
					}
				}
			}
		}
	} else {
		variable.ZapLog.Error("获取用户在redis缓存的 token 值出错：", zap.Error(err))
	}
	return
}

// UserTokenCacheIsExists 登录用户是否有效
func (u *userToken) UserTokenCacheIsExists(token string) (exists bool) {
	// s, _ := u.redisClient.Execute("mget", token)
	// fmt.Printf("%s", s)
	if s, err := u.redisClient.Strings(u.redisClient.Execute("mget", token)); err == nil && s[0] != "" {
		// 更新token在redis中的过期时间
		after30, _ := time.ParseDuration("15m")
		t := time.Now().Add(after30)
		u.RefreshTokenExpire(token, t.Unix())
		return true
	} else if err != nil {
		variable.ZapLog.Error("获取用户在redis缓存的 token 值出错：", zap.Error(err))
	} else {
		return false
	}
	return false
}

// SetUserTokenExpire 设置用户的 usertoken 键过期时间
// 参数： 时间戳
func (u *userToken) SetUserTokenExpire(ts int64) bool {
	if _, err := u.redisClient.Execute("expireAt", u.userTokenKey, ts); err == nil {
		return true
	}
	return false
}

// SetToken SET key
func (u *userToken) SetToken(key string, value string, expireAt int64) bool {
	if _, err := u.redisClient.Execute("setex", key, expireAt, value); err == nil {
		return true
	} else {
		fmt.Errorf(err.Error())
	}
	return false
}

// ClearUserToken 清除某个用户的全部缓存，当用户更改密码或者用户被禁用则删除该用户的全部缓存
func (u *userToken) ClearUserToken() bool {
	if _, err := u.redisClient.Execute("del", u.userTokenKey); err == nil {
		return true
	}
	return false
}

// ReleaseRedisConn 释放redis
func (u *userToken) ReleaseRedisConn() {
	u.redisClient.ReleaseOneRedisClient()
}
