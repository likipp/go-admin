package models

import "time"

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
