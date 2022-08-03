package auth

import "goskeleton/app/model"

func CreateRoleFactory(sqlType string) *RoleModel {
	return &RoleModel{BaseModel: model.BaseModel{DB: model.UseDbConn(sqlType)}}
}

// RoleModel 角色模型
type RoleModel struct {
	model.BaseModel
	RoleKey  string       `gorm:"column:role_key" json:"roleKey"`
	RoleName string       `gorm:"column:role_name" json:"roleName"`
	Menus    []*MenuModel `gorm:"many2many:tb_auth_role_menus;" gorm:"comment:角色菜单关联" json:"menus"`
}

// TableName 表名
func (u *RoleModel) TableName() string {
	return "tb_auth_role"
}
