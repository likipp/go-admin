package models

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	orm "go-admin/init/database"
	"go-admin/init/globalID"
	"go-admin/internal/schema"
)

type KPI struct {
	BaseModel
	UUID   string `json:"uuid"`
	Name   string `gorm:"column:name" json:"name"`
	Unit   string `gorm:"column:unit" json:"unit"`
	Status int    `gorm:"column:status" json:"status"`
}

type KPIQueryParam struct {
	schema.PaginationParam
	Name   string `form:"name"`
	Status int    `form:"status"`
}

func (KPI) TableName() string {
	return "kpi"
}

func (k *KPI) CreateKPI() (err error, KPI *KPI) {
	hasKPI := orm.DB.Where("name = ?", k.Name).RecordNotFound()
	if hasKPI {
		return errors.New("KPI名称重复,请检查"), nil
	} else {
		k.UUID, err = initID.GetID()
		if err != nil {
			return
		}
		err = orm.DB.Create(k).Error
	}
	return err, k
}

func PagingServer(pageParams KPIQueryParam, db *gorm.DB) {
	var total int
	limit := pageParams.PageSize
	offset := pageParams.PageSize * (pageParams.Current - 1)
	_ = db.Model(&KPI{}).Count(&total).Error
	db.Limit(limit).Offset(offset).Order("id desc")
}

func (k *KPI) GetKPIList(params KPIQueryParam) (err error, KPIList []KPI) {
	var db *gorm.DB
	if v := params.Name; v != "" {
		db = orm.DB.Where("name = ?", v).Find(&KPIList)
	}
	if v := params.Status; v > 0 {
		db = orm.DB.Where("status =?", v).Find(&KPIList)
	}
	if params.Status <= 0 && params.Name == "" {
		db = orm.DB.Find(&KPIList)
	}
	params.Pagination = true
	PagingServer(params, db)
	return err, KPIList
}

func (k *KPI) GetKPIByUUID(uuid string) (err error, KPI *KPI) {
	return nil, nil
}

func (k *KPI) DeleteKPIByUUID(uuid string) error {
	return nil
}
