package main

import (
	"go-admin/config"
	"go-admin/global"
	"go-admin/init/ccasbin"
	"go-admin/init/cookies"
	orm "go-admin/init/database"
	globalID "go-admin/init/globalID"
	initTableStruct "go-admin/init/tableStruct"
	"go-admin/router"
)

func main() {
	orm.InitMySQL(config.AdminConfig.MysqlAdmin)
	sqlDB, _ := global.GDB.DB()
	defer sqlDB.Close()
	globalID.Init(1)
	initTableStruct.InitTableStruct(global.GDB)
	cookies.InitSession(config.AdminConfig.RedisAdmin)
	ccasbin.InitCasBin()
	router.InitRouter()
}
