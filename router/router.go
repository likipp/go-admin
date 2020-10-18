package router

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go-admin/controller/apis"
	//_ "go-admin/docs"
	"net/http"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.StaticFS("/images", http.Dir("./static/images"))
	r.StaticFile("/favicon.ico", "./static/images/default.jpg")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//r.POST("/login", auth.Authenticator)
	//r.Use(middleware.JWTAuth())
	r.POST("/api/v1/base/login", apis.Login)

	//baseRouter := r.Group("/api/v1/base").Use(middleware.JWTAuth())
	baseRouter := r.Group("/api/v1/base")
	//.Use(middleware.JWTAuth())
	{
		// 用户登录
		//baseRouter.POST("login", apis.Login)
		baseRouter.GET("currentUser", apis.GetCurrentUser)
		// 用户设置router
		baseRouter.POST("user", apis.CreateUser)
		baseRouter.GET("users/:uuid", apis.GetUserByUUID)
		// 后端需要这样的格式 base/users?page=1&pageSize=3
		baseRouter.GET("users", apis.GetUserList)
		baseRouter.DELETE("users/:uuid", apis.DeleteUser)
		baseRouter.PATCH("users/:uuid", apis.UpdateUser)
		// 传递status值， 1代表禁用， 2代表启用
		baseRouter.PATCH("users/:uuid/:status", apis.EnableOrDisableUser)

		// 角色设置router
		baseRouter.GET("roles", apis.GetRoleList)
		baseRouter.POST("role", apis.CreateRole)

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
	}

	return r
}
