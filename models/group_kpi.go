package models

import (
	"errors"
	"fmt"
	orm "go-admin/init/database"
	initID "go-admin/init/globalID"
	"go-admin/internal/entity"
	"go-admin/internal/schema"
	"gorm.io/gorm"
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

type Dept struct {
	DeptName string `json:"dept_name"`
	DeptID   string `json:"dept_id"`
	KPI      string `json:"kpi"`
}

type KPIDeptQueryParam struct {
	Dept string `form:"dept"`
	KPI  string `form:"kpi"`
}

func (GroupKPI) TableName() string {
	return "group_kpi"
}

func GroupKpiPagingServer(pageParams GroupKpiQueryParam, db *gorm.DB) {
	var total int64
	limit := pageParams.PageSize
	offset := pageParams.PageSize * (pageParams.Current - 1)
	_ = db.Count(&total).Error
	db.Limit(int(limit)).Offset(int(offset)).Order("id desc")
}

func GetGroupKpiDB(db *gorm.DB) *gorm.DB {
	return entity.GetDBWithModel(db, new(GroupKPI))
}

func (g *GroupKPI) CreateGroupKPI() (err error, gK *GroupKPI) {
	var result GroupKPI
	db := GetGroupKpiDB(orm.DB)
	hasGroupKpi := db.Where("dept = ? AND kpi = ?", g.Dept, g.KPI).First(&result).Error
	hasGroupKpiResult := errors.Is(hasGroupKpi, gorm.ErrRecordNotFound)
	if !hasGroupKpiResult {
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

func (g *GroupKPI) GetGroupKPIDept(params KPIDeptQueryParam) (err error, dept []Dept) {
	var result []Dept
	db := GetGroupKpiDB(orm.DB).Distinct("sys_dept.dept_id, sys_dept.dept_name").Select("sys_dept.dept_id, sys_dept.dept_name, group_kpi.kpi").Joins("join sys_dept on sys_dept.dept_id = group_kpi.dept")
	if params.Dept == "" && params.KPI == "" {
		db = db.Scan(&result)
	}
	if v := params.Dept; v != "" {
		db = db.Where("group_kpi.dept = ?", v).Scan(&result)
	}

	if v := params.KPI; v != "" {
		db = db.Where("group_kpi.kpi = ?", v).Scan(&result)
	}
	//_ = GetGroupKpiDB(orm.DB).Distinct("sys_dept.dept_id, sys_dept.dept_name").Select("sys_dept.dept_id, sys_dept.dept_name").Joins("join sys_dept on sys_dept.dept_id = group_kpi.dept").Scan(&result)
	return nil, result
}
