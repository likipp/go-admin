package apis

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-admin/config"
	"go-admin/global"
	"go-admin/init/cookies"
	"go-admin/middleware"
	"go-admin/models"
	"go-admin/utils"
	"go-admin/utils/response"
	"net/http"
	"strconv"
	"time"
)

func CreateUser(c *gin.Context) {
	var U models.SysUser
	//var _ = c.ShouldBind(&U)
	err := c.ShouldBindBodyWith(&U, binding.JSON)
	if err != nil {
		response.Result(http.StatusBadRequest, nil, "获取前段数据失败", 0, false, c)
		return
	}

	U.Password = "123456"
	err, user := U.CreateUser()
	if err != nil {
		response.Result(http.StatusBadRequest, nil, "用户创建失败", 1, false, c)
		return
	} else {
		response.Result(http.StatusOK, user, "用户创建成功", 1, true, c)
	}
}

func GetUserByUUID(c *gin.Context) {
	var U models.SysUser
	uid := c.Param("uuid")
	//U.UUID, _ = uuid.FromString(uid)
	U.UUID = uid
	user, err := U.GetUserByUUID()
	if err != nil {
		response.Result(http.StatusBadRequest, nil, "用户查询失败", 1, false, c)
		return
	} else {
		response.Result(http.StatusOK, user, "查询用户成功", 0, true, c)
	}
}

func GetUserList(c *gin.Context) {
	//var pageInfo page.InfoPage
	var userQuery models.UserQuery
	// 使用Query方法
	//pageInfo.PageSize, _ = strconv.Atoi(c.Query("pageSize"))
	//pageInfo.Page, _ = strconv.Atoi(c.Query("page"))

	// 需要使用POST方法
	//_ = c.BindJSON(&pageInfo)

	// 结构体中需要定义form Tag
	status := c.PostForm("status")
	if status == "" {
		userQuery.Status = 3
	} else {
		userQuery.Status = utils.StringConvInt(status)
	}

	_ = c.BindQuery(&userQuery)
	fmt.Println(userQuery, "用户条件")
	//_ = c.ShouldBindJSON(&pageInfo)
	err, list, total := new(models.SysUser).GetList(userQuery)
	if err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"code": 400, "data": err.Error()})
		response.Result(http.StatusBadRequest, nil, "获取用户数据失败", 0, false, c)
	} else {
		response.ResultWithPageInfo(list, "获取数据成功", 0, true, total, userQuery.Page, userQuery.PageSize, c)
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
		response.Result(http.StatusBadRequest, nil, "更新失败", 0, false, c)
	} else {
		response.Result(http.StatusOK, nil, "更新成功", 0, true, c)
	}
}

func DeleteUser(c *gin.Context) {
	var U models.SysUser
	uid := c.Param("uuid")
	U.UUID = uid
	err := U.DeleteUser()
	if err != nil {
		response.Result(http.StatusBadRequest, nil, "删除失败", 0, false, c)
	} else {
		response.Result(http.StatusOK, nil, "删除成功", 0, true, c)
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
			response.Result(http.StatusBadRequest, nil, "启用失败", 0, false, c)
			return
		} else {
			response.Result(http.StatusOK, nil, "启用成功", 0, true, c)
		}
	} else if status == 0 {
		if err != nil {
			response.Result(http.StatusBadRequest, nil, "禁用失败", 0, false, c)
			return
		} else {
			response.Result(http.StatusOK, nil, "禁用成功", 0, true, c)
		}
	} else {
		response.Result(http.StatusBadRequest, nil, "未知状态", 0, false, c)
	}
}

// Login 用户登录
func Login(c *gin.Context) {
	var L models.Login
	_ = c.ShouldBindJSON(&L)
	err, user := models.UserLogin(&L)
	if err != nil {
		response.Result(http.StatusBadRequest, nil, err.Error(), 0, false, c)
	}
	//token := GetToken(c, *user)
	//session, _ := cookies.GetSession(c)
	//session.Values["nickname"] = user.Username
	//session.Values["name"] = user.NickName
	//session.Values["avatar"] = "https://gw.alipayobjects.com/zos/antfincdn/XAosXuNZyF/BiazfanxmamNRoxxVxka.png"
	//session.Values["uuid"] = user.UUID
	//session.Values["access"] = "admin"
	//session.Values["token"] = token
	//global.GUser.Username = user.Username
	////fmt.Println(user.Username, "&user.Username")
	//global.GUser.NickName = user.NickName
	////global.GUser.DeptID = &string(user.DeptID)
	//cookies.SaveSession(c)
	GetTokenAndSession(c, *user)
	response.Result(http.StatusOK, user, "保存数据成功", 0, true, c)

}

// GetCurrentUser 获取当前登录用户信息
func GetCurrentUser(c *gin.Context) {
	var user models.CurrentUser
	fmt.Println("进入获取用户信息页面")
	session, err := cookies.GetSession(c)
	fmt.Println(err, "页面错误")
	user.Avatar = session.Values["avatar"].(string)
	user.UUID = session.Values["uuid"].(string)
	user.Nickname = session.Values["nickname"].(string)
	user.Access = "admin"
	user.Name = session.Values["name"].(string)
	c.JSONP(http.StatusOK, user)
}

func Logout(c *gin.Context) {
	cookies.DeleteSession(c)
}

func GetTokenAndSession(c *gin.Context, user models.SysUser) {
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
		response.Result(http.StatusBadRequest, nil, "获取Token失败", 0, false, c)
		return
	}
	session, _ := cookies.GetSession(c)
	session.Values["nickname"] = user.Username
	session.Values["name"] = user.NickName
	session.Values["avatar"] = "https://gw.alipayobjects.com/zos/antfincdn/XAosXuNZyF/BiazfanxmamNRoxxVxka.png"
	session.Values["uuid"] = user.UUID
	session.Values["access"] = "admin"
	session.Values["token"] = token
	global.GUser.Username = user.Username
	global.GUser.NickName = user.NickName
	//global.GUser.DeptID = &string(user.DeptID)
	cookies.SaveSession(c)
	//return token
}

//func Sessions(c *gin.Context, user *models.SysUser) () {
//	token := GetToken(c, *user)
//	session, _ := cookies.GetSession(c)
//	session.Values["nickname"] = user.Username
//	session.Values["name"] = user.NickName
//	session.Values["avatar"] = "https://gw.alipayobjects.com/zos/antfincdn/XAosXuNZyF/BiazfanxmamNRoxxVxka.png"
//	session.Values["uuid"] = user.UUID
//	session.Values["access"] = "admin"
//	session.Values["token"] = token
//	global.GUser.Username = user.Username
//	global.GUser.NickName = user.NickName
//	//global.GUser.DeptID = &string(user.DeptID)
//	cookies.SaveSession(c)
//}

func getUserUUID(c *gin.Context) string {
	session, _ := cookies.GetSession(c)
	return session.Values["uuid"].(string)
}
