package models

import (
	"fmt"
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
	Name   string `from:"name"`
	Status int    `from:"status"`
}

func (KPI) TableName() string {
	return "kpi"
}

func (k *KPI) CreateKPI() (err error, KPI *KPI) {
	hasKPI := orm.DB.Where("name = ?", k.Name).RecordNotFound()
	if hasKPI {
		return errors.New("用户名已经注册"), nil
	} else {
		k.UUID, err = initID.GetID()
		if err != nil {
			return
		}
		err = orm.DB.Create(k).Error
	}
	return err, k
}

func (k *KPI) GetKPIList(params KPIQueryParam) (err error, KPIList []*KPI) {
	var db *gorm.DB
	fmt.Println(db, db)
	if v := params.Name; v != "" {
		db = orm.DB.Where("name = ?", v)
	}
	if v := params.Status; v > 0 {
		db = orm.DB.Where("status =?", v)
	}
	//var list []KPI
	//schema.QueryPaging(params, list, db)
	return nil, nil
}

func (k *KPI) GetKPIByUUID(uuid string) (err error, KPI *KPI) {
	return nil, nil
}

func (k *KPI) DeleteKPIByUUID(uuid string) error {
	return nil
}
