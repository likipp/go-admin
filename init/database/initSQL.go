package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go-admin/config"
	"log"
)

var DB *gorm.DB

func InitMySQL(admin config.MySQL) *gorm.DB {
	if db, err := gorm.Open("mysql", admin.Username+":"+admin.Password+"@("+admin.Path+")/"+admin.DBName+"?"+admin.Config); err != nil {
		log.Printf("DEFAULTDB数据库启动异常%v", err)
	} else {
		DB = db
		DB.DB().SetMaxIdleConns(10)
		DB.DB().SetMaxOpenConns(100)
	}
	return DB
}
