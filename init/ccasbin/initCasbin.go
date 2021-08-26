package ccasbin

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go-admin/config"
	orm "go-admin/init/database"
	"log"
)

var SyncedEnforcer *casbin.SyncedEnforcer

func InitCasBin() (*casbin.SyncedEnforcer, error) {
	//a, err := gormadapter.NewAdapterByDBWithCustomTable(orm.DB, &CasbinRule{})
	//a, err := gormadapter.NewAdapter("mysql", "xiaom:Server@1234.com@tcp(nas.xiaom.work:3306)/qmPlus", true)
	a, err := gormadapter.NewAdapterByDB(orm.DB)
	if err != nil {
		log.Fatalf("error: adapter: %s", err)
	}
	SyncedEnforcer, err = casbin.NewSyncedEnforcer(config.AdminConfig.Casbin.ModelPath, a)
	if err != nil {
		log.Fatalf("error: syncedEnforcer: %s", err)
		return nil, err
	}
	err = SyncedEnforcer.LoadPolicy()
	if err != nil {
		log.Fatalf("error: adapter: %s", err)
		return nil, err
	}
	return SyncedEnforcer, err
}
