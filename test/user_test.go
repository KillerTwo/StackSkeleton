package test

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"testing"
)

func TestUser(t *testing.T) {
	variable.GormDbPostgreSql.AutoMigrate(&model.UsersModel{})
}
