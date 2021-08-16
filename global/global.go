package global

import (
	"gorm.io/gorm"
)

type User struct {
	Username string `json:"username"`
	NickName string `json:"nickname" gorm:"default:'匿名用户'"`
	DeptID   string `json:"deptID"`
}

var (
	GDB   *gorm.DB
	GUser User
)
