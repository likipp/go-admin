package cookies

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"go-admin/config"
	"go-admin/utils/response"
	"gopkg.in/boj/redistore.v1"
)

var RS *redistore.RediStore

func InitSession(admin config.Redis) {
	store, err := redistore.NewRediStore(10, "tcp", admin.Path, "", []byte("secret-key"))
	if err != nil {
		panic("Redis启动异常")
	}
	store.SetMaxAge(10 * 24 * 3600)
	store.Options.Secure = true
	store.Options.HttpOnly = true
	RS = store
}

func GetSession(c *gin.Context) *sessions.Session {
	session, err := RS.Get(c.Request, "session")
	if err != nil {
		response.FailWithMessage("获取session失败", c)
		c.Abort()
		return nil
	}
	return session
}

func SaveSession(c *gin.Context) {
	err := sessions.Save(c.Request, c.Writer)
	if err = sessions.Save(c.Request, c.Writer); err != nil {
		response.FailWithMessage("保存session失败!", c)
	}
}

func DeleteSession(c *gin.Context) {
	session := GetSession(c)
	session.Options.MaxAge = -1
	if err := sessions.Save(c.Request, c.Writer); err != nil {
		response.FailWithMessage("清除session失败.", c)
		c.Abort()
		return
	}
	response.OkWithMessage("退出成功", c)
}
