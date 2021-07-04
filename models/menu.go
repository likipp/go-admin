package models

import (
	orm "go-admin/init/database"
	initID "go-admin/init/globalID"
)

type BaseMenu struct {
	BaseModel
	//gorm.Model
	//gorm:"size:256;not null;unique"
	UUID       string `json:"uuid" `
	Icon       string `json:"icon" gorm:"column:icon"`
	MenuLevel  uint   `json:"menu_level"`
	Sequence   int    `json:"sequence" gorm:"column:sequence;index;default:0;not null;"`
	Path       string `json:"path" gorm:"column:path;"`
	ParentId   string `json:"parent_id" gorm:"comment:父菜单ID;index"`
	ParentPath string `json:"parent_path" gorm:"comment:路由path;index"`
	//Routers    []BaseMenu `json:"routers" gorm:"foreignKey:UUID"`
	Name       string  `json:"name" gorm:"comment:路由name;index"`
	ShowStatus int     `json:"show_status" gorm:"show_status;index;default:0;not null"`
	Hidden     bool    `json:"hidden" gorm:"comment:是否在列表隐藏;index"`
	Status     int     `json:"status" gorm:"column:status;index;default:0;not null;"`
	Memo       *string `json:"memo" gorm:"column:memo;size:1024;"`
	Component  string  `json:"component" gorm:"comment:对应前端文件路径;index"`
	//SysRoles     []SysRole     `json:"roles" gorm:"many2many:sys_authority_menus;"`
}

func (BaseMenu) TableName() string {
	return "base_menu"
}

func (m *BaseMenu) CreateBaseMenu() (err error, menu *BaseMenu) {
	m.UUID, err = initID.NewID()
	if err != nil {
		return err, m
	}
	err = orm.DB.Create(m).Error
	return err, m
}
