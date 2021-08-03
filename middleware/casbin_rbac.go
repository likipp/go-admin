package middleware

import (
	"github.com/gin-gonic/gin"
	"go-admin/models"
	"go-admin/service"
	"go-admin/utils/response"
)

func CasbinHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		claims, _ := context.Get("claims")
		waitUse, ok := claims.(*models.CustomClaims)
		if !ok {
			response.FailWithMessage("获取用户失败", context)
			context.Abort()
			return
		}
		obj := context.Request.URL.RequestURI()
		act := context.Request.Method
		e := service.Casbin()
		service.HasPermissions("359681968171909121", "324851701305573377")
		// 先查看用户是否拥有权限, 如果已经拥有了权限, 则不查看所属是否拥有权限
		sub := waitUse.UUID
		ok, _ = e.Enforce(sub, obj, act)
		//if err != nil {
		//	response.FailWithMessage("权限不足", context)
		//	context.Abort()
		//	return
		//}
		//context.Next()
		if ok {
			context.Next()
		} else {
			response.FailWithMessage("权限不足", context)
			context.Abort()
			return
		}
		//service.AddRolesForUser(sub, "default")
		//if ok {
		//	context.Next()
		//	//return
		//}
		// 再循环用户所拥有的角色是否拥有权限, 如果用户属于管理员组，则默认拥有所有权
		//for _, v := range waitUse.Roles {
		//	sub := v.RoleName
		//	ok, _ = e.Enforce(sub, obj, act)
		//}
		//if ok {
		//	context.Next()
		//} else {
		//	response.FailWithMessage("权限不足", context)
		//	context.Abort()
		//	return
		//}
		//response.FailWithMessage("权限不足", context)
		//context.Abort()
		//return
	}
}
