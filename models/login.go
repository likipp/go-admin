package models

import (
	"errors"
	"fmt"
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

func UserLogin(l *Login) (err error, userInter *SysUser) {
	var user SysUser
	err = orm.DB.Table("sys_user").Where("username = ?", l.Username).Find(&user).Error
	if err != nil {
		return errors.New("用户不存在"), &user
	}
	fmt.Println(l.Password, "ll", user.Password)
	//err = bcrypt.CompareHashAndPassword([]byte("d41d8cd98f00b204e9800998ecf8427e"), []byte("d41d8cd98f00b204e9800998ecf8427e"))
	hash, _ := utils.PasswordHash(l.Password)

	match := utils.PasswordVerify(user.Password, hash)
	fmt.Println("验证:", match)
	if err != nil {
		fmt.Println(err, "err")
		return errors.New("密码不正确"), &user
	}
	//if user.Password != l.Password {
	//	return errors.New("密码不正确"), &user
	//}
	return err, &user
}
