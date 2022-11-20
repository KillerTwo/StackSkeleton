package users

import commondatatype "goskeleton/app/http/request/common/data_type"

type BaseField struct {
	UserName string `form:"user_name" json:"user_name"  binding:"required,min=1"` // 必填、对于文本,表示它的长度>=1
	Pass     string `form:"pass" json:"pass" binding:"required,min=6,max=20"`     //  密码为 必填，长度>=6
}

type Id struct {
	Id float64 `form:"id"  json:"id" binding:"required,min=1"`
}

type Register struct {
	UserId   int64  `form:"userId" json:"userId" binding:"-"`
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=20"`
	RealName string `form:"realName" json:"realName" binding:"-"`
	Phone    string `form:"phone" json:"phone" binding:"-"`
	Remark   string `form:"remark" json:"remark" binding:"-"`
}

type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Phone    string `form:"phone" json:"phone" binding:"-"`
}

type UserParam struct {
	commondatatype.Page
	Username string `form:"username" json:"username" binding:"-"`
}
