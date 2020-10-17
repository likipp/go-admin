package method

import (
	"github.com/jinzhu/gorm"
	"go-admin/models"
)

func PagingServer(pageParams models.KPIQueryParam, db *gorm.DB) {
	var total int
	limit := pageParams.PageSize
	offset := pageParams.PageSize * (pageParams.Current - 1)
	_ = db.Model(&models.KPI{}).Count(&total).Error
	db.Limit(limit).Offset(offset).Order("id desc")
}
