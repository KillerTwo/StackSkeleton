package service

import (
	"goskeleton/app/http/response"
	"goskeleton/app/model/auth"
)

func MenuServiceFactory() *MenusCrud {
	return &MenusCrud{
		menuModel: auth.CreateMenuFactory(""),
	}
}

type MenusCrud struct {
	menuModel *auth.MenuModel
}

// GetMenusByRoles 获取指定角色对应的菜单
func (m *MenusCrud) GetMenusByRoles(roleKeys []string) ([]auth.MenuModel, error) {
	menus, err := m.menuModel.FindByRoles(roleKeys)
	if err != nil {
		return nil, err
	}
	list := m.getChildPerms(menus, 0)
	return list, nil
}

// SelectTreeRoutesByRole 用户所有的路由菜单
func (m *MenusCrud) SelectTreeRoutesByRole(roleKeys []string) ([]response.Route, error) {
	menus, err := m.menuModel.FindByRoles(roleKeys)
	if err != nil {
		return nil, err
	}
	list := m.getChildPerms(menus, 0)
	return m.buildTreeRoute(list), nil
}

// SelectTreeMenus 树形菜单
func (m *MenusCrud) SelectTreeMenus() ([]auth.MenuModel, error) {
	menus, err := m.menuModel.FindAll()
	if err != nil {
		return nil, err
	}
	list := m.getChildPerms(menus, 0)
	return list, nil
}

// SelectTreeRoutes 树形菜单
func (m *MenusCrud) SelectTreeRoutes() ([]response.Route, error) {
	menus, err := m.menuModel.FindAll()
	if err != nil {
		return nil, err
	}
	list := m.getChildPerms(menus, 0)
	return m.buildTreeRoute(list), nil
}

func (m *MenusCrud) getChildPerms(menus []auth.MenuModel, parentId int64) []auth.MenuModel {
	var list []auth.MenuModel
	for i, _ := range menus {
		if menus[i].ParentId == parentId {
			// 设置menu的子菜单
			m.recursionFn(menus, &menus[i])
			list = append(list, menus[i])
		}
	}
	return list
}

func (m *MenusCrud) recursionFn(menus []auth.MenuModel, model *auth.MenuModel) {
	children := m.getChildList(menus, model)
	model.Children = children
	for i, _ := range children {
		if m.hasChild(menus, &children[i]) {
			m.recursionFn(menus, &children[i])
		}
	}
}

func (m *MenusCrud) getChildList(menus []auth.MenuModel, model *auth.MenuModel) []auth.MenuModel {
	var list []auth.MenuModel
	for i, _ := range menus {
		if menus[i].ParentId == model.Id {
			list = append(list, menus[i])
		}
	}
	return list
}

func (m *MenusCrud) hasChild(menus []auth.MenuModel, model *auth.MenuModel) bool {
	list := m.getChildList(menus, model)
	return len(list) > 0
}

func (m *MenusCrud) toRoute(menu auth.MenuModel) response.Route {
	route := &response.Route{
		Path:      menu.Path,
		Name:      menu.Name,
		Component: menu.Component,
		Redirect:  menu.Redirect,
		Meta: response.Meta{
			Title:              menu.Title,
			Icon:               menu.Icon,
			HideMenu:           menu.HideMenu,
			FrameSrc:           menu.FrameSrc,
			IgnoreKeepAlive:    !menu.IsCache,
			CurrentActiveMenu:  menu.CurrentActiveMenu,
			ShowMenu:           !menu.HideMenu,
			HideChildrenInMenu: menu.HideChildrenInMenu,
			HideBreadcrumb:     menu.HideBreadcrumb,
		},
	}
	if menu.Children != nil && len(menu.Children) > 0 {
		var children []response.Route
		for _, item := range menu.Children {
			r := m.toRoute(item)
			children = append(children, r)
		}
		route.Children = children
	}
	return *route
}

func (m *MenusCrud) buildTreeRoute(menus []auth.MenuModel) []response.Route {
	var routes []response.Route
	for _, menu := range menus {
		route := m.toRoute(menu)
		routes = append(routes, route)
	}
	return routes
}
