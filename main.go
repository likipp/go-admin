package main

import (
	"go-admin/config"
	orm "go-admin/init/database"
	globalID "go-admin/init/globalID"
	initTableStruct "go-admin/init/tableStruct"
	"go-admin/router"
)

func main() {
	db := orm.InitMySQL(config.AdminConfig.MysqlAdmin)
	defer db.Close()
	err := globalID.Init(1)
	if err != nil {
		panic("ID生成器初始化失败")
	}
	initTableStruct.InitTableStruct(db)
	_ = router.InitRouter().Run()
}
