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
	RValue   int    `gorm:"column:r_value"    json:"r_value"`
	InTime   string `gorm:"column:in_time"    json:"in_time"`
	User     string `gorm:"column:user"       json:"user"`
	Dept     string `gorm:"column:dept"       json:"dept"`
	GroupKPI string `gorm:"column:group_kpi"  json:"group_kpi"`
}

type KpiDataResult struct {
	Dept    string
	KpiData map[string]interface{}
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
	hasKpiData := db.Where("group_kpi = ? AND in_time = ? AND dept = ?", k.GroupKPI, k.InTime, k.Dept).First(&result).RecordNotFound()
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

func (k *KpiData) GetKpiData(params KpiDataQueryParam) (err error, kd []*KpiData) {
	db := GetKpiDataDB(orm.DB)
	if v := params.Dept; v != "" {
		db = db.Select("r_value, in_time, user, dept, group_kpi").Order("in_time desc").Where("dept = ?", v).Find(&kd)
	}
	if v := params.GroupKPI; v != "" {
		db = db.Select("r_value, in_time, user, dept, group_kpi").Order("in_time desc").Where("group_kpi = ?", v).Group("dept").Find(&kd)
	}
	if v := params.User; v != "" {
		db = db.Select("r_value, in_time, user, dept, group_kpi").Order("in_time desc").Where("user = ?", v).Group("dept").Find(&kd)
	}
	if params.Dept == "" && params.GroupKPI == "" && params.User == "" {
		db = db.Limit(100).Select("dept, any_value(user) as user, any_value(group_kpi) as group_kpi, any_value(r_value) as r_value, any_value(in_time) as in_time").Group("dept").Find(&kd)
	}
	// 根据月份进行排序
	//params.Pagination = true
	//db.Order("in_time desc").Find(&kd)
	//KpdDataPagingServer(params, db)
	return nil, kd
}
