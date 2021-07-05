package models

import (
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	orm "go-admin/init/database"
	initID "go-admin/init/globalID"
)

type BaseMenu struct {
	BaseModel
	//gorm.Model
	//gorm:"size:256;not null;unique"
	UUID       string `json:"uuid"`
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

type MenuTrees []*MenuTree

type MenuTree struct {
	UUID       string     `yaml:"-" json:"uuid"`
	Name       string     `yaml:"name" json:"name"`
	Icon       string     `yaml:"icon" json:"icon"`
	Path       string     `yaml:"path,omitempty" json:"path"`
	ParentID   string     `yaml:"-" json:"parent_id"`
	ParentPath string     `yaml:"-" json:"parent_path"`
	Sequence   int        `yaml:"sequence" json:"sequence"`
	ShowStatus int        `yaml:"-" json:"show_status"`
	Status     int        `yaml:"-" json:"status"`
	Children   *MenuTrees `yaml:"children,omitempty" json:"children,omitempty"`
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

func MenuToTree(m MenuTrees) MenuTrees {
	mi := make(map[string]*MenuTree)
	for _, item := range m {
		mi[item.UUID] = item
	}

	var list MenuTrees
	for _, item := range m {
		if item.ParentID == "" {
			list = append(list, item)
			continue
		}
		if pitem, ok := mi[item.ParentID]; ok {
			if pitem.Children == nil {
				children := MenuTrees{item}
				pitem.Children = &children
				continue
			}
			*pitem.Children = append(*pitem.Children, item)
		}
	}
	return list
}

func MenusAddChildren(ms []BaseMenu) (err error, mts MenuTrees) {
	var mt MenuTree
	for _, item := range ms {
		err = copier.Copy(&item, &mt)
		if err != nil {
			return errors.New("转换失败"), nil
		}
		mts = append(mts, &mt)
	}
	fmt.Println(mts, "菜单列表")
	return nil, mts
}

func (m *BaseMenu) GetBaseMenu() (err error, trees MenuTree) {
	var menus []BaseMenu
	err = orm.DB.Find(&menus).Error
	if err != nil {
		return err, trees
	}
	err, list := MenusAddChildren(menus)
	if err != nil {
		return err, MenuTree{}
	}
	list2 := MenuToTree(list)
	fmt.Println(list2, "菜单列表转换后")
	return err, trees
}
