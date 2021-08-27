package ccasbin

import (
	"github.com/casbin/casbin/v2"
	gormAdapter "github.com/casbin/gorm-adapter/v3"
	"go-admin/config"
	"go-admin/global"
	orm "go-admin/init/database"
	"log"
)

var syncedEnforcer *casbin.SyncedEnforcer

func InitCasBin() {
	//a, err := gormAdapter.NewAdapter("mysql", "xiaom:Server@1234.com@tcp(nas.xiaom.work:3306)/qmPlus", true)
	a, err := gormAdapter.NewAdapterByDB(orm.DB)
	if err != nil {
		log.Printf("创建Casbin适配器失败:%v", err)
	}
	syncedEnforcer, err = casbin.NewSyncedEnforcer(config.AdminConfig.Casbin.ModelPath, a)
	if err != nil {
		log.Printf("创建Casbin调度器失败:%v", err)
	}
	global.GSyncedEnforcer = syncedEnforcer
	err = global.GSyncedEnforcer.LoadPolicy()
	if err != nil {
		log.Printf("加载Casbin策略失败:%v", err)
	}
}
