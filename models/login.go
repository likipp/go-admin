package models

import (
	orm "go-admin/init/database"
	"go-admin/utils"
	"log"
)

type Login struct {
	Username string `form:"UserName" json:"username" binding:"required"`
	Password string `form:"Password" json:"password" binding:"required"`
}

func (l *Login) GetUser() (user SysUser, role SysRole, err error) {
	err = orm.DB.Table("sys_user").Where("username = ?", l.Username).Find(&user).Error
	if err != nil {
		log.Println(err)
		return
	}
	if utils.MD5V([]byte(user.Password)) == l.Password {
		return
	}
	//_, err = util.CompareHashAndPassword(user.Password, l.Password)
	//if err != nil {
	//	return
	//}
	return
}
