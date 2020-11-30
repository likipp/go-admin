package server

import (
	orm "go-admin/init/database"
	"go-admin/models/page"
	"gorm.io/gorm"
)

func PagingServer(paging page.Paging, infoPage page.InfoPage) (err error, db *gorm.DB, total int64) {
	limit := infoPage.PageSize
	offset := infoPage.PageSize * (infoPage.Page - 1)
	err = orm.DB.Model(paging).Count(&total).Error
	db = orm.DB.Limit(limit).Offset(offset).Order("id desc")
	return err, db, total
}
