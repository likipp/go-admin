package models

import (
	"github.com/jinzhu/gorm"
	"go-admin/controller/server"
	orm "go-admin/init/database"
	"go-admin/models/page"
)

type SysRole struct {
	gorm.Model
	RoleId   string `json:"roleId" gorm:"not null;unique"`
	RoleName string `json:"roleName"`
	ParentId string `json:"parentId"`
	//DataRoleId []SysRole `json:"dataRoleId" gorm:"many2many:sys_data_role_id;association_jointable_foreignkey:data_id"`
	Children []SysRole  `json:"children" gorm:"many2many:children_roles;association_jointable_foreignkey:role_id"`
	Users    []*SysUser `gorm:"many2many:user_role;"`
	//UserID []string
}

func (SysRole) TableName() string {
	return "sys_role"
}

func (r *SysRole) CreateRole() (role *SysRole, err error) {
	err = orm.DB.Create(r).Error
	return r, err
}

func (r *SysRole) GetList(info page.InfoPage) (err error, list interface{}, total int64) {
	err, db, total := server.PagingServer(r, info)
	var roles []SysRole
	if err != nil {
		return
	} else {
		err := db.Preload("Users").Find(&roles).Error
		return err, roles, total
	}
}
