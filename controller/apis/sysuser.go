package apis

import (
	"github.com/gin-gonic/gin"
	"go-admin/models"
	"go-admin/models/page"
	"go-admin/utils"
	"net/http"
)

//type RegisterStruct struct {
//	Username string `json:"username"`
//	Password string `json:"password"`
//	NickName string `json:"nickname" gorm:"default:'QMPlusUser'"`
//	HeaderImg string `json:"headerImg" gorm:"default:'http://www.henrongyi.top/avatar/lufu.jpg'"`
//	AuthorityId string `json:"authorityId" gorm:"default:888"`
//}

// @Tags Base
// @Summary 用户注册账号
// @Produce  application/json
// @Param data body model.SysUser true "用户注册接口"
// @Success 200 {string} string "{"code":200,"data":{},"msg":"创建成功"}"
// @Router /base/user [post]
func CreateUser(c *gin.Context) {
	//var R RegisterStruct
	var U models.SysUser
	_ = c.ShouldBindJSON(&U)
	err, user := U.CreateUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "创建失败", "data": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建成功", "data": user})
	}
}

func GetUserByUUID(c *gin.Context) {
	var U models.SysUser
	uid := c.Param("uuid")
	//U.UUID, _ = uuid.FromString(uid)
	U.UUID = utils.StringConUint(uid)
	user, err := U.GetUserByUUID()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "data": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": user,
		})
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
	err, list, total := new(models.SysUser).GetList(pageInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "data": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":     200,
			"userList": list,
			"total":    total,
			"page":     pageInfo.Page,
			"pageSize": pageInfo.PageSize,
		})
	}
}

func UpdateUser(c *gin.Context) {
	// U代表原有的用户信息
	var U models.SysUser
	// N代表前端传递过来的用户修改信息
	var N models.SysUser
	_ = c.ShouldBindJSON(&N)
	uid := c.Param("uuid")
	//U.UUID, _ = uuid.FromString(uid)
	U.UUID = utils.StringConUint(uid)
	err := U.UpdateUser(N)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "更新失败", "data": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功"})
	}
}

func DeleteUser(c *gin.Context) {
	var U models.SysUser
	//id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	uid := c.Param("uuid")
	//U.UUID, _ = uuid.FromString(uid)
	U.UUID = utils.StringConUint(uid)
	//_ = c.BindJSON(&U)
	err := U.DeleteUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "删除失败", "data": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "删除成功"})
	}
}
