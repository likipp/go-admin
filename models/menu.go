package models

import (
	orm "go-admin/init/database"
)

type BaseMenu struct {
	BaseModel
	Icon       string  `json:"icon" gorm:"column:icon"`
	MenuLevel  uint    `json:"-"`
	Sequence   int     `json:"sequence" gorm:"column:sequence;index;default:0;not null;"`
	Router     string  `json:"router" gorm:"column:router;size:255;" binding:"required"`
	ParentId   string  `json:"parent_id" gorm:"comment:父菜单ID;index"`
	ParentPath string  `json:"parent_path" gorm:"comment:路由path;index"`
	Name       string  `json:"name" gorm:"comment:路由name;index" binding:"required"`
	ShowStatus int     `json:"show_status" gorm:"show_status;index;default:0;not null" binding:"required"`
	Hidden     bool    `json:"hidden" gorm:"comment:是否在列表隐藏;index" binding:"required,max=2,min=1" binding:"required""`
	Status     int     `json:"status" gorm:"column:status;index;default:0;not null;" binding:"required,max=2,min=1"`
	Memo       *string `json:"memo" gorm:"column:memo;size:1024;"`
	Component  string  `json:"component" gorm:"comment:对应前端文件路径;index"`
	//SysRoles     []SysRole     `json:"roles" gorm:"many2many:sys_authority_menus;"`
}

func (BaseMenu) TableName() string {
	return "base_menu"
}

func (m *BaseMenu) CreateBaseMenu() (err error, menu *BaseMenu) {
	//var menu BaseMenu
	err = orm.DB.Create(menu).Error
	return err, menu
}
