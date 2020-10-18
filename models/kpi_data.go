package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	orm "go-admin/init/database"
	"go-admin/internal/entity"
	"time"
)

type KpiData struct {
	BaseModel
	RValue   int        `gorm:"column:r_value" json:"r_value"`
	Month    time.Month `gorm:"column: month"  json:"month"`
	User     string     `gorm:"user"           json:"user"`
	GroupKPI string     `gorm:"group_kpi"      json:"group_kpi"`
}

func (KpiData) TableName() string {
	return "kpi_data"
}

func GetKpiDataDB(db *gorm.DB) *gorm.DB {
	return entity.GetDBWithModel(db, new(KpiData))
}

func (k *KpiData) CreateKpiData() (err error, kd *KpiData) {
	var result KpiData
	db := GetKpiDataDB(orm.DB)
	hasKpiData := db.Where("group_kpi = ?", k.GroupKPI).First(&result).RecordNotFound()
	if !hasKpiData {
		return errors.New("KPI数据已经录入"), kd
	}
	err = db.Create(&k).Error
	if err != nil {
		return errors.New("创建KPI数据失败"), kd
	}
	return nil, k

}
