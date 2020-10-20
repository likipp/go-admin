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

type KpiDataQueryParam struct {
	schema.PaginationParam
	GroupKPI string `form:"group_kpi"`
	User     string `form:"user"`
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

func GroupByDept(kd []*KpiData) {
	//var result []*KpiData
	//var t map[string]interface{}
	for i, v := range kd {
		fmt.Println(i, "数据", v)
		//t[v.Dept] = make([]interface)

	}
}

type Result struct {
	RValue int `json:"r_value"`
	//InTime   string `gorm:"column:in_time"    json:"in_time"`
	User string `json:"user"`
	//GroupKPI string `gorm:"column:group_kpi"  json:"group_kpi"`
	Dept   string `json:"deptID"`
	KPI    string `json:"KpiID"`
	ULimit int    `json:"u_limit"`
	LLimit int    `json:"l_limit"`
	TLimit int    `json:"t_limit"`
	//Unit   string `json:"unit"`
}

func (k *KpiData) GetKpiData(params KpiDataQueryParam) (err error, kd []Result) {
	db := GetKpiDataDB(orm.DB)
	var result []Result
	//if v := params.Dept; v != "" {
	//	db = db.Select("r_value, in_time, user, dept, group_kpi").Order("dept desc, in_time").Where("dept = ?", v).Find(&kd)
	//}
	if v := params.GroupKPI; v != "" {
		db = db.Select("r_value, in_time, user, group_kpi").Order("in_time").Where("group_kpi = ?", v).Find(&kd)
	}
	if v := params.User; v != "" {
		db = db.Select("r_value, in_time, user, group_kpi").Order("in_time").Where("user = ?", v).Find(&kd)
	}
	if params.GroupKPI == "" && params.User == "" {
		//db = db.Limit(100).Select("dept, any_value(user) as user, any_value(group_kpi) as group_kpi, any_value(r_value) as r_value, any_value(in_time) as in_time").Group("dept").Find(&kd)
		//db = db.Order("in_time").Find(&kd)
		db = db.Select("group_kpi.dept, group_kpi.kpi, group_kpi.l_limit, group_kpi.t_limit, group_kpi.u_limit, kpi_data.r_value, kpi_data.user").Joins("join group_kpi on kpi_data.group_kpi = group_kpi.uuid").Scan(&result)

	}
	// 根据月份进行排序
	params.Pagination = true
	//KpdDataPagingServer(params, db)
	//var tem = make(map[string]map[string][]KpiChild)
	for _, v := range result {
		fmt.Println(v)
	}
	return nil, result
}
