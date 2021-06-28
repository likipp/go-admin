package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-admin/models"
	"go-admin/service"
	"go-admin/utils/errors"
)

func CasbinHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		claims, _ := context.Get("claims")
		waitUse, ok := claims.(*models.CustomClaims)
		//if !ok {
		//	return
		//}
		obj := context.Request.URL.RequestURI()
		act := context.Request.Method
		fmt.Println(waitUse.Roles[0])
		sub := waitUse.Roles[0].RoleName
		fmt.Println(sub)
		e := service.Casbin()
		ok, _ = e.Enforce(sub, obj, act)
		if ok {
			context.Next()
		} else {
			errors.FailWithMessage("权限不足", context)
			context.Abort()
			return
		}
	}
}
