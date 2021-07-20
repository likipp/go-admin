package initTableStruct

import (
	"go-admin/models"
	"gorm.io/gorm"
)

func InitTableStruct(db *gorm.DB) {
	_ = db.AutoMigrate(
		models.BaseMenu{},
		models.SysRole{},
		models.GroupKPI{},
		models.KpiData{},
		models.KPI{},
		models.SysUser{},
		models.SysDept{})
	//models.MenuResource{},
	//models.MenuMethod{})
}
