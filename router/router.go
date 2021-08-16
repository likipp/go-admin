package router

import (
	"context"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go-admin/controller/apis"
	"go-admin/middleware"
	"log"
	"os"
	"os/signal"
	"time"
	//_ "go-admin/docs"
	"net/http"
)

func InitRouter() {
	r := gin.Default()
	r.StaticFS("/images", http.Dir("./static/images"))
	r.StaticFile("/favicon.ico", "./static/images/default.jpg")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//r.POST("/login", auth.Authenticator)
	//r.Use(middleware.JWTAuth())
	r.POST("/api/v1/base/login", apis.Login)
	r.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	//baseRouter := r.Group("/api/v1/base").Use(middleware.JWTAuth())
	baseRouter := r.Group("/api/v1/base")
	//.Use(middleware.JWTAuth())
	{
		// 用户登录
		//baseRouter.POST("login", apis.Login)
		baseRouter.GET("currentUser", apis.GetCurrentUser)
		baseRouter.POST("logout", apis.Logout)
		// 用户设置router
		baseRouter.POST("users", apis.CreateUser)
		baseRouter.GET("users/:uuid", apis.GetUserByUUID)
		// 后端需要这样的格式 base/users?page=1&pageSize=3
		baseRouter.GET("users", apis.GetUserList)
		baseRouter.DELETE("users/:uuid", apis.DeleteUser)
		baseRouter.PATCH("users/:uuid", apis.UpdateUser)
		// 传递status值， 1代表禁用， 2代表启用
		baseRouter.PATCH("users/:uuid/:status", apis.EnableOrDisableUser)

		// 角色设置router
		baseRouter.GET("roles", apis.GetRoleList)
		baseRouter.POST("roles", apis.CreateRole)
		baseRouter.GET("roles/:id", apis.GetRoleByQuery)
		// 部门设置router
		baseRouter.POST("dept", apis.CreateDept)
		baseRouter.GET("dept", apis.GetAll)
		baseRouter.GET("dept/:uuid", apis.GetByUUID)
		baseRouter.GET("dept-tree", apis.GetDepTree)
		baseRouter.GET("dept-tree/:name", apis.GetDepTreeByName)
		//baseRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		// KPI设置kpi
		baseRouter.POST("kpi", apis.CreateKPI)
		baseRouter.GET("kpi", apis.GetKPIList)
		baseRouter.GET("kpi/:uuid", apis.GetKPIByUUID)
		baseRouter.PATCH("kpi/:uuid", apis.UpdateKPIByUUID)

		// GroupKPI
		baseRouter.POST("group-kpi", apis.CreateGroupKPI)
		baseRouter.GET("group-kpi", apis.GetGroupKPI)
		baseRouter.GET("group-kpi-dept", apis.GetGroupKPIDept)

		// KpiData
		baseRouter.POST("kpi-data", apis.CreateKPIData)
		baseRouter.GET("kpi-data", apis.GetKpiDataList)
		baseRouter.GET("kpi-line", apis.GetKpiDateLine)

		// Menu
		baseRouter.POST("menus", apis.CreateBaseMenu)
		baseRouter.GET("menus", apis.GetMenusTree)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
	//return r
}
