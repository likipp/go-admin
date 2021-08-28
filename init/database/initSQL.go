package database

import (
	"go-admin/config"
	"go-admin/global"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitMySQL(admin config.MySQL) {
	if db, err := gorm.Open(mysql.Open(admin.Username+":"+admin.Password+"@("+admin.Path+")/"+admin.DBName+"?"+admin.Config), &gorm.Config{}); err != nil {
		log.Printf("DEFAULTDB数据库启动异常%v", err)
	} else {
		//DB = db
		global.GDB = db
		sqlDb, _ := global.GDB.DB()
		sqlDb.SetMaxIdleConns(10)
		sqlDb.SetMaxOpenConns(100)
	}
}
