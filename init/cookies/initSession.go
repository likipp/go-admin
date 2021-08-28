package cookies

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"go-admin/config"
	"go-admin/utils/response"
	"gopkg.in/boj/redistore.v1"
	"net/http"
)

var RS *redistore.RediStore

func InitSession(admin config.Redis) {
	store, err := redistore.NewRediStore(10, "tcp", admin.Path, admin.Password, []byte("secret-key"))
	if err != nil {
		panic("Redis启动异常")
	}
	//store.SetMaxAge(10 * 24 * 3600)

	//store.Pool.IdleTimeout = 60*60*24*7
	//store.Options.MaxAge = 60*60*24*7
	store.SetMaxAge(60 * 60 * 24 * 7)
	store.Options.Secure = true
	store.Options.HttpOnly = true
	RS = store
}

func GetSession(c *gin.Context) (*sessions.Session, error) {
	session, err := RS.Get(c.Request, "session")
	if err != nil {
		response.Result(http.StatusBadRequest, nil, "获取session失败", 0, false, c)
		c.Abort()
		return nil, err
	}
	return session, nil
}

func SaveSession(c *gin.Context) {
	err := sessions.Save(c.Request, c.Writer)
	if err = sessions.Save(c.Request, c.Writer); err != nil {
		response.Result(http.StatusBadRequest, nil, "保存session失败", 0, false, c)
	}
}

func DeleteSession(c *gin.Context) {
	session, err := GetSession(c)
	if err != nil {
		response.Result(http.StatusBadRequest, nil, "获取session失败", 0, false, c)
		c.Abort()
		return
	}
	session.Options.MaxAge = -1
	if err := sessions.Save(c.Request, c.Writer); err != nil {
		response.Result(http.StatusBadRequest, nil, "清除session失败", 0, false, c)
		c.Abort()
		return
	}
	response.Result(http.StatusOK, nil, "退出成功", 0, true, c)
}
