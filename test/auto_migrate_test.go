package test

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/model/auth"
	"testing"
)

func TestAutoMigrate(t *testing.T) {
	variable.GormDbPostgreSql.AutoMigrate(&model.UsersModel{}, &auth.CasbinRule{}, &auth.MenuModel{}, &auth.RoleModel{})
}
