package models

import (
	"fmt"
	"github.com/pkg/errors"
	orm "go-admin/init/database"
	"go-admin/init/globalID"
	"go-admin/internal/entity"
	"go-admin/internal/schema"
	"gorm.io/gorm"
)

type KPI struct {
	BaseModel
	UUID   string `json:"uuid"`
	Name   string `gorm:"column:name" json:"name"`
	Unit   string `gorm:"column:unit" json:"unit"`
	Status int    `gorm:"column:status" json:"status"`
}

type KPIQueryParam struct {
	schema.PaginationParam
	Name   string `form:"name"`
	Status int    `form:"status"`
}

func (KPI) TableName() string {
	return "kpi"
}

func (k *KPI) CreateKPI() (err error, KPI *KPI) {
	hasKPI := orm.DB.Where("name = ?", k.Name).Error
	hasKPIResult := errors.Is(hasKPI, gorm.ErrRecordNotFound)
	if hasKPIResult {
		return errors.New("KPI名称重复,请检查"), nil
	} else {
		k.UUID, err = initID.NewID()
		if err != nil {
			return
		}
		err = orm.DB.Create(k).Error
	}
	return err, k
}

func PagingServer(pageParams KPIQueryParam, db *gorm.DB) {
	var total int64
	limit := pageParams.PageSize
	offset := pageParams.PageSize * (pageParams.Current - 1)
	_ = db.Count(&total).Error
	db.Limit(int(limit)).Offset(int(offset)).Order("id desc")
}

func GetKpiDB(db *gorm.DB) *gorm.DB {
	return entity.GetDBWithModel(db, new(KPI))
}

func (k *KPI) GetKPIList(params KPIQueryParam) (err error, KPIList []KPI) {
	db := GetKpiDB(orm.DB)
	if v := params.Name; v != "" {
		db = db.Where("name = ?", v).Find(&KPIList)
	}
	if v := params.Status; v > 0 {
		db = db.Where("status =?", v).Find(&KPIList)
	}
	if params.Status <= 0 && params.Name == "" {
		db = db.Find(&KPIList)
	}
	params.Pagination = true
	PagingServer(params, db)
	return err, KPIList
}

func (k *KPI) GetKPIByUUID() (KPI KPI, err error) {
	db := GetKpiDB(orm.DB)
	result := db.Where("uuid = ?", k.UUID).First(&KPI)
	if result.Error != nil {
		return KPI, result.Error
	}
	return KPI, nil
}

//func GetKPIName(uuid string) (KPIName string, err error) {
//	var kpi KPI
//	db := GetKpiDB(orm.DB)
//	hasKPI := db.Where("uuid = ?", k.UUID).First(&kpi).RecordNotFound()
//	if !hasKPI {
//		return KPIName, response.New("未找到此KPI")
//	}
//	return kpi.Name, nil
//}

//func (k *KPI) BeforeUpdate(tx *gorm.DB) (err error) {
//	if tx.Statement.Changed("Name") {
//		tx.Statement.SetColumn("Name", k.Name)
//	}
//	return nil
//}

func (k *KPI) UpdateKPIByUUID() (KR KPI, err error) {
	db := GetKpiDB(orm.DB)
	fmt.Println(k, "k")
	err = db.Where("uuid", k.UUID).Model(&KPI{}).Updates(k).Error
	if err != nil {
		return KR, errors.New("KPI修改失败")
	}
	return KR, nil
}

func (k *KPI) DeleteKPIByUUID(uuid string) error {
	return nil
}
