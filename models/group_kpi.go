package models

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	orm "go-admin/init/database"
	initID "go-admin/init/globalID"
	"go-admin/internal/entity"
	"go-admin/internal/schema"
)

type GroupKPI struct {
	BaseModel
	UUID   string `gorm:"column:uuid"         json:"uuid"`
	Dept   string `gorm:"column:dept"         json:"deptID"`
	KPI    string `gorm:"column:kpi"          json:"KpiID"`
	ULimit int    `gorm:"column:u_limit"      json:"u_limit"`
	LLimit int    `gorm:"column:l_limit"      json:"l_limit"`
	TLimit int    `gorm:"column:t_limit"      json:"t_limit"`
	Status string `gorm:"column:status"       json:"status"`
}

type GroupKpiQueryParam struct {
	schema.PaginationParam
	Dept   string `form:"dept"`
	KPI    string `form:"kpi"`
	Status string `form:"status"`
}

type GroupKPIWithName struct {
	UUID     string `json:"uuid"`
	Dept     string `json:"deptID"`
	DeptName string `json:"deptName"`
	KPI      string `json:"KpiID"`
	KPIName  string `json:"KPIName"`
	ULimit   int    `json:"u_limit"`
	LLimit   int    `json:"l_limit"`
	TLimit   int    `json:"t_limit"`
	Status   string `json:"status"`
}

func (GroupKPI) TableName() string {
	return "group_kpi"
}

func GroupKpiPagingServer(pageParams GroupKpiQueryParam, db *gorm.DB) {
	var total int
	limit := pageParams.PageSize
	offset := pageParams.PageSize * (pageParams.Current - 1)
	_ = db.Count(&total).Error
	db.Limit(limit).Offset(offset).Order("id desc")
}

func GetGroupKpiDB(db *gorm.DB) *gorm.DB {
	return entity.GetDBWithModel(db, new(GroupKPI))
}

func (g *GroupKPI) CreateGroupKPI() (err error, gK *GroupKPI) {
	var result GroupKPI
	db := GetGroupKpiDB(orm.DB)
	hasGroupKpi := db.Where("dept = ? AND kpi = ?", g.Dept, g.KPI).First(&result).RecordNotFound()
	if !hasGroupKpi {
		return errors.New("部门KPI已经关联"), nil
	}
	fmt.Println(g, "前端数据")
	g.UUID, err = initID.GetID()
	if err != nil {
		return
	}
	err = orm.DB.Create(g).Error
	if err != nil {
		return err, g
	}
	return nil, g
}

func (g *GroupKPI) GetGroupKPI() (err error, gk []GroupKPIWithName) {
	var results []GroupKPI
	var resultsWithName []GroupKPIWithName
	db := GetGroupKpiDB(orm.DB)
	db.Find(&results)
	for _, v := range results {
		var dept SysDept
		var kpi KPI
		orm.DB.Where("dept_id  = ?", v.Dept).First(&dept)
		orm.DB.Where("uuid = ?", v.KPI).First(&kpi)
		resultWithName := GroupKPIWithName{
			UUID:     v.UUID,
			Dept:     v.Dept,
			DeptName: dept.DeptName,
			KPI:      v.KPI,
			KPIName:  kpi.Name,
			ULimit:   v.ULimit,
			LLimit:   v.LLimit,
			TLimit:   v.TLimit,
			Status:   v.Status,
		}
		resultsWithName = append(resultsWithName, resultWithName)
	}
	return nil, resultsWithName
}
