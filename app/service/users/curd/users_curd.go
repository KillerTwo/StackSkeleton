package curd

import (
	"goskeleton/app/http/response"
	"goskeleton/app/model"
	"goskeleton/app/utils/md5_encrypt"
	"strconv"
)

func CreateUserCurdFactory() *UsersCurd {
	return &UsersCurd{model.CreateUserFactory("")}
}

type UsersCurd struct {
	userModel *model.UsersModel
}

func (u *UsersCurd) FindById(userId int64) *response.UserResponse {
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

func (u *UsersCurd) Register(userName, pass, userIp string) bool {
	pass, err := md5_encrypt.BcryptEncode(pass) // 预先处理密码加密，然后存储在数据库
	if err != nil {
		return false
	}
	return u.userModel.Register(userName, pass, userIp)
}

func (u *UsersCurd) Store(name string, pass string, realName string, phone string, remark string) bool {

	pass, err := md5_encrypt.BcryptEncode(pass) // 预先处理密码加密，然后存储在数据库
	if err != nil {
		return false
	}
	return u.userModel.Store(name, pass, realName, phone, remark)
}

func (u *UsersCurd) Update(id int, name string, pass string, realName string, phone string, remark string, clientIp string) bool {
	//预先处理密码加密等操作，然后进行更新
	pass, err := md5_encrypt.BcryptEncode(pass) // 预先处理密码加密，然后存储在数据库
	if err != nil {
		return false
	}
	return u.userModel.Update(id, name, pass, realName, phone, remark, clientIp)
}
