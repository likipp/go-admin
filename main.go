package main

import (
	"go-admin/config"
	"go-admin/init/cookies"
	orm "go-admin/init/database"
	globalID "go-admin/init/globalID"
	"go-admin/router"
)

func main() {
	db := orm.InitMySQL(config.AdminConfig.MysqlAdmin)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	globalID.Init(1)
	//initTableStruct.InitTableStruct(db)
	cookies.InitSession(config.AdminConfig.RedisAdmin)
	//ccasbin.InitCasBin()
	router.InitRouter()

	//defer store.Close()
}
