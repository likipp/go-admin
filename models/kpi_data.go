package models

import (
	"errors"
	"fmt"
	orm "go-admin/init/database"
	"go-admin/internal/entity"
	"go-admin/internal/schema"
	"go-admin/utils"
	"gorm.io/gorm"
	"strings"
	"time"
)

type KpiData struct {
	BaseModel
	UUID     string `gorm:"column:uuid"       json:"uuid"`
	RValue   int    `gorm:"column:r_value"    json:"r_value"`
	InTime   string `gorm:"column:in_time"    json:"in_time"`
	User     string `gorm:"column:user"       json:"user"`
	GroupKPI string `gorm:"column:group_kpi"  json:"group_kpi"`
}

type Result struct {
	ID     string `json:"id"`
	Dept   string `json:"dept"`
	KPI    string `json:"kpi"`
	Name   string `json:"name"`
	ULimit int    `json:"u_limit"`
	LLimit int    `json:"l_limit"`
	TLimit int    `json:"t_limit"`
	RValue int    `json:"r_value"`
	Unit   string `json:"unit"`
	User   string `json:"user"`
	InTime string `json:"in_time"`
}

type ResultLine struct {
	Name   string `json:"type"`
	ULimit int    `json:"u_limit"`
	LLimit int    `json:"l_limit"`
	TLimit int    `json:"t_limit"`
	RValue int    `json:"value"`
	Unit   string `json:"unit"`
	InTime string `json:"date"`
}

type KpiDataQueryParam struct {
	schema.PaginationParam
	GroupKPI string `form:"group_kpi"`
	User     string `form:"user"`
	Dept     string `form:"dept"`
	KPI      string `form:"kpi"`
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
	hasKpiData := db.Where("group_kpi = ? AND in_time = ?", k.GroupKPI, k.InTime).First(&result).Error
	hasKpiDataResult := errors.Is(hasKpiData, gorm.ErrRecordNotFound)
	if !hasKpiDataResult {
		return errors.New("KPI数据已经录入"), kd
	}
	err = db.Create(&k).Error
	if err != nil {
		return errors.New("创建KPI数据失败"), kd
	}
	return nil, k

}

func KpdDataPagingServer(pageParams KpiDataQueryParam, db *gorm.DB) {
	var total int64
	limit := pageParams.PageSize
	offset := pageParams.PageSize * (pageParams.Current - 1)
	_ = db.Count(&total).Error
	db.Limit(int(limit)).Offset(int(offset)).Order("in_time desc")
}

func (k *KpiData) GetKpiData(params KpiDataQueryParam) (err error, kd []map[string]interface{}) {
	var month = time.Now().Format("2006-01")
	var selectData = "kpi_data.id, group_kpi.dept, group_kpi.kpi, kpi.name, group_kpi.l_limit, group_kpi.t_limit, group_kpi.u_limit, kpi_data.r_value, kpi.unit, kpi_data.user, kpi_data.in_time"
	var joinData = "join group_kpi on kpi_data.group_kpi = group_kpi.uuid join kpi on group_kpi.kpi = kpi.uuid"
	var orderData = "group_kpi.kpi desc, group_kpi.dept, kpi_data.in_time"
	db := GetKpiDataDB(orm.DB).Select(selectData).Joins(joinData).Order(orderData).Limit(12)
	//kDB := GetKpiDataDB(orm.DB).Select("group_kpi.kpi").Joins(joinData).Order(orderData).Limit(12)
	var result []Result
	if v := params.User; v != "" {
		db = db.Where("kpi_data.user = ? and in_time < ?", v, month).Scan(&result)
	}
	if v := params.Dept; v != "" {
		db = db.Where("group_kpi.dept = ? and in_time < ?", v, month).Scan(&result)
	}
	if v := params.GroupKPI; v != "" {
		db = db.Where("group_kpi.uuid = ? and in_time < ?", v, month).Scan(&result)
	}

	if v := params.KPI; v != "" {
		db = db.Where("kpi.uuid = ? and in_time < ?", v, month).Scan(&result)
	}

	if params.GroupKPI == "" && params.User == "" && params.Dept == "" {
		db = db.Scan(&result)
	}
	// 根据月份进行排序
	params.Pagination = true
	kd = GroupBy(result, time.Now())
	return nil, kd
}

func getKeys(m map[string]interface{}) []string {
	// 数组默认长度为map长度,后面append时,不需要重新申请内存和拷贝,效率很高
	keys := make([]string, 0, len(m))
	for k := range m {
		if strings.Contains(k, "/") {
			keys = append(keys, k)
		}
	}
	return keys
}

