package models

import (
	"errors"
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
	if utils.PasswordVerify(user.Password, l.Password) {
		return
	}
	return
}

func UserLogin(l *Login) (err error, userInter *SysUser) {
	var user SysUser
	err = orm.DB.Table("sys_user").Where("username = ?", l.Username).Find(&user).Error
	if err != nil {
		return errors.New("用户不存在"), &user
	}
	if utils.PasswordVerify(user.Password, l.Password) != true {
		return errors.New("密码不正确"), &user
	}
	return err, &user
}
