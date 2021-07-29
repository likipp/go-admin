package service

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	_ "github.com/go-sql-driver/mysql"
	"go-admin/config"
	orm "go-admin/init/database"
	"strings"
	"sync"
)

var (
	syncedEnforcer *casbin.SyncedEnforcer
	once           sync.Once
)

func Casbin() *casbin.SyncedEnforcer {
	once.Do(func() {
		a, _ := gormadapter.NewAdapterByDB(orm.DB)
		syncedEnforcer, _ = casbin.NewSyncedEnforcer(config.AdminConfig.Casbin.ModelPath, a)
		//syncedEnforcer.AddFunction("ParamsMatch", ParamsMatchFunc)
		//syncedEnforcer.AddFunction("AdminMatch", AdminMatchFunc)
	})
	_ = syncedEnforcer.LoadPolicy()
	return syncedEnforcer
}

func ParamsMatch(fullNameKey1 string, key2 string) bool {

	key1 := strings.Split(fullNameKey1, "?")[0]
	// 剥离路径后再使用casbin的keyMatch2
	return util.KeyMatch2(key1, key2)
}

func ParamsMatchFunc(args ...interface{}) (interface{}, error) {
	name1 := args[0].(string)
	name2 := args[1].(string)
	return ParamsMatch(name1, name2), nil
}

// AdminMatch 设置系统默认管理员组, 在v1(path), v2(method)是*时, 即拥有系统的所有权
//func AdminMatch(r, p string) bool {
//	if p == "*" {
//		return true
//	}
//	return r == p
//}
//
//func AdminMatchFunc(args ...interface{}) (interface{}, error) {
//	path1 := args[0].(string)
//	path2 := args[1].(string)
//	return AdminMatch(path1, path2), nil
//}

func AddRolesForUser(user, role string) bool {
	e := Casbin()
	ok, _ := e.HasRoleForUser(user, role)
	if ok {
		return ok
	}
	ok, _ = e.AddRoleForUser(user, role)
	return ok
}
