package apis

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-admin/config"
	orm "go-admin/init/database"
	"go-admin/middleware"
	"go-admin/models"
	"go-admin/utils"
	"go-admin/utils/response"
	"net/http"
	"strconv"
	"time"
)

func CreateUser(c *gin.Context) {
	//var R RegisterStruct
	var U models.SysUser
	var _ = c.ShouldBind(&U)
	err := c.ShouldBindBodyWith(&U, binding.JSON)
	if err != nil {
		response.FailWithMessage("获取前段数据失败", c)
		return
	}

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
	} else {
		userFilter.Status = utils.StringConvInt(status)
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

	if status == 1 {
		if err != nil {
			response.FailWithMessage("启用失败", c)
			return
		} else {
			response.OkWithMessage("启用成功", c)
		}
	} else if status == 0 {
		if err != nil {
			response.FailWithMessage("禁用失败", c)
			return
		} else {
			response.OkWithMessage("禁用成功", c)
		}
	} else {
		response.FailWithMessage("未知状态", c)
	}
}

// Login 用户登录
func Login(c *gin.Context) {
	var L models.Login
	_ = c.ShouldBindJSON(&L)
	if err, user := models.UserLogin(&L); err != nil {
		response.FailWithMessage(fmt.Sprintf("%v", err), c)
	} else {
		session, err := orm.Store.Get(c.Request, "session")
		if err != nil {
			fmt.Println("初始化session成功")
		}
		session.Values["username"] = user.Username
		session.Values["nickname"] = user.NickName
		session.Values["avatar"] = user.Avatar
		session.Values["uuid"] = user.UUID
		orm.Store.Save(c.Request, c.Writer, session)
		http.Error(c.Writer, "", http.StatusOK)
		GetToken(c, *user)
	}

}

// GetCurrentUser 获取当前登录用户信息
func GetCurrentUser(c *gin.Context) {
	var user models.CurrentUser
	if claims, exists := c.Get("claims"); !exists {
		response.FailWithMessage("获取Token失败", c)
	} else {
		waitUse := claims.(*models.CustomClaims)
		user.Avatar = "https://gw.alipayobjects.com/zos/antfincdn/XAosXuNZyF/BiazfanxmamNRoxxVxka.png"
		user.UUID = waitUse.UUID
		user.Nickname = waitUse.Username
		user.Access = "admin"
		user.Name = waitUse.NickName
	}
	//fmt.Println(user, "打印用户")
	//response.OkWithData(user, c)
	//c.JSON(200, gin.H{"data": user})
	c.JSONP(http.StatusOK, user)
}

func GetToken(c *gin.Context, user models.SysUser) {
	j := &middleware.JWT{
		SigningKey: []byte(config.AdminConfig.JWT.SigningKey),
	}
	clams := models.CustomClaims{
		UUID:       user.UUID,
		ID:         user.ID,
		NickName:   user.NickName,
		Username:   user.Username,
		Roles:      user.Roles,
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

func getUserUUID(c *gin.Context) string {
	if claims, exists := c.Get("claims"); !exists {
		response.FailWithMessage("获取Token失败", c)
		return ""
	} else {
		waitUse := claims.(*models.CustomClaims)
		return waitUse.UUID
	}
}
