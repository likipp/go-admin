package models

import (
	"errors"
	"github.com/jinzhu/copier"
	"go-admin/controller/server"
	orm "go-admin/init/database"
	initID "go-admin/init/globalID"
	"go-admin/models/page"
)

type SysDept struct {
	BaseModel
	DeptID string `gorm:"column:dept_id" json:"deptID"`
	// 自引用时, 必须是指针, 否则创建数据会失败
	//ManageID *uint
	ParentId uint   `gorm:"column:parent_id" json:"parent_id"`
	DeptName string `gorm:"column:dept_name" json:"deptName"`
	//DeptPath string    `gorm:"column:dept_path" json:"deptPath"`
	Sort   int    `gorm:"column:sort" json:"sort"`
	Leader int    `gorm:"column:leader" json:"leader"`
	Status string `gorm:"column:status" json:"status"`
	//Children []SysDept `gorm:"foreignKey:ParentId" json:"children"`
	//Children []SysDept `json:"children"`
	Users []SysUser `gorm:"foreignKey:DeptID" json:"users"`
}

type SysDeptInfo struct {
	ID       uint          `json:"key"`
	ParentId uint          `json:"parent_id"`
	DeptName string        `json:"title"`
	DeptPath string        `json:"deptPath"`
	Sort     int           `json:"sort"`
	Leader   int           `json:"leader"`
	Status   string        `json:"status"`
	Children []SysDeptInfo `json:"children"`
	//Leaf     bool    `json:"leaf`
	EnableUsersCount  int `json:"en_users_count"`
	DisableUsersCount int `json:"dis_users_count"`
}

func (SysDept) TableName() string {
	return "sys_dept"
}

type DeptLabel struct {
	ID       uint64      `json:"id" gorm:"-"`
	Label    string      `json:"label" gorm:"-"`
	Children []DeptLabel `json:"children" gorm:"-"`
}

// Create 创建部门
func (d *SysDept) Create() (*SysDept, error) {
	var dept SysDept
	result := orm.DB.Where("dept_name = ?", d.DeptName).First(&dept).Error
	if result == nil {
		err := errors.New("部门已存在")
		return d, err
	} else {
		d.DeptID, _ = initID.NewID()
	}
	//if d.ParentId != "" {
	//	var ParDept SysDept
	//	orm.DB.Where("dept_id = ?", d.ParentId).First(&ParDept)
	//	d.DeptPath = ParDept.DeptPath + d.DeptName
	//} else {
	//	d.DeptPath = "/" + d.DeptName
	//}
	err := orm.DB.Create(&d).Error
	return d, err
}

// GetList 获取带分页的部门列表, GetInfoList
func (d *SysDept) GetList(info page.InfoPage) (err error, list interface{}, total int64) {
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

// DeptTree 获取部门的组织架构
func (d *SysDept) DeptTree() ([]SysDeptInfo, error) {
	var list []SysDept
	err := orm.DB.Order("sort").Find(&list).Error
	m := make([]SysDeptInfo, 0)
	for i := 0; i < len(list); i++ {
		if list[i].ParentId != 0 {
			continue
		}
		info := DeptOrder(&list, list[i])
		m = append(m, info)
	}
	return m, err
}

func (d *SysDept) DeptTreeByName() ([]SysDeptInfo, error) {
	var list []SysDept
	//var childList []SysDept
	err := orm.DB.Where("dept_name LIKE ?", "%"+d.DeptName+"%").Order("sort").Find(&list).Error
	m := make([]SysDeptInfo, 0)
	for i := 0; i < len(list); i++ {
		info := DeptOrder(&list, list[i])
		m = append(m, info)
	}
	m = DeptCompare(m)
	return m, err
}

// DeptCompare 过滤重复的部门
func DeptCompare(deptList []SysDeptInfo) []SysDeptInfo {
	result := make([]SysDeptInfo, 0)
	for i := 0; i < len(deptList); i++ {
		repeat := false
		for j := 0; j < len(deptList); j++ {
			if deptList[i].ParentId == deptList[j].ID {
				repeat = true
			}
		}
		if !repeat {
			result = append(result, deptList[i])
		}
	}
	return result
}

// DeptOrder 对部门的组织树进行排列
func DeptOrder(deptList *[]SysDept, dept SysDept) SysDeptInfo {
	list := *deptList
	min := make([]SysDeptInfo, 0)
	deptInfo := SysDeptInfo{}
	// copier可以拷贝相同类型的结构
	err := copier.Copy(&deptInfo, &dept)
	//deptInfo.ID = dept.ID
	if err != nil {
		return SysDeptInfo{}
	}
	deptInfo.EnableUsersCount = int(orm.DB.Model(dept).Where("status = ?", 1).Association("Users").Count())
	deptInfo.DisableUsersCount = int(orm.DB.Model(dept).Where("status = ?", 2).Association("Users").Count())
	for i := 0; i < len(list); i++ {
		if dept.ID != list[i].ParentId {
			continue
		}
		mi := SysDeptInfo{}
		mi.ParentId = list[i].ParentId
		mi.DeptName = list[i].DeptName
		mi.Sort = list[i].Sort
		mi.Leader = list[i].Leader
		mi.Status = list[i].Status
		mi.ID = list[i].ID
		mi.DisableUsersCount = int(orm.DB.Model(list[i]).Where("status = ?", 1).Association("Users").Count())
		mi.EnableUsersCount = int(orm.DB.Model(list[i]).Where("status = ?", 2).Association("Users").Count())
		mi.Children = []SysDeptInfo{}
		ms := DeptOrder(deptList, list[i])
		min = append(min, ms)
	}
	deptInfo.Children = min
	return deptInfo
}
