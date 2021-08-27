package global

import (
	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"
)

type user struct {
	Username string `json:"username"`
	NickName string `json:"nickname" gorm:"default:'匿名用户'"`
	DeptID   uint   `json:"deptID"`
}

var (
	GDB             *gorm.DB
	GUser           user
	GSyncedEnforcer *casbin.SyncedEnforcer
)
