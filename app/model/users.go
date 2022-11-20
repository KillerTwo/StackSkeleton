package model

import (
	"errors"
	"go.uber.org/zap"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/md5_encrypt"
)

// 本文件针对 postgresql 数据库有效，请手动使用本文件的所有代码替换同目录的 users.go 中的所有代码即可
// 针对数据库选型为 postgresql 的开发者使用

// 操作数据库喜欢使用gorm自带语法的开发者可以参考 GinSkeleton-Admin 系统相关代码
// Admin 项目地址：https://gitee.com/daitougege/gin-skeleton-admin-backend/
// gorm_v2 提供的语法+ ginskeleton 实践 ：  http://gitee.com/daitougege/gin-skeleton-admin-backend/blob/master/app/model/button_cn_en.go

// UserModelFactory 创建 userFactory
// 参数说明： 传递空值，默认使用 配置文件选项：UseDbType（mysql）
func UserModelFactory(sqlType string) *UsersModel {
	return &UsersModel{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

type UsersModel struct {
	BaseModel
	UserName    string `gorm:"column:user_name" json:"user_name"`
	Pass        string `json:"-"`
	Phone       string `json:"phone"`
	RealName    string `gorm:"column:real_name" json:"real_name"`
	Status      int    `json:"status"`
	Token       string `json:"token"`
	LastLoginIp string `gorm:"column:last_login_ip" json:"last_login_ip"`
	Avatar      string `gorm:"avatar" json:"avatar"`
}

// TableName 表名
func (u *UsersModel) TableName() string {
	return "tb_auth_users"
}

// FindById 工具ID查询
func (u *UsersModel) FindById(userId int64) (user *UsersModel) {
	UserModelFactory("").DB.First(&user, userId)
	return user
}

// Register 用户注册（写一个最简单的使用账号、密码注册即可）
func (u *UsersModel) Register(userName, pass, userIp string) bool {
	user := UsersModel{
		UserName:    userName,
		Pass:        pass,
		LastLoginIp: userIp,
	}
	result := UserModelFactory("").DB.Create(&user)
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

// Login 用户登录,
func (u *UsersModel) Login(userName string, pass string) (*UsersModel, error) {
	result := UserModelFactory("").DB.Where("user_name=?", userName).First(u)
	if result.Error == nil {
		// 账号密码验证成功
		/*if len(u.Pass) > 0 && (u.Pass == md5_encrypt.Base64Md5(pass)) {
			return u
		}*/
		if len(u.Pass) > 0 && md5_encrypt.BcryptCompare(u.Pass, pass) {
			return u, nil
		} else {
			return nil, errors.New("用户名或密码错误")
		}
	} else {
		variable.ZapLog.Error("根据账号查询单条记录出错:", zap.Error(result.Error))
	}
	return nil, errors.New("用户名或密码错误")
}

// UpdateUserLoginInfo 更新用户登陆次数、最近一次登录ip、最近一次登录时间
func (u *UsersModel) UpdateUserLoginInfo(lastLoginIp string, userId int64) {
	UserModelFactory("").DB.Model(&UsersModel{}).Where("id=?", &userId).Updates(&UsersModel{LastLoginIp: lastLoginIp})
}

// ShowOneItem 根据用户ID查询一条信息
func (u *UsersModel) ShowOneItem(userId int) (*UsersModel, error) {
	result := UserModelFactory("").DB.First(u, userId)
	if result.Error == nil {
		return u, nil
	} else {
		return nil, result.Error
	}
}

// counts 查询数据之前统计条数
func (u *UsersModel) counts(userName string) (counts int64) {

	res := UserModelFactory("").DB.Where("user_name=?", userName).Count(&counts)

	if res.Error != nil {
		variable.ZapLog.Error("UsersModel - counts 查询数据条数出错", zap.Error(res.Error))
	}
	return counts
}

// Show 查询（根据关键词模糊查询）
func (u *UsersModel) Show(userName string, limitStart, limitItems int) (counts int64, temp []UsersModel) {
	if counts = u.counts(userName); counts > 0 {
		res := UserModelFactory("").DB.Where("user_name LIKE ? AND limit ? offset ?", "%"+userName+"%", limitStart, limitItems).Find(&temp)
		if res.RowsAffected > 0 {
			return counts, temp
		}
	}
	return 0, nil
}

// Store 新增
func (u *UsersModel) Store(userName string, pass string, realName string, phone string) bool {
	user := UsersModel{
		UserName: userName,
		Pass:     pass,
		RealName: realName,
		Phone:    phone,
	}
	result := UserModelFactory("").DB.Create(&user)
	if result.RowsAffected > 0 {
		return true
	}
	return false
}

// UpdateDataCheckUserNameIsUsed 更新前检查新的用户名是否已经存在（避免和别的账号重名）
func (u *UsersModel) UpdateDataCheckUserNameIsUsed(userId int, userName string) (exists int64) {
	UserModelFactory("").DB.Where("id <> ? AND user_name=?", userId, userName).First(&exists)
	return exists
}

// Update 更新
func (u *UsersModel) Update(id int, userName string, pass string, realName string, phone string) bool {
	result := UserModelFactory("").DB.Where("status = 1 AND id = ?", id).Updates(UsersModel{UserName: userName, Pass: pass, RealName: realName, Phone: phone})
	if result.RowsAffected >= 0 {
		return true
	}
	return false
}

// Destroy 删除用户以及关联的token记录
func (u *UsersModel) Destroy(id int) bool {
	u.DB.Delete(&UsersModel{}, id)
	return true
}
