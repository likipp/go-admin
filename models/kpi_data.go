package models

import (
	"errors"
	"github.com/jinzhu/copier"
	orm "go-admin/init/database"
	"go-admin/internal/entity"
	"go-admin/internal/schema"
	"go-admin/utils"
	"gorm.io/gorm"
	"strconv"
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
	// 获取当前服务器时间, 推算前11个月，用于数据库Between使用
	var nowMonth = time.Now().Format("2006-01")
	var beforeMonth = time.Now().AddDate(0, -11, 0).Format("2006-01")
	var selectData = "kpi_data.id, group_kpi.dept, group_kpi.kpi, kpi.name, group_kpi.l_limit, group_kpi.t_limit, group_kpi.u_limit, kpi_data.r_value, kpi.unit, kpi_data.user, kpi_data.in_time, kpi.unit"
	var joinData = "join group_kpi on kpi_data.group_kpi = group_kpi.uuid join kpi on group_kpi.kpi = kpi.uuid"
	var orderData = "group_kpi.kpi desc, group_kpi.dept, kpi_data.in_time"
	db := GetKpiDataDB(orm.DB).Select(selectData).Joins(joinData).Where("kpi_data.in_time BETWEEN ? AND ?", beforeMonth, nowMonth).Order(orderData)
	var result []Result
	if v := params.User; v != "" {
		db = db.Where("kpi_data.user = ?", v).Scan(&result)
	}
	if v := params.Dept; v != "" {
		db = db.Where("group_kpi.dept = ?", v).Scan(&result)
	}
	if v := params.GroupKPI; v != "" {
		db = db.Where("group_kpi.uuid = ?", v).Scan(&result)
	}

	if v := params.KPI; v != "" {
		db = db.Where("kpi.uuid = ?", v).Scan(&result)
	}

	if params.GroupKPI == "" && params.User == "" && params.Dept == "" {
		db = db.Scan(&result)
	}
	// 根据月份进行排序
	params.Pagination = true
	kd = GroupBy(result, time.Now())
	return nil, kd
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
		for _, v := range data {
			if s[i] == v.KPI {
				month[utils.ChangeDate(v.InTime)] = v.RValue
				month["id"] = v.ID
				month["kpi"] = v.KPI
				month["name"] = v.Name
				// 在Table中显示单位
				month["lLimit"] = strconv.Itoa(v.LLimit) + v.Unit
				month["uLimit"] = strconv.Itoa(v.ULimit) + v.Unit
				month["tValue"] = strconv.Itoa(v.TLimit) + v.Unit
			}
		}
		kList = append(kList, month)
	}
	for _, v := range kList {
		monthMap := utils.CompareByMonth(date)
		for i, a := range v {
			monthMap[i] = a
		}
		monthList = append(monthList, monthMap)
	}
	return monthList
}

func (k *KpiData) GetKPIDataForLine(params KpiDataQueryParam) (err error, r []ResultLine) {
	var nowMonth = time.Now().Format("2006-01")
	var beforeMonth = time.Now().AddDate(0, -11, 0).Format("2006-01")
	var selectData = "kpi_data.id, group_kpi.dept, group_kpi.kpi, kpi.name, group_kpi.l_limit, group_kpi.t_limit, group_kpi.u_limit, kpi_data.r_value, kpi.unit, kpi_data.user, kpi_data.in_time"
	var joinData = "join group_kpi on kpi_data.group_kpi = group_kpi.uuid join kpi on group_kpi.kpi = kpi.uuid"
	var orderData = "group_kpi.kpi, kpi_data.in_time asc"
	db := GetKpiDataDB(orm.DB).Select(selectData).Joins(joinData).Where("kpi_data.in_time BETWEEN ? AND ?", beforeMonth, nowMonth).Order(orderData)
	var result []ResultLine
	if v := params.User; v != "" {
		db = db.Where("kpi_data.user = ?", v).Scan(&result)
	}
	if v := params.Dept; v != "" {
		db = db.Where("group_kpi.dept = ?", v).Scan(&result)
	}
	if v := params.GroupKPI; v != "" {
		db = db.Where("group_kpi.uuid = ?", v).Scan(&result)
	}

	if v := params.KPI; v != "" {
		db = db.Where("kpi.uuid = ?", v).Scan(&result)
	}
	// 根据月份进行排序
	params.Pagination = true
	for i := 0; i < len(result); i++ {
		result[i].InTime = utils.ChangeDate(result[i].InTime)
	}
	result = GroupByLine(result, time.Now())
	return nil, result
}

func GroupByLine(result []ResultLine, date time.Time) []ResultLine {
	var monthStringList []ResultLine
	monthsList := utils.GetFullMonths(date)
	var a ResultLine
	for _, i := range monthsList {
		for _, v := range result {
			if v.InTime != i.Format("2006/01") {
				a.InTime = i.Format("2006/01")
				a.Name = v.Name
				a.LLimit = v.LLimit
				a.TLimit = v.TLimit
				a.Unit = v.Unit
				a.RValue = 0
				a.ULimit = v.ULimit
			}
		}
		monthStringList = append(monthStringList, a)
	}
	for i, mi := range monthStringList {
		for _, ri := range result {
			if mi.InTime == ri.InTime {
				monthStringList[i].RValue = ri.RValue
			}
		}
	}
	var b ResultLine
	var bList []ResultLine
	var temp = make(map[string]string)
	for _, v := range result {
		temp[v.Name] = v.Name
	}
	for _, v := range temp {
		for _, i := range monthStringList {
			b.InTime = i.InTime
			b.Name = v
			b.LLimit = i.LLimit
			b.TLimit = i.TLimit
			b.Unit = i.Unit
			b.ULimit = i.ULimit
			bList = append(bList, b)
		}
	}
	for i, _ := range bList {
		for _, v := range result {
			if bList[i].Name == v.Name && bList[i].InTime == v.InTime {
				_ = copier.Copy(&bList[i], &v)
			}
		}
	}
	return bList
}
