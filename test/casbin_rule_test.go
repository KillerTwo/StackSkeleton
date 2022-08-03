package test

import (
	"fmt"
	"goskeleton/app/global/variable"
	"testing"
)

// 测试是否有权限
func TestAuthorization(t *testing.T) {
	isPass, err := variable.Enforcer.Enforce("Alice", "/roles", "*")
	if err != nil {
		t.Errorf("单元测试失败，错误明细:%s\n", err)
	}
	if !isPass {
		t.Errorf("验证不通过:\n")
	} else {
		fmt.Printf("验证通过")
	}
}

// 测试获取用户角色
func TestGetRoleByUser(t *testing.T) {
	roles, err := variable.Enforcer.GetRolesForUser("Alice")
	if err != nil {
		t.Errorf("查询失败")
	}
	fmt.Printf("%q\n", roles)
}

// 测试分配角色
func TestAddAuthorize(t *testing.T) {
	success, err := variable.Enforcer.AddGroupingPolicy("Alice", "ROLE_EMPLOYEE")
	if err != nil || !success {
		t.Errorf(err.Error())
	}
}

// 测试分配菜单
func TestAddAuthorizeForMenu(t *testing.T) {
	success, err := variable.Enforcer.AddPolicy("ROLE_EMPLOYEE", "/menus", "GET")
	if err != nil || !success {
		t.Errorf(err.Error())
	}
}
