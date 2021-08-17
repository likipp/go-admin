package initTableStruct

import (
	"go-admin/models"
	"gorm.io/gorm"
	"log"
)

func InitTableStruct(db *gorm.DB) {
	err := db.AutoMigrate(
		models.SysDept{},
		models.BaseMenu{},
		models.SysRole{},
		models.GroupKPI{},
		models.KpiData{},
		models.KPI{},
		models.SysUser{})
	if err != nil {
		log.Printf("AutoMigrate数据库失败%v", err)
	}
	//models.MenuResource{},
	//models.MenuMethod{})
}
