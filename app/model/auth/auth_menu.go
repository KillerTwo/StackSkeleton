package auth

import "goskeleton/app/model"

// MenuModel 菜单实体
type MenuModel struct {
	model.BaseModel
	Name               string                 `gorm:"column:name" gorm:"comment:菜单名称" json:"name"`
	ParentName         string                 `gorm:"column:parent_name" gorm:"comment:父菜单名称" json:"parentName"`
	ParentId           int64                  `gorm:"column:parent_id" gorm:"comment:父菜单ID" json:"parentId"`
	OrderNo            uint16                 `gorm:"column:order_no" gorm:"comment:父菜单ID" json:"orderNo"`
	Path               string                 `gorm:"column:path" gorm:"comment:菜单路径" json:"path"`
	Icon               string                 `gorm:"column:icon" gorm:"comment:菜单图标" json:"icon"`
	HideMenu           bool                   `gorm:"column:hide_menu" gorm:"comment:是否隐藏" json:"hideMenu"`
	Component          string                 `gorm:"column:component" gorm:"comment:前端组件" json:"component"`
	IsFrame            bool                   `gorm:"column:is_iframe" gorm:"comment:是否是外链" json:"isFrame"`
	FrameSrc           string                 `gorm:"column:frame_src" gorm:"comment:iframe链接" json:"frameSrc"`
	IsCache            bool                   `gorm:"column:is_cache" gorm:"comment:是否缓存" json:"isCache"`
	MenuType           string                 `gorm:"column:menu_type" gorm:"comment:菜单类型" json:"menuType"`
	Title              string                 `gorm:"column:title" gorm:"comment:菜单标题" json:"title"`
	Redirect           string                 `gorm:"column:redirect" gorm:"comment:重定向地址" json:"redirect"`
	HideChildrenInMenu bool                   `gorm:"column:hide_children_in_menu" gorm:"comment:是否隐藏子菜单" json:"hideChildrenInMenu"`
	CurrentActiveMenu  string                 `gorm:"column:current_active_menu" gorm:"comment:当前激活的菜单" json:"currentActiveMenu"`
	HideBreadcrumb     bool                   `gorm:"column:hide_breadcrumb" gorm:"comment:在面包屑不显示" json:"hideBreadcrumb"`
	Children           []MenuModel            `gorm:"-" json:"children"`
	Meta               map[string]interface{} `gorm:"-" json:"meta"`
	Roles              []*RoleModel           `gorm:"many2many:tb_auth_role_menus;" gorm:"comment:角色菜单关联" json:"roles"`
}

func CreateMenuFactory(sqlType string) *MenuModel {
	return &MenuModel{BaseModel: model.BaseModel{DB: model.UseDbConn(sqlType)}}
}

// TableName 表名
func (u *MenuModel) TableName() string {
	return "tb_auth_menu"
}

// Create 新增菜单
func (m *MenuModel) Create(menuModel *MenuModel) (*MenuModel, error) {
	result := m.DB.Create(menuModel)
	if result.Error != nil {
		return nil, result.Error
	}
	return menuModel, nil
}

func (m *MenuModel) FindAll() (menus []MenuModel, err error) {
	result := m.DB.Order("order_no").Find(&menus)
	if result.Error != nil {
		return nil, result.Error
	}
	return menus, nil
}

func (m *MenuModel) FindByRoles(roleKeys []string) (menus []MenuModel, err error) {
	// result := m.DB.Order("order_no").Find(&menus)
	result := m.DB.Order("order_no").Joins("JOIN tb_auth_role_menus b ON tb_auth_menu.id=b.menu_model_id").
		Joins("JOIN tb_auth_role c ON b.role_model_id=c.id").Where("c.role_key IN ?", roleKeys).Find(&menus)

	if result.Error != nil {
		return nil, result.Error
	}
	return menus, nil
}

/*func (m *MenuModel) findByUserId(userId uint64) (menus []MenuModel, err error) {
	m.DB.Where("")
}*/
