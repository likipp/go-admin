package initTableStruct

import (
	models "go-admin/models"
	"gorm.io/gorm"
)

func InitTableStruct(db *gorm.DB) {
	_ = db.AutoMigrate(
		models.SysUser{},
		models.SysDept{},
		models.SysRole{},
		models.GroupKPI{},
		models.KpiData{},
		models.KPI{})
	//models.Menu{},
	//models.MenuResource{},
	//models.MenuMethod{})
}
