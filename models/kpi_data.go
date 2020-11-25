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
	//Dept    string
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
	InTime string `json:"in_time"`
	//GroupKPI string `gorm:"column:group_kpi"  json:"group_kpi"`
}

type ResultWithMonth struct {
	KPI    string `json:"kpi"`
	ULimit int    `json:"u_limit"`
	LLimit int    `json:"l_limit"`
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

func (k *KpiData) GetKpiData(params KpiDataQueryParam) (err error, kd []map[string]interface{}) {
	var selectData = "group_kpi.dept, group_kpi.kpi, group_kpi.l_limit, group_kpi.t_limit, group_kpi.u_limit, kpi_data.r_value, kpi.unit, kpi_data.user, kpi_data.in_time"
	var joinData = "join group_kpi on kpi_data.group_kpi = group_kpi.uuid join kpi on group_kpi.kpi = kpi.uuid"
	var orderData = "group_kpi.kpi desc, group_kpi.dept, kpi_data.in_time"
	db := GetKpiDataDB(orm.DB).Select(selectData).Joins(joinData).Order(orderData).Limit(12)
	//kDB := GetKpiDataDB(orm.DB).Select("group_kpi.kpi").Joins(joinData).Order(orderData).Limit(12)
	var result []Result
	if v := params.User; v != "" {
		db = db.Where("kpi_data.user = ?", v).Scan(&result)
	}
	if v := params.Dept; v != "" {
		db = db.Where("group_kpi.dept = ?", v).Scan(&result)
	}
	if v := params.GroupKPI; v != "" {
		db = db.Where("group_kpi.kpi = ?", v).Scan(&result)
	}

	if params.GroupKPI == "" && params.User == "" && params.Dept == "" {
		db = db.Scan(&result)
	}

	// 根据月份进行排序
	params.Pagination = true
	//KpdDataPagingServer(params, db
	//res := GetKPICate(result)
	kd = GroupBy(result)
	return nil, kd
}

func GroupBy(data []Result) []map[string]interface{} {
	var kList = make([]map[string]interface{}, 0)

	//var s = []string{"324858678177955841", "324858629188485121", "324858517754216449"}
	var s []string
	var temp = map[string]bool{}

	for i := 0; i < len(data); i++ {
		if _, ok := temp[data[i].KPI]; ok {

		} else {
			temp[data[i].KPI] = true
			s = append(s, data[i].KPI)
		}
	}
	fmt.Println(s, "list")
	for i := 0; i < len(s); i++ {
		var month = make(map[string]interface{})
		for _, v := range data {
			if s[i] == v.KPI {
				month[v.InTime] = v.RValue
				month["KPI"] = v.KPI
				month["LLimit"] = v.LLimit
				month["ULimit"] = v.ULimit
				month["TLimit"] = v.TLimit
			}
		}
		kList = append(kList, month)
	}
	return kList
}
