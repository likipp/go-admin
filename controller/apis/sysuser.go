package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-admin/models"
	"go-admin/models/page"
	"go-admin/utils/response"
	"net/http"
	"strconv"
)

//type RegisterStruct struct {
//	Username string `json:"username"`
//	Password string `json:"password"`
//	NickName string `json:"nickname" gorm:"default:'QMPlusUser'"`
//	HeaderImg string `json:"headerImg" gorm:"default:'http://www.henrongyi.top/avatar/lufu.jpg'"`
//	AuthorityId string `json:"authorityId" gorm:"default:888"`
//}

// @title Register
// @Summary 用户注册账号
// @Produce  application/json
// @Param u models.SysUser
// @Success 200 {string} string "{"code":200,"data":{},"msg":"创建成功"}"
// @Router /api/v1/base/user [post]
func CreateUser(c *gin.Context) {
	//var R RegisterStruct
	var U models.SysUser
	_ = c.ShouldBindBodyWith(&U, binding.JSON).Error()
	fmt.Println(&U, "前端传递的用户信息")
	err, user := U.CreateUser()
	if err != nil {
		response.FailWithMessage("用户创建失败", c)
		return
	} else {
		response.OkWithData(user, c)
	}
}

func GetUserByUUID(c *gin.Context) {
	var U models.SysUser
	uid := c.Param("uuid")
	//U.UUID, _ = uuid.FromString(uid)
	U.UUID = uid
	user, err := U.GetUserByUUID()
	if err != nil {
		response.FailWithMessage("用户查询失败", c)
		return
	} else {
		fmt.Println(user, "user")
		response.OkWithData(user, c)
	}
}

func GetUserList(c *gin.Context) {
	var pageInfo page.InfoPage

	// 使用Query方法
	//pageInfo.PageSize, _ = strconv.Atoi(c.Query("pageSize"))
	//pageInfo.Page, _ = strconv.Atoi(c.Query("page"))

	// 需要使用POST方法
	//_ = c.BindJSON(&pageInfo)

	// 结构体中需要定义form Tag
	_ = c.BindQuery(&pageInfo)
	//_ = c.ShouldBindJSON(&pageInfo)
	err, list, total := new(models.SysUser).GetList(pageInfo)
	//list["deptName"] = "信息部"
	//test := list.(map[string]interface{})["deptName"]
	//data := list.([]models.SysUser)
	//for _, value := range data {
	//	fmt.Println(value, "value")
	//}
	if err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"code": 400, "data": err.Error()})
		response.FailWithMessage(fmt.Sprintf("获取用户数据失败, %v", err), c)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":     200,
			"data":     list,
			"total":    total,
			"page":     pageInfo.Page,
			"pageSize": pageInfo.PageSize,
		})
		//response.OkWithData(response.PageResult{
		//	Data:     list,
		//	Total:    total,
		//	Page:     pageInfo.Page,
		//	PageSize: pageInfo.PageSize,
		//}, c)
	}
}

func UpdateUser(c *gin.Context) {
	// U代表原有的用户信息
	var U models.SysUser
	// N代表前端传递过来的用户修改信息
	var N models.SysUser
	_ = c.ShouldBindJSON(&N)
	fmt.Println(N.Roles, "roles", N.NickName)
	uid := c.Param("uuid")
	//U.UUID, _ = uuid.FromString(uid)
	U.UUID = uid
	err := U.UpdateUser(N)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "更新失败", "data": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功"})
	}
}

func DeleteUser(c *gin.Context) {
	var U models.SysUser
	uid := c.Param("uuid")
	U.UUID = uid
	err := U.DeleteUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "删除失败", "data": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "删除成功"})
	}
}

func EnableOrDisableUser(c *gin.Context) {
	var U models.SysUser
	uid := c.Param("uuid")
	status, _ := strconv.Atoi(c.Param("status"))
	U.UUID = uid
	err := U.EnableOrDisableUser(status)
	if status == 0 {
		if err != nil {
			//c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "禁用失败", "data": err.Error()})
			response.FailWithMessage("禁用失败", c)
			return
		} else {
			//c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "禁用成功"})
			response.OkWithMessage("禁用成功", c)
		}
	} else if status == 1 {
		if err != nil {
			response.FailWithMessage("启用失败", c)
			return
		} else {
			response.OkWithMessage("启用成功", c)
		}
	} else {
		response.FailWithMessage("未知状态", c)
	}
}

// 用户登录
func Login(c *gin.Context) {
	var L models.Login
	_ = c.ShouldBindJSON(&L)
	fmt.Println(&L)
	if err, _ := models.UserLogin(&L); err != nil {
		response.FailWithMessage(fmt.Sprintf("%v", err), c)
	} else {
		response.OkWithData(L, c)
	}

}

// 获取当前登录用户信息
func GetCurrentUser(c *gin.Context) {
	response.OkWithMessage("获取成功", c)
}
