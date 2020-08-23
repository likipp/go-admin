package models

import (
	"errors"
	"go-admin/controller/server"
	orm "go-admin/init/database"
	globalID "go-admin/init/globalID"
	"go-admin/models/page"
	"go-admin/utils"
)

type SysUser struct {
	//UUID      uuid.UUID `json:"uuid"`
	BaseModel
	UUID     string    `json:"uuid"`
	Username string    `json:"username"`
	Password string    `json:"-"`
	NickName string    `json:"nickname" gorm:"default:'匿名用户'"`
	Avatar   string    `json:"avatar" gorm:"default:'/favicon.ico'"`
	Roles    []SysRole `json:"roles" gorm:"many2many:user_role;"`
	DeptID   string    `json:"deptID"`
	PostID   int       `json:"postID"`
	//SysDept   SysDept   `json:"dept"`
	//SysDeptId string    `json:"deptID"`
	Sex      int `json:"sex"`
	LeaderId string
	Remark   string `json:"remark"`
	Status   *int   `json:"status" gorm:"type:int(1);default:1"`
}

type UserInfo struct {
	SysUser
	DeptName string
}

func (SysUser) TableName() string {
	return "sys_user"
}

func (u *SysUser) CreateUser() (err error, userInter *SysUser) {
	var user SysUser
	//mysql.DB.Model(&u).Association("Roles").Find(&u.Roles)
	hasUser := orm.DB.Where("username = ?", u.Username).First(&user).RecordNotFound()
	if !hasUser {
		return errors.New("用户名已经注册"), nil
	} else {
		u.Roles = u.GetRoleList()
		u.UUID, err = globalID.GetID()
		if err != nil {
			return
		}
		u.Password = utils.MD5V([]byte(u.Password))
		err = orm.DB.Create(u).Error
	}
	//orm.DB.Model(&u).Related(&u.SysDept)
	//orm.DB.Model(&u).Association("Roles").Find(&u.Roles)
	return err, u
}

// 前端传递JSON格式的[]SysRole表时, 遍历获取到具体的SysRole {"roles": [{"id": 1}, {"id": 2}]}
func (u *SysUser) GetRoleList() []SysRole {
	var roles []SysRole
	for index, _ := range u.Roles {
		var role SysRole
		orm.DB.Where(&u.Roles[index]).First(&role)
		roles = append(roles, role)
	}
	return roles
}

func (u *SysUser) GetUserByUUID() (user SysUser, err error) {
	if err := orm.DB.Where("uuid = ?", u.UUID).Select("id, uuid, username, nick_name, dept_id, status").Find(&user).Error; err != nil {
		return user, errors.New("找不到该用户")
	}
	//orm.DB.Model(&user).Related(&user.SysDept)
	// orm.DB.Model(&user).Association("Roles").Count() 获取用户关联的角色组数量
	orm.DB.Model(&user).Association("Roles").Find(&user.Roles)
	return user, nil
}

func (u *SysUser) GetList(info page.InfoPage) (err error, list interface{}, total int) {
	err, db, total := server.PagingServer(u, info)
	if err != nil {
		return
	} else {
		var userList []SysUser
		// 获取用户关联的部门与角色
		var userInfoList []UserInfo
		var userInfo UserInfo
		err = db.Preload("Roles").Find(&userList).Error
		for _, value := range userList {
			var dept SysDept
			var _ = orm.DB.Where("dept_id = ?", value.DeptID).First(&dept)
			//fmt.Println(dept.DeptName, value.DeptID)
			userInfo = UserInfo{
				value,
				dept.DeptName,
			}
			//db.Model(&value).Association("DeptID").Find(&dept)
			userInfoList = append(userInfoList, userInfo)
		}
		//fmt.Println(userInfoList)
		return err, userInfoList, total
	}
}

func (u *SysUser) UpdateUser(user SysUser) (err error) {
	if err = orm.DB.Model(&u).Updates(&user).Error; err != nil {
		return errors.New("修改用户失败")
	}
	return nil
}

func (u *SysUser) DeleteUser() (err error) {
	var user SysUser
	if err = orm.DB.Where("uuid = ?", u.UUID).First(&user).Error; err == nil {
		if err = orm.DB.Unscoped().Delete(&user).Error; err != nil {
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
	err = orm.DB.Model(&u).Update("status", status).Error
	return err
}
