package models

import (
	"errors"
	"go-admin/controller/server"
	orm "go-admin/init/database"
	initID "go-admin/init/globalID"
	"go-admin/models/page"
)

type SysDept struct {
	BaseModel
	DeptID   string    `gorm:"column:dept_id" json:"deptID"`
	ParentId string    `gorm:"column:parent_id" json:"parent_id"`
	DeptName string    `gorm:"column:dept_name" json:"deptName"`
	DeptPath string    `gorm:"column:dept_path" json:"deptPath"`
	Sort     int       `gorm:"column:sort" json:"sort"`
	Leader   int       `gorm:"column:leader" json:"leader"`
	Status   string    `gorm:"column:status" json:"status"`
	Children []SysDept `json:"children"`
	Users    []SysUser `gorm:"foreignkey:DeptID;association_foreignkey:DeptID"`
}

func (SysDept) TableName() string {
	return "sys_dept"
}

type DeptLabel struct {
	ID       uint64      `json:"id" gorm:"-"`
	Label    string      `json:"label" gorm:"-"`
	Children []DeptLabel `json:"children" gorm:"-"`
}

// 创建部门
func (d *SysDept) Create() (*SysDept, error) {
	var dept SysDept
	result := orm.DB.Where("dept_name = ?", d.DeptName).First(&dept).Error
	if result == nil {
		err := errors.New("部门已存在")
		return d, err
	} else {
		d.DeptID, _ = initID.GetID()
		//err := orm.DB.Create(&d).Error
		//return d, nil
	}
	if d.DeptID != "" {
		var ParDept SysDept
		orm.DB.Where("dept_id = ?", d.ParentId).First(&ParDept)
		d.DeptPath = ParDept.DeptPath + d.DeptPath
	} else {
		d.DeptPath = "/" + d.DeptName
	}
	err := orm.DB.Create(&d).Error
	return d, err
	//result := mysql.DB.Table("sys_dept").Create(&d)
	//if result.Error != nil {
	//	err := result.Error
	//	return dept, err
	//}
	//dept = *d
	//return dept, nil
}

// 获取带分页的用户列表, GetInfoList
func (d *SysDept) GetList(info page.InfoPage) (err error, list interface{}, total int) {
	err, db, total := server.PagingServer(d, info)
	if err != nil {
		return
	} else {
		var depList []SysDept
		//err = db.Preload("Role").Find(&userList).Error
		err = db.Find(&depList).Error
		return err, depList, total
	}
}

func (d *SysDept) GetByUUID() (D SysDept, err error) {
	err = orm.DB.Where("dept_id = ?", d.DeptID).First(&D).Error
	if err != nil {
		return D, errors.New("未找到该部门")
	}
	return D, nil
}

// 获取部门的组织架构
func (d *SysDept) SetDept() ([]SysDept, error) {
	//list, err := e.GetPage(bl)
	var list []SysDept
	err := orm.DB.Find(&list).Error
	m := make([]SysDept, 0)
	for i := 0; i < len(list); i++ {
		if list[i].ParentId != "" {
			continue
		}
		info := DeptOrder(&list, list[i])
		m = append(m, info)
	}
	return m, err
}

// 对部门的组织树进行排列
func DeptOrder(deptList *[]SysDept, menu SysDept) SysDept {
	list := *deptList
	min := make([]SysDept, 0)
	for i := 0; i < len(list); i++ {
		if menu.DeptID != list[i].ParentId {
			continue
		}
		mi := SysDept{}
		mi.ID = list[i].ID
		mi.ParentId = list[i].ParentId
		mi.DeptName = list[i].DeptName
		mi.Sort = list[i].Sort
		mi.Leader = list[i].Leader
		mi.Status = list[i].Status
		mi.DeptID = list[i].DeptID
		//mi.CreatedAt = list[i].CreatedAt
		//mi.UpdatedAt = list[i].UpdatedAt
		mi.Children = []SysDept{}
		//ms := DeptOrder(deptList, list[i])
		min = append(min, mi)
	}
	menu.Children = min
	return menu
}
