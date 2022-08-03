package service

import (
	"encoding/json"
	"fmt"
	"goskeleton/app/global/variable"
	"goskeleton/app/http/response"
	"goskeleton/app/model"
	token2 "goskeleton/app/service/users/token"
	"goskeleton/app/utils/md5_encrypt"
	"strconv"
)

func UserServiceFactory() *UsersService {
	return &UsersService{model.UserModelFactory("")}
}

type UsersService struct {
	userModel *model.UsersModel
}

func (u *UsersService) FindById(userId int64) *response.UserResponse {
	user := u.userModel.FindById(userId)
	if user != nil {
		return &response.UserResponse{
			UserId:   strconv.FormatInt(user.Id, 10),
			Username: user.UserName,
			RealName: user.RealName,
			Avatar:   user.Avatar,
			Desc:     "manager",
			Password: "",
			Token:    "fakeToken1",
			HomePath: "/dashboard/analysis",
			Roles: []response.RoleResponse{{
				RoleName: "Super Admin",
				Value:    "super",
			}},
		}
	}

	return nil
}

func (u *UsersService) Register(userName, pass, userIp string) bool {
	pass, err := md5_encrypt.BcryptEncode(pass) // 预先处理密码加密，然后存储在数据库
	if err != nil {
		return false
	}
	return u.userModel.Register(userName, pass, userIp)
}

func (u *UsersService) Store(name string, pass string, realName string, phone string, remark string) bool {

	pass, err := md5_encrypt.BcryptEncode(pass) // 预先处理密码加密，然后存储在数据库
	if err != nil {
		return false
	}
	return u.userModel.Store(name, pass, realName, phone, remark)
}

func (u *UsersService) Update(id int, name string, pass string, realName string, phone string, remark string) bool {
	//预先处理密码加密等操作，然后进行更新
	pass, err := md5_encrypt.BcryptEncode(pass) // 预先处理密码加密，然后存储在数据库
	if err != nil {
		return false
	}
	//如果用户新旧密码一致，直接返回true，不需要处理
	userItem, err := u.userModel.ShowOneItem(id)
	if ok := u.userModel.Update(id, name, pass, realName, phone, remark); ok {
		if userItem != nil && err == nil && userItem.Pass == pass {
			return true
		} else if userItem != nil {
			// 如果用户密码被修改，那么redis中的token值也清除
			if variable.ConfigYml.GetInt("Token.IsCacheToRedis") == 1 {
				return token2.CreateUserTokenFactory().DelTokenCacheFromRedis(int64(id))
			}
		}
	}
	return false
}

// OauthLoginToken 记录登录Token
func (u *UsersService) OauthLoginToken(userId int64, token string, expiresAt int64, clientIp string) bool {
	if ok := u.userModel.OauthLoginToken(userId, token, expiresAt, clientIp); ok {
		// 异步缓存用户有效的token到redis
		if variable.ConfigYml.GetInt("Token.IsCacheToRedis") == 1 {
			// go u.ValidTokenCacheToRedis(userId)
			user := u.FindById(userId)
			if user != nil {
				if value, err := json.Marshal(user); err == nil {
					return token2.CreateUserTokenFactory().StoreToken(userId, expiresAt, token, string(value))
				} else {
					fmt.Errorf(err.Error())
				}
			}
		}
	}
	return false
}

// Destroy 删除用户以及关联的token记录
func (u *UsersService) Destroy(id int) bool {
	u.userModel.Destroy(id)
	// 删除用户时，清除用户缓存在redis的全部token
	if variable.ConfigYml.GetInt("Token.IsCacheToRedis") == 1 {
		go token2.CreateUserTokenFactory().DelTokenCacheFromRedis(int64(id))
	}
	return true
}
