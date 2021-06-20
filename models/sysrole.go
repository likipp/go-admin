package models

import (
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"go-admin/controller/server"
	orm "go-admin/init/database"
	"go-admin/models/page"
	"gorm.io/gorm"
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

type RoleList struct {
	ID       uint   `form:"id" json:"id"`
	RoleName string `form:"roleName" json:"roleName"`
	Members  int    `form:"members" json:"members"`
}

type RoleQuery struct {
	Base     bool `form:"base"`
	Perm     bool `form:"perm"`
	Members  bool `form:"members"`
	Page     int  `form:"current"`
	PageSize int  `form:"pageSize"`
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
	var roleList []RoleList
	var role RoleList
	if err != nil {
		return
	}
	err = db.Preload("Users").Find(&roles).Error
	if err != nil {
		return
	}
	for _, v := range roles {
		err := copier.Copy(&role, &v)
		if err != nil {
			return err, nil, 0
		}
		role.Members = len(v.Users)
		roleList = append(roleList, role)
	}
	return err, roleList, total
}

func (r *SysRole) GetRoleByQuery(rq RoleQuery) (result interface{}, total int, err error) {
	var role SysRole
	var users []SysUser
	//var db *gorm.DB
	if err := orm.DB.Where("id = ?", r.ID).Find(&role).Error; err != nil {
		return result, 0, errors.New("找不到该角色")
	}

	if rq.Base {
		orm.DB.Select("roleId, roleName").Where("ID = ?", r.ID).Find(role)
	}
	if rq.Members {
		orm.DB.Select("Users").Preload("Users").Where("ID = ?", r.ID).Find(&role)
		err := orm.DB.Model(&role).Association("Users").Find(&users)
		fmt.Println(&users)
		if err != nil {
			return users, len(users), err
		}
	}
	return users, 0, nil
}
