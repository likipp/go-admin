package router

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go-admin/controller/apis"
	"net/http"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.StaticFS("/images", http.Dir("./static/images"))
	r.StaticFile("/favicon.ico", "./static/images/default.jpg")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//r.POST("/login", auth.Authenticator)
	//.Use(middleware.JWTAuth())

	baseRouter := r.Group("/api/v1/base")
	//.Use(middleware.JWTAuth())
	{
		baseRouter.POST("user", apis.CreateUser)
		baseRouter.GET("user/:uuid", apis.GetUserByUUID)
		// base/list?page=2&pageSize=3
		baseRouter.GET("users", apis.GetUserList)
		baseRouter.DELETE("user/:uuid", apis.DeleteUser)
		baseRouter.PATCH("user/:uuid", apis.UpdateUser)

		baseRouter.POST("role", apis.CreateRole)

		baseRouter.POST("dept", apis.CreateDept)
		baseRouter.GET("dept", apis.GetAll)
		baseRouter.GET("dept/:uuid", apis.GetByUUID)
		baseRouter.GET("dept-tree", apis.GetDepTree)
	}

	return r
}
