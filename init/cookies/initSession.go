package cookies

import (
	"go-admin/config"
	"gopkg.in/boj/redistore.v1"
)

var RS *redistore.RediStore

func InitSession(admin config.Redis) {
	store, err := redistore.NewRediStore(10, "tcp", admin.Path, "", []byte("secret-key"))
	if err != nil {
		panic("Redis启动异常")
	} else {
		store.SetMaxAge(10 * 24 * 3600)
		store.Options.Secure = true
		store.Options.HttpOnly = true
		// Delete session
		// session.Options.MaxAge = -1
		// 由于 Cookie 已经过期，已拒绝 Cookie “session”。
		// store.Options.MaxAge = -1
		RS = store
		// 添加Close()登录直接报错，暂时未找到原因
		//defer RS.Close()
	}

}
