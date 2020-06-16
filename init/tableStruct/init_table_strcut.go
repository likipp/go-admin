package initTableStruct

import (
	"github.com/jinzhu/gorm"
	models "go-admin/models"
)

func InitTableStruct(db *gorm.DB) {
	db.AutoMigrate(
		models.SysUser{},
		models.SysDept{},
		models.SysRole{})
	//models.Menu{},
	//models.MenuResource{},
	//models.MenuMethod{})
}
