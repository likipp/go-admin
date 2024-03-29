package models

import (
	"errors"
	"github.com/jinzhu/copier"
	orm "go-admin/init/database"
	"go-admin/init/globalID"
	"go-admin/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SysUser struct {
	BaseModel
	UUID     string    `json:"uuid"`
	Username string    `json:"username"`
	Password string    `json:"-"`
	NickName string    `json:"nickname" gorm:"default:'匿名用户'"`
	Avatar   string    `json:"avatar" gorm:"default:'/favicon.ico'"`
	Roles    []SysRole `json:"roles" gorm:"many2many:users_roles;"`
	DeptID   uint      `json:"deptID"`
	PostID   int       `json:"postID"`
	Sex      int       `json:"sex"`
	LeaderId string
	Remark   string `json:"remark"`
	Status   int    `json:"status" gorm:"type:int(1);default:1"`
}

type UserInfo struct {
	SysUser
	DeptName string
}

type LoginResponse struct {
	User      SysUser `json:"user"`
	Token     string  `json:"token"`
	ExpiresAt int64   `json:"expiresAt"`
}

type UserQuery struct {
	Page     int    `form:"current"`
	PageSize int    `form:"pageSize"`
	Status   int    `form:"status"`
	Username string `form:"username"`
	NickName string `form:"nickname"`
	Sex      int    `form:"sex"`
	DeptID   string `form:"deptID"`
}

type UserFilterNoPage struct {
	Status   int    `json:"status"`
	Username string `json:"username"`
	NickName string `json:"nickname"`
	Sex      int    `json:"sex"`
}

type CurrentUser struct {
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
	Name     string `json:"name"`
	Title    string `json:"title"`
	UUID     string `json:"uuid"`
	Access   string `json:"access"`
}

func (SysUser) TableName() string {
	return "sys_user"
}

func Paging(query UserQuery, model interface{}) (err error, db *gorm.DB, total int64) {
	limit := query.PageSize
	offset := query.PageSize * (query.Page - 1)
	var user SysUser
	if query.Status == 3 {
		err = copier.Copy(&user, &query)
		user.Status = 0
		db = orm.DB.Where(&user).Find(model).Count(&total).Limit(limit).Offset(offset).Order("id desc")
	} else {
		err = copier.Copy(&user, &query)
		db = orm.DB.Where(&user).Find(model).Count(&total).Limit(limit).Offset(offset).Order("id desc")
		db = db.Where(&user, "Status").Find(model).Count(&total).Limit(limit).Offset(offset).Order("id desc")
	}

	return err, db, total
}

func (u *SysUser) CreateUser() (err error, userInter *SysUser) {
	var user SysUser
	//mysql.DB.Model(&u).Association("Roles").Find(&u.Roles)
	hasUser := orm.DB.Where("username = ?", u.Username).First(&user).Error
	hasUserResult := errors.Is(hasUser, gorm.ErrRecordNotFound)
	if !hasUserResult {
		return errors.New("用户名已经注册"), nil
	} else {
		u.Roles = u.GetRoleList()
		u.UUID, err = initID.NewID()
		if err != nil {
			return
		}
		//u.Password = utils.MD5V([]byte(u.Password))
		u.Password = utils.PasswordHash(u.Password)
		err = orm.DB.Create(u).Error
	}
	//orm.DB.Model(&u).Related(&u.SysDept)
	//orm.DB.Model(&u).Association("Roles").Find(&u.Roles)
	return err, u
}

// GetRoleList 前端传递JSON格式的[]SysRole表时, 遍历获取到具体的SysRole {"roles": [{"id": 1}, {"id": 2}]}
func (u *SysUser) GetRoleList() []SysRole {
	var roles []SysRole
	for index := range u.Roles {
		var role SysRole
		orm.DB.Where(&u.Roles[index]).First(&role)
		roles = append(roles, role)
	}
	return roles
}

func (u *SysUser) GetUserByUUID() (userInfo UserInfo, err error) {
	var user SysUser
	var dept SysDept
	// 查询出部分字段 orm.DB.Where("uuid = ?", u.UUID).Select("id, uuid, username, nick_name, dept_id, status, sex, created_at").Find(&user)
	if err := orm.DB.Where("uuid = ?", u.UUID).Find(&user).Error; err != nil {
		return userInfo, errors.New("找不到该用户")
	}
	var _ = orm.DB.Where("dept_id = ?", user.DeptID).First(&dept)
	//orm.DB.Model(&user).Related(&user.SysDept)
	// orm.DB.Model(&user).Association("Roles").Count() 获取用户关联的角色组数量
	orm.DB.Model(&user).Association("Roles").Find(&user.Roles)
	userInfo.SysUser = user
	userInfo.DeptName = dept.DeptName
	return userInfo, nil
}

func (u *SysUser) GetList(query UserQuery) (err error, list interface{}, total int64) {
	var userList []SysUser
	// 获取用户关联的部门与角色
	var userInfoList []UserInfo
	var userInfo UserInfo
	err, db, total := Paging(query, &userList)
	if err != nil {
		return
	} else {
		//if filters.Status != 3 {
		//	err = db.Preload("Roles").Find(&userList).Where("status = ?", filters.Status).Error
		//}
		err = db.Preload("Roles").Find(&userList).Error
		for _, value := range userList {
			var dept SysDept
			var _ = orm.DB.Where("dept_id = ?", value.DeptID).First(&dept)
			userInfo = UserInfo{
				value,
				dept.DeptName,
			}
			//db.Model(&value).Association("DeptID").Find(&dept)
			userInfoList = append(userInfoList, userInfo)
		}
		return err, userInfoList, total
	}
}

func (u *SysUser) UpdateUser(user SysUser) (err error) {
	//orm.DB.Set("gorm:association_autocreate", false).Save(&user)
	//err = orm.DB.Model(&u).Updates(&user).Error
	if err = orm.DB.Where("uuid = ?", u.UUID).Model(u).Updates(&user.Roles).Error; err != nil {
		return errors.New("修改用户失败")
	}
	return nil
}

func (u *SysUser) DeleteUser() (err error) {
	var user SysUser
	if err = orm.DB.Where("uuid = ?", u.UUID).First(&user).Error; err == nil {
		if err = orm.DB.Select(clause.Associations).Unscoped().Delete(&user).Error; err != nil {
			return errors.New("删除用户失败")
		}
		return err
	} else {
		return errors.New("未找到要删除的用户")
	}
}

func (u *SysUser) EnableOrDisableUser(status int) (err error) {
	var user SysUser
	if err = orm.DB.Where("uuid = ?", u.UUID).First(&user).Error; err != nil {
		return errors.New("未找到此用户")
	}
	// 根据前端传递的status值, 更新用户的状态信息
	// 使用Update时，数据库执行时间过长，SLOW SQL >= 200ms, 后面更改成UpdateColumn
	//err = orm.DB.Model(&user).Update("status", status).Error
	//单个Update时，需要传递id主键值，所以需要传递整个use结构体，或者传递id
	err = orm.DB.Model(&user).UpdateColumn("status", status).Error
	return err
}
