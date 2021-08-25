package ccasbin

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go-admin/config"
	orm "go-admin/init/database"
	"log"
	"sync"
)

var (
	syncedEnforcer *casbin.SyncedEnforcer
	once           sync.Once
	err            error
)

func InitCasBin() (*casbin.SyncedEnforcer, error) {
	_, err := gormadapter.NewAdapterByDB(orm.DB)
	if err != nil {
		log.Fatalf("error: adapter: %s", err)
	}
	syncedEnforcer, err = casbin.NewSyncedEnforcer(config.AdminConfig.Casbin.ModelPath)
	if err != nil {
		log.Fatalf("error: syncedEnforcer: %s", err)
	}
	//fmt.Println(err)
	//err = syncedEnforcer.InitWithModelAndAdapter(syncedEnforcer.GetModel(), a)
	//if err != nil {
	//	log.Fatalf("error: InitWithModelAndAdapter: %s", err)
	//}
	return syncedEnforcer, err
}
