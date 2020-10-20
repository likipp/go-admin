package models

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	orm "go-admin/init/database"
	"go-admin/internal/entity"
	"go-admin/internal/schema"
)

type KpiData struct {
	BaseModel
	UUID     string `gorm:"column:uuid"       json:"uuid"`
	RValue   int    `gorm:"column:r_value"    json:"r_value"`
	InTime   string `gorm:"column:in_time"    json:"in_time"`
	User     string `gorm:"column:user"       json:"user"`
	GroupKPI string `gorm:"column:group_kpi"  json:"group_kpi"`
}

type KpiDataResult struct {
	Dept    string
	KpiData map[string]interface{}
}

type KpiChild struct {
	RValue int `json:"r_value"`
	ULimit int `json:"u_limit"`
	LLimit int `json:"l_limit"`
	TLimit int `json:"t_limit"`
	//Unit   string `json:"unit"`
	InTime string `json:"in_time"`
	User   string `json:"user"`
}

type Result struct {
	Dept   string `json:"dept"`
	KPI    string `json:"kpi"`
	ULimit int    `json:"u_limit"`
	LLimit int    `json:"l_limit"`
	TLimit int    `json:"t_limit"`
	RValue int    `json:"r_value"`
	Unit   string `json:"unit"`
	User   string `json:"user"`
	//InTime   string `gorm:"column:in_time"    json:"in_time"`
	//GroupKPI string `gorm:"column:group_kpi"  json:"group_kpi"`
}

type KpiDataQueryParam struct {
	schema.PaginationParam
	GroupKPI string `form:"group_kpi"`
	User     string `form:"user"`
	Dept     string `form:"dept"`
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
	hasKpiData := db.Where("group_kpi = ? AND in_time = ?", k.GroupKPI, k.InTime).First(&result).RecordNotFound()
	if !hasKpiData {
		return errors.New("KPI数据已经录入"), kd
	}
	fmt.Println(k, "k 数据")
	err = db.Create(&k).Error
	if err != nil {
		return errors.New("创建KPI数据失败"), kd
	}
	return nil, k

}

func KpdDataPagingServer(pageParams KpiDataQueryParam, db *gorm.DB) {
	var total int
	limit := pageParams.PageSize
	offset := pageParams.PageSize * (pageParams.Current - 1)
	_ = db.Count(&total).Error
	fmt.Println(total, "total")
	db.Limit(limit).Offset(offset).Order("in_time desc")
}

func GroupByDept(kd []Result) {
	//var result []*KpiData
	//var t = make(map[string]interface{})
	for i, v := range kd {
		fmt.Println(i, "数据", v)
		//t[v.Dept]

	}
}

func (k *KpiData) GetKpiData(params KpiDataQueryParam) (err error, kd []Result) {
	db := GetKpiDataDB(orm.DB)
	//var selectData = "group_kpi.dept, group_kpi.kpi, group_kpi.l_limit, group_kpi.t_limit, group_kpi.u_limit, kpi_data.r_value, kpi.unit, kpi_data.user, kpi_data.in_time"
	//var joinData = "join group_kpi on kpi_data.group_kpi = group_kpi.uuid join kpi on group_kpi.kpi = kpi.uuid"
	//var orderData = "group_kpi.dept desc, group_kpi.kpi, kpi_data.in_time"
	var selectData = "g.dept, g.kpi, g.l_limit, g.t_limit, g.u_limit, kpi_data.r_value, k.unit, kpi_data.user, kpi_data.in_time"
	var joinData = "join group_kpi g on kpi_data.group_kpi = g.uuid join kpi k on g.kpi = k.uuid"
	var orderData = "g.dept desc, g.kpi, kpi_data.in_time"
	var result []Result
	if v := params.Dept; v != "" {
		db = db.Select(selectData).Joins(joinData).Order(orderData).Where("g.dept = ?", v).Scan(&result)
	}
	if v := params.GroupKPI; v != "" {
		db = db.Select(selectData).Joins(joinData).Order(orderData).Where("g.kpi = ?", v).Scan(&result)
	}
	if v := params.User; v != "" {
		db = db.Select(selectData).Joins(joinData).Order(orderData).Where("kpi_data.user = ?", v).Scan(&result)
	}
	if params.GroupKPI == "" && params.User == "" && params.Dept == "" {
		db = db.Select(selectData).Joins(joinData).Order(orderData).Scan(&result)
	}
	//select g.dept, g.kpi, g.l_limit, g.t_limit, g.u_limit, kpi_data.r_value, kpi.unit, kpi_data.user, kpi_data.in_time from kpi_data
	//	join group_kpi g on kpi_data.group_kpi = g.uuid
	//	join kpi on g.kpi = kpi.uuid
	//	where dept = '323404962476326913'
	//	order by g.dept desc, g.kpi, kpi_data.in_time
	// 根据月份进行排序
	params.Pagination = true
	//KpdDataPagingServer(params, db)
	//var tem = make(map[string]map[string][]KpiChild)
	GroupByDept(result)
	return nil, result
}
