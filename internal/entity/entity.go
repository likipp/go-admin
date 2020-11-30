package entity

import (
	"gorm.io/gorm"
)

func GetDBWithModel(db *gorm.DB, i interface{}) *gorm.DB {
	return db.Model(i)
}
