package test

import (
	"fmt"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/auth"
	"goskeleton/app/service/menus/service"
	"testing"
)

func TestSelectTreeMenu(t *testing.T) {
	/*factory := crud.CreateMenuCrudFactory()
	  menus, err := factory.SelectTreeMenus()
	  if err != nil {
	  	t.Errorf("单元测试失败，错误明细:%s\n", err)
	  }
	  fmt.Println(menus)*/
	// variable.GormDbPostgreSql.AutoMigrate(&auth.RoleModel{})
	variable.GormDbPostgreSql.AutoMigrate(&auth.MenuModel{})

}

func TestGetMenusByRoles(t *testing.T) {
	factory := service.MenuServiceFactory()
	menus, err := factory.GetMenusByRoles([]string{"ROLE_ADMIN"})
	if err != nil {
		t.Errorf("单元测试失败，错误明细:%s\n", err)
	}
	fmt.Println(menus)
}
