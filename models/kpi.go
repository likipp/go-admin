package models

import "time"

type KPI struct {
	BaseModel
	KpiId  string    `gorm:"column:kpi_id" json:"kpiID"`
	Name   string    `gorm:"column:name" json:"name"`
	Unit   string    `gorm:"column:unit" json:"unit"`
	Status string    `gorm:"column:status" json:"status"`
	InTime time.Time `gorm:"column:in_time" json:"inTime"`
	MoTime time.Time `gorm:"column:mo_time" json:"moTime"`
}

type GroupKPI struct {
	BaseModel
	Dept   string `gorm:"column:dept"  json:"deptID"`
	KPI    string `gorm:"column:kpi"   json:"KpiID"`
	ULimit int    `gorm:"u_limit"      json:"u_limit"`
	LLimit int    `gorm:"l_limit"      json:"l_limit"`
	TLimit int    `gorm:"t_limit"      json:"t_limit"`
	Status string `gorm:"status"       json:"status"`
}

type KpiInput struct {
	BaseModel
	RValue   int        `gorm:"column:r_value" json:"r_value"`
	Month    time.Month `gorm:"column: month"  json:"month"`
	User     string     `gorm:"user"           json:"user"`
	GroupKPI GroupKPI   `gorm:"group_kpi"      json:"group_kpi"`
}
