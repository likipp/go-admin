package main

import (
	"fmt"
	"go-admin/config"
	orm "go-admin/init/database"
	globalID "go-admin/init/globalID"
	initTableStruct "go-admin/init/tableStruct"
	"go-admin/router"
)

func main() {
	db, store := orm.InitMySQL(config.AdminConfig.MysqlAdmin)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	err := globalID.Init(1)
	if err != nil {
		panic("ID生成器初始化失败")
	}
	fmt.Print("store", store, "main")
	initTableStruct.InitTableStruct(db)
	_ = router.InitRouter().Run()
	//defer store.Close()
}