func GroupBy(data []Result, date time.Time) []map[string]interface{} {
	var kList = make([]map[string]interface{}, 0)
	var monthList = make([]map[string]interface{}, 0)
	var s []string
	var temp = map[string]bool{}

	for i := 0; i < len(data); i++ {
		if _, ok := temp[data[i].KPI]; ok {

		} else {
			temp[data[i].KPI] = true
			s = append(s, data[i].KPI)
		}
	}
	if s == nil {
		return monthList
	}
	for i := 0; i < len(s); i++ {
		var month = make(map[string]interface{})
		var monthTemp = make(map[string]interface{})
		for _, v := range data {
			if s[i] == v.KPI {
				month[utils.ChangeDate(v.InTime)] = v.RValue
				month["id"] = v.ID
				month["kpi"] = v.KPI
				month["name"] = v.Name
				month["lLimit"] = v.LLimit
				month["uLimit"] = v.ULimit
				month["tValue"] = v.TLimit
				//month["className"] = SetClassName(v.LLimit, v.TLimit, v.RValue)
			}
		}

		kList = append(kList, month)
		if len(month) < 18 {
			monthTemp = utils.CompareByMonth(date)
			for _, v := range kList {
				for i, v := range v {
					monthTemp[i] = v
				}
			}
		}
		monthList = append(monthList, monthTemp)
	}
	return monthList
}

func (k *KpiData) GetKPIDataForLine(params KpiDataQueryParam) (err error, r []ResultLine) {
	var month = time.Now().Format("2006-01")
	var selectData = "kpi_data.id, group_kpi.dept, group_kpi.kpi, kpi.name, group_kpi.l_limit, group_kpi.t_limit, group_kpi.u_limit, kpi_data.r_value, kpi.unit, kpi_data.user, kpi_data.in_time"
	var joinData = "join group_kpi on kpi_data.group_kpi = group_kpi.uuid join kpi on group_kpi.kpi = kpi.uuid"
	var orderData = "kpi_data.in_time asc, group_kpi.kpi"
	db := GetKpiDataDB(orm.DB).Select(selectData).Joins(joinData).Order(orderData).Limit(12)
	var result []ResultLine
	//if v := params.User; v != "" {
	//	db = db.Where("kpi_data.user = ?", v).Scan(&result)
	//}
	//if v := params.Dept; v != "" {
	//	db = db.Where("group_kpi.dept = ?", v).Scan(&result)
	//}
	//if v := params.GroupKPI; v != "" {
	//	db = db.Where("group_kpi.uuid = ?", v).Scan(&result)
	//}
	if v := params.User; v != "" {
		db = db.Where("kpi_data.user = ? and in_time < ?", v, month).Scan(&result)
	}
	if v := params.Dept; v != "" {
		db = db.Where("group_kpi.dept = ? and in_time < ?", v, month).Scan(&result)
	}
	if v := params.GroupKPI; v != "" {
		db = db.Where("group_kpi.uuid = ? and in_time < ?", v, month).Scan(&result)
	}

	//if params.GroupKPI == "" && params.User == "" && params.Dept == "" {
	//	db = db.Scan(&result)
	//}
	if v := params.KPI; v != "" {
		db = db.Where("kpi.uuid = ? and in_time < ?", v, month).Scan(&result)
	}
	fmt.Println(result, "result")
	// 根据月份进行排序
	params.Pagination = true
	for i := 0; i < len(result); i++ {
		result[i].InTime = utils.ChangeDate(result[i].InTime)
	}
	GroupByLine(result, time.Now())
	return nil, result
}

func GroupByLine(result []ResultLine, date time.Time) {
	var monthTimeList []time.Time
	//monthMap := make(map[string]interface{}, 12)
	var monthStringList []ResultLine
	fmt.Println(result, "result")
	for i := 1; i <= 12; i++ {
		m := date.AddDate(0, -i, 0)
		monthTimeList = append(monthTimeList, m)
	}
	for n := 0; n <= len(monthTimeList); n++ {
		for i := 1; i < len(monthTimeList)-n; i++ {
			if monthTimeList[i].Before(monthTimeList[i-1]) {
				monthTimeList[i], monthTimeList[i-1] = monthTimeList[i-1], monthTimeList[i]
			}
		}
	}

	for _, i := range monthTimeList {
		var a ResultLine
		//monthMap["date"] = i.Format("2006/01")
		//monthStringList = append(monthStringList, monthMap)
		for _, v := range result {
			if v.InTime != i.Format("2006/01") {
				a.InTime = i.Format("2006/01")
				a.Name = v.Name
				a.LLimit = v.LLimit
				a.TLimit = 0
				a.Unit = v.Unit
				a.RValue = v.RValue
				a.ULimit = v.ULimit
			}
		}
		monthStringList = append(monthStringList, a)
	}
}
