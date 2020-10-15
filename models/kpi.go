package models

import (
	"github.com/pkg/errors"
	orm "go-admin/init/database"
	"go-admin/init/globalID"
)

type KPI struct {
	BaseModel
	UUID   string `json:"uuid"`
	Name   string `gorm:"column:name" json:"name"`
	Unit   string `gorm:"column:unit" json:"unit"`
	Status string `gorm:"column:status" json:"status"`
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
