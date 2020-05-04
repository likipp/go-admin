package apis

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go-admin/config"
	"go-admin/middleware"
	"go-admin/models"
	"net/http"
	"time"
)

func Authenticator(c *gin.Context) {
	var login models.Login
	_ = c.BindJSON(&login)
	user, _, err := login.GetUser()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":  401,
			"error": err.Error(),
		})
	} else {
		tokenNext(c, user)
	}
}

func tokenNext(c *gin.Context, user models.SysUser) {
	j := &middleware.JWT{
		SigningKey: []byte(config.AdminConfig.JWT.SigningKey)}
	claims := middleware.CustomClaims{
		UUID:     user.UUID,
		NickName: user.NickName,
		StandardClaims: jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000),       // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 60*60*24*7), // 过期时间 一周
			Issuer:    "xiaom",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "获取token失败",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "登录成功", "token": token, "user": user})
	}
}
