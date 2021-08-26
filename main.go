package main

import (
	"fmt"
	"go-admin/config"
	"go-admin/init/ccasbin"
	"go-admin/init/cookies"
	orm "go-admin/init/database"
	globalID "go-admin/init/globalID"
	"go-admin/router"
)

func main() {
	db := orm.InitMySQL(config.AdminConfig.MysqlAdmin)
	fmt.Println("DB信息:", db)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	globalID.Init(1)
	//initTableStruct.InitTableStruct(db)
	cookies.InitSession(config.AdminConfig.RedisAdmin)

	_, err := ccasbin.InitCasBin()
	if err != nil {
		fmt.Println(err, "错误信息")
		return
	}
	router.InitRouter()

	//defer store.Close()
}
