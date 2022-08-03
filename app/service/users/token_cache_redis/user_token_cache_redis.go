package token_cache_redis

import (
	"fmt"
	"go.uber.org/zap"
	"goskeleton/app/global/variable"
	userstoken "goskeleton/app/service/users/token"
	"goskeleton/app/utils/redis_factory"
	"strconv"
	"strings"
	"time"
)

func CreateUsersTokenCacheFactory(userId int64) *userTokenCacheRedis {
	redCli := redis_factory.GetOneRedisClient()
	if redCli == nil {
		return nil
	}
	return &userTokenCacheRedis{redisClient: redCli, userTokenKey: "token_userid_" + strconv.FormatInt(userId, 10)}
}

type userTokenCacheRedis struct {
	redisClient  *redis_factory.RedisClient
	userTokenKey string
}

// SetTokenCache 设置缓存
func (u *userTokenCacheRedis) SetTokenCache(tokenExpire int64, token string) bool {
	if _, err := u.redisClient.Int(u.redisClient.Execute("zAdd", u.userTokenKey, tokenExpire, token)); err == nil {
		return true
	} else {
		variable.ZapLog.Error("缓存用户token到redis出错", zap.Error(err))
	}
	return false
}

// DelOverMaxOnlineCache 删除缓存,删除超过系统允许最大在线数量之外的用户
func (u *userTokenCacheRedis) DelOverMaxOnlineCache() bool {
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
func (u *userTokenCacheRedis) OauthCheckTokenIsOk() bool {
	onlineUsers := variable.ConfigYml.GetInt("Token.JwtTokenOnlineUsers")
	if count, err := u.redisClient.Int64(u.redisClient.Execute("zCard", u.userTokenKey)); err == nil {
		return int64(onlineUsers) >= count
	}
	return false
}

// TokenCacheIsExists 查询token是否在redis存在
func (u *userTokenCacheRedis) TokenCacheIsExists(token string) (exists bool) {
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
func (u *userTokenCacheRedis) UserTokenCacheIsExists(token string) (exists bool) {
	// s, _ := u.redisClient.Execute("mget", token)
	// fmt.Printf("%s", s)
	if s, err := u.redisClient.Strings(u.redisClient.Execute("mget", token)); err == nil && s[0] != "" {
		// 更新token在redis中的过期时间
		after30, _ := time.ParseDuration("15m")
		t := time.Now().Add(after30)
		userstoken.CreateUserFactory().RefreshTokenExpire(token, t.Unix())
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
func (u *userTokenCacheRedis) SetUserTokenExpire(ts int64) bool {
	if _, err := u.redisClient.Execute("expireAt", u.userTokenKey, ts); err == nil {
		return true
	}
	return false
}

// SetToken SET key
func (u *userTokenCacheRedis) SetToken(key string, value string, expireAt int64) bool {
	if _, err := u.redisClient.Execute("setex", key, expireAt, value); err == nil {
		return true
	} else {
		fmt.Errorf(err.Error())
	}
	return false
}

// ClearUserToken 清除某个用户的全部缓存，当用户更改密码或者用户被禁用则删除该用户的全部缓存
func (u *userTokenCacheRedis) ClearUserToken() bool {
	if _, err := u.redisClient.Execute("del", u.userTokenKey); err == nil {
		return true
	}
	return false
}

// ReleaseRedisConn 释放redis
func (u *userTokenCacheRedis) ReleaseRedisConn() {
	u.redisClient.ReleaseOneRedisClient()
}
