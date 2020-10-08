package apis

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-admin/config"
	"go-admin/models"
	"go-admin/models/request"
	"go-admin/utils/jwtauth"
	"go-admin/utils/response"
	"net/http"
	"strconv"
	"time"
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
	var _ = c.ShouldBind(&U)
	fmt.Println(&U, "ShouldBind")
	_ = c.ShouldBindBodyWith(&U, binding.JSON).Error()
	U.Password = "123456"
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
		response.OkWithData(user, c)
	}
}

func GetUserList(c *gin.Context) {
	//var pageInfo page.InfoPage
	var userFilter models.UserFilter
	// 使用Query方法
	//pageInfo.PageSize, _ = strconv.Atoi(c.Query("pageSize"))
	//pageInfo.Page, _ = strconv.Atoi(c.Query("page"))

	// 需要使用POST方法
	//_ = c.BindJSON(&pageInfo)

	// 结构体中需要定义form Tag
	status := c.PostForm("status")
	if status == "" {
		userFilter.Status = 3
	}
	_ = c.BindQuery(&userFilter)
	//_ = c.ShouldBindJSON(&pageInfo)
	err, list, total := new(models.SysUser).GetList(userFilter)
	if err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"code": 400, "data": err.Error()})
		response.FailWithMessage(fmt.Sprintf("获取用户数据失败, %v", err), c)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":     200,
			"data":     list,
			"total":    total,
			"page":     userFilter.Page,
			"pageSize": userFilter.PageSize,
		})
		//response.OkWithData(response.PageResult{
		//	Data:     list,
		//	Total:    total,
		//	Page:     userFilter.Page,
		//	PageSize: userFilter.PageSize,
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
	if err, user := models.UserLogin(&L); err != nil {
		response.FailWithMessage(fmt.Sprintf("%v", err), c)
	} else {
		GetToken(c, *user)
	}

}

// 获取当前登录用户信息
func GetCurrentUser(c *gin.Context) {
	response.OkWithMessage("获取成功", c)
}

func GetToken(c *gin.Context, user models.SysUser) {
	j := &jwtauth.JWT{
		SigningKey: []byte(config.AdminConfig.JWT.SigningKey),
	}
	clams := request.CustomClaims{
		UUID:       user.UUID,
		ID:         user.ID,
		NickName:   user.NickName,
		Username:   user.Username,
		BufferTime: 60 * 60 * 24, // 缓冲时间1天 缓冲时间内会获得新的token刷新令牌 此时一个用户会存在两个有效令牌 但是前端只留一个 另一个会丢失
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,       // 签名生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*7, // 过期时间 7天
			Issuer:    "xiao",                         // 签名的发行者
		},
	}
	token, err := j.CreateToken(clams)
	if err != nil {
		response.FailWithMessage("获取Token失败", c)
		return
	}
	response.OkWithData(models.LoginResponse{
		User:      user,
		Token:     token,
		ExpiresAt: clams.StandardClaims.ExpiresAt * 1000,
	}, c)
}
