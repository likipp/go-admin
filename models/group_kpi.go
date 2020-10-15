package models

type GroupKPI struct {
	BaseModel
	Dept   string `gorm:"column:dept"  json:"deptID"`
	KPI    string `gorm:"column:kpi"   json:"KpiID"`
	ULimit int    `gorm:"u_limit"      json:"u_limit"`
	LLimit int    `gorm:"l_limit"      json:"l_limit"`
	TLimit int    `gorm:"t_limit"      json:"t_limit"`
	Status string `gorm:"status"       json:"status"`
}

func (GroupKPI) TableName() string {
	return "group_kpi"
}
