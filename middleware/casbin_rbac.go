package middleware

import (
	"github.com/gin-gonic/gin"
	"go-admin/models"
	"go-admin/service"
	"go-admin/utils/errors"
)

func CasbinHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		claims, _ := context.Get("claims")
		waitUse, ok := claims.(*models.CustomClaims)
		if !ok {
			errors.FailWithMessage("获取用户失败", context)
			context.Abort()
			return
		}
		obj := context.Request.URL.RequestURI()
		act := context.Request.Method
		e := service.Casbin()
		// 先查看用户是否拥有权限, 如果已经拥有了权限, 则不查看所属是否拥有权限
		sub := waitUse.Username
		ok, _ = e.Enforce(sub, obj, act)
		if ok {
			context.Next()
			return
		}
		// 再循环用户所拥有的角色是否拥有权限
		for _, v := range waitUse.Roles {
			sub := v.RoleName
			// 如果用户属于管理员组，则默认拥有所有权
			if sub == "管理员组" {
				sub = "root"
			}
			ok, _ = e.Enforce(sub, obj, act)
		}
		if ok {
			context.Next()
		} else {
			errors.FailWithMessage("权限不足", context)
			context.Abort()
			return
		}
	}
}
