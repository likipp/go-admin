package models

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	orm "go-admin/init/database"
	"go-admin/internal/entity"
	"go-admin/internal/schema"
	"go-admin/utils"
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

func Test(kd []Result, listSort map[string]int, dept map[string]map[string]int, kpi string) {
	for _, v := range kd {
		if v.RValue != 0 {
			if utils.StringConvInt(v.InTime[5:7]) >= 10 {
				listSort[utils.StringConvJoin(v.InTime[2:4], v.InTime[5:7])] = v.RValue
			} else {
				listSort[utils.StringConvJoin(v.InTime[2:4], v.InTime[6:7])] = v.RValue
			}
		}
	}
}

func GetKPICate(kd []Result) map[string][]map[string]map[string]int {
	var temp = map[string]map[string]int{}
	var result = map[string][]map[string]map[string]int{}
	var ss = make(map[string]int)
	var dept string
	for i := 0; i < len(kd); i++ {

		if utils.StringConvInt(kd[i].InTime[5:7]) >= 10 {
			ss[utils.StringConvJoin(kd[i].InTime[2:4], kd[i].InTime[5:7])] = kd[i].RValue
		} else {
			ss[utils.StringConvJoin(kd[i].InTime[2:4], kd[i].InTime[6:7])] = kd[i].RValue
		}

		temp[kd[i].KPI] = ss
		dept = kd[i].Dept
	}
	result[dept] = append(result[dept], temp)
	return result
}

func (k *KpiData) GetKpiData(params KpiDataQueryParam) (err error, kd map[string][]map[string]map[string]int) {
	var selectData = "group_kpi.dept, group_kpi.kpi, group_kpi.l_limit, group_kpi.t_limit, group_kpi.u_limit, kpi_data.r_value, kpi.unit, kpi_data.user, kpi_data.in_time"
	var joinData = "join group_kpi on kpi_data.group_kpi = group_kpi.uuid join kpi on group_kpi.kpi = kpi.uuid"
	var orderData = "group_kpi.dept desc, group_kpi.kpi, kpi_data.in_time"
	db := GetKpiDataDB(orm.DB).Select(selectData).Joins(joinData).Order(orderData).Limit(12)
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
	GroupBy(result)
	return nil, nil
}

func GroupBy(data []Result) {
	var i = 0
	var j int
	var kList = make([]map[string]interface{}, 2)
	var month = make(map[string]interface{})
	for j = i; j < len(data) && data[i].KPI == data[j].KPI; j++ {
		month["KPI"] = data[j].KPI
		month["LLimit"] = data[j].LLimit
		month["ULimit"] = data[j].ULimit
		month["TLimit"] = data[j].TLimit
		month[data[j].InTime] = data[j].RValue
	}

	kList = append(kList, month)
	fmt.Println(kList)
}
