package database

import (
	"fmt"
	"github.com/wader/gormstore/v2"
	"go-admin/config"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var DB *gorm.DB
var Store *gormstore.Store

func InitMySQL(admin config.MySQL) (*gorm.DB, *gormstore.Store) {
	//var store *gormstore.Store
	if db, err := gorm.Open(mysql.Open(admin.Username+":"+admin.Password+"@("+admin.Path+")/"+admin.DBName+"?"+admin.Config), &gorm.Config{}); err != nil {
		log.Printf("DEFAULTDB数据库启动异常%v", err)
	} else {
		store := gormstore.New(db, []byte("secret"))
		store.SessionOpts.Secure = true
		store.SessionOpts.HttpOnly = true
		store.SessionOpts.MaxAge = 60 * 60 * 24 * 60
		quit := make(chan struct{})
		go store.PeriodicCleanup(1*time.Hour, quit)
		DB = db
		Store = store
		sqlDb, _ := DB.DB()
		sqlDb.SetMaxIdleConns(10)
		sqlDb.SetMaxOpenConns(100)
	}
	fmt.Println(Store, "store")
	return DB, Store
}
