package service

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	_ "github.com/go-sql-driver/mysql"
	"go-admin/config"
	"go-admin/global"
	"go-admin/utils"
	//orm "go-admin/init/database"
	"gorm.io/gorm"
	"log"
	//orm "go-admin/init/database"
	//"log"
	"strings"
	//"sync"
)

var (
	SyncedEnforcer *casbin.SyncedEnforcer
	//once           sync.Once
	err error
)

func Casbin(db *gorm.DB) {
	//once.Do(func() {
	//
	//})
	//a, err := gormadapter.NewAdapterByDB(db)
	//if err != nil {
	//	log.Fatalf("error: adapter: %s", err)
	//}
	a, _ := gormadapter.NewAdapter("mysql", "xiaom:Server@1234.com@tcp(nas.xiaom.work:3306)/") // Your driver and data source.
	//fmt.Println("Model路径:", config.AdminConfig.Casbin.ModelPath, a)
	fmt.Println(a, "a的值")
	//syncedEnforcer.AddFunction("ParamsMatch", ParamsMatchFunc)
	//syncedEnforcer.AddFunction("AdminMatch", AdminMatchFunc)
	SyncedEnforcer, err = casbin.NewSyncedEnforcer(config.AdminConfig.Casbin.ModelPath, a)
	if err != nil {
		log.Fatalf("error: syncedEnforcer: %s", err)
	}
	err = SyncedEnforcer.LoadPolicy()
	if err != nil {
		log.Fatalf("error: loadPolicy: %s", err)
	}
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

//func AddRolesForUser(user, role string) bool {
//	e, _ := Casbin()
//	ok, _ := e.HasRoleForUser(user, role)
//	if ok {
//		return ok
//	}
//	ok, _ = e.AddRoleForUser(user, role)
//	return ok
//}
//
//func HasPermissions(user, permission, method string) bool {
//	e, _ := Casbin()
//	ok := e.HasPermissionForUser(user, permission, method)
//	return ok
//}

func HasPermissionList(user, module string) []string {
	allObjects, _ := global.GSyncedEnforcer.GetImplicitPermissionsForUser(user)
	//a, _ := ccasbin.SyncedEnforcer.GetImplicitRolesForUser(user)
	fmt.Println(allObjects)
	//fmt.Println(a)
	var ss []string
	for _, v := range allObjects {
		sList := strings.Split(v[1], ":")
		if sList[0] == module {
			ss = append(ss, sList[1])
		}
	}
	rss := utils.RemoveDuplicateByMap(ss)
	return rss
}
