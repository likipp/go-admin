package models

import (
	"errors"
	"go-admin/global"
	"go-admin/utils"
)

type Login struct {
	Username string `form:"UserName" json:"username" binding:"required"`
	Password string `form:"Password" json:"password" binding:"required"`
}

func (l *Login) GetUser() (user SysUser, role SysRole, err error) {
	err = global.GDB.Table("sys_user").Where("username = ?", l.Username).Find(&user).Error
	if err != nil {
		return user, role, errors.New("查找用户失败")
	}
	if utils.PasswordVerify(user.Password, l.Password) {
		return user, role, errors.New("密码不匹配")
	}
	return
}

func UserLogin(l *Login) (err error, userInter *SysUser) {
	var user SysUser
	err = global.GDB.Table("sys_user").Where("username = ?", l.Username).Find(&user).Error
	if err != nil {
		return errors.New("用户不存在"), &user
	}
	if utils.PasswordVerify(user.Password, l.Password) != true {
		return errors.New("密码不正确"), &user
	}
	err = global.GDB.Model(&user).Association("Roles").Find(&user.Roles)
	if err != nil {
		return errors.New("查找关联角色失败"), nil
	}
	return err, &user
}
