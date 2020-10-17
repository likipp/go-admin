package entity

import (
	"github.com/jinzhu/gorm"
)

func GetDBWithModel(db *gorm.DB, i interface{}) *gorm.DB {
	return db.Model(i)
}
