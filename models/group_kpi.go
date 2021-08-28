package models

import (
	"errors"
	"go-admin/global"
	orm "go-admin/init/database"
	initID "go-admin/init/globalID"
	"go-admin/internal/entity"
	"go-admin/internal/schema"
	"gorm.io/gorm"
	"time"
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
}

type DeptKPIResult struct {
	UUID     string `json:"kpi"`
	Name     string `json:"name"`
	DeptName string `json:"dept_name"`
	DeptID   string `json:"dept_id"`
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
	g.UUID, err = initID.NewID()
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
	db := GetGroupKpiDB(global.GDB)
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

func (g *GroupKPI) GetGroupKPIDept(params KPIDeptQueryParam) (err error, result []DeptKPIResult) {
	var month = time.Now().Format("2006-01")
	var selectDept = "sys_dept.dept_id, sys_dept.dept_name"
	var joinQuery = "join sys_dept on sys_dept.dept_id = group_kpi.dept join kpi on group_kpi.kpi = kpi.uuid join kpi_data on group_kpi.uuid = kpi_data.group_kpi"
	db := GetGroupKpiDB(orm.DB).Distinct("sys_dept.dept_id, sys_dept.dept_name").Joins(joinQuery)
	if params.Dept == "" && params.KPI == "" {
		err := db.Distinct("sys_dept.dept_id, sys_dept.dept_name").Select(selectDept).Scan(&result).Error
		if err != nil {
			return err, result
		}
	}
	if params.Dept != "" {
		err := db.Distinct("sys_dept.dept_id, sys_dept.dept_name").Select("kpi.uuid, kpi.name").Where("group_kpi.dept = ? and kpi_data.in_time < ?", params.Dept, month).Scan(&result).Error
		if err != nil {
			return err, result
		}
	}
	return nil, result
}
