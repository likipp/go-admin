package main

import (
	"go-admin/config"
	"go-admin/init/cookies"
	orm "go-admin/init/database"
	globalID "go-admin/init/globalID"
	initTableStruct "go-admin/init/tableStruct"
	"go-admin/router"
)

func main() {
	db := orm.InitMySQL(config.AdminConfig.MysqlAdmin)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	globalID.Init(1)
	cookies.InitSession(config.AdminConfig.RedisAdmin)
	initTableStruct.InitTableStruct(db)
	_ = router.InitRouter().Run()
	//defer store.Close()
}
