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
		waitUse := claims.(*models.CustomClaims)
		obj := context.Request.URL.RequestURI()
		act := context.Request.Method
		sub := waitUse.RoleId
		e := service.Casbin()
		ok, err := e.Enforce(sub, obj, act)
		fmt.Println(err, "错误信息")
		if ok {
			context.Next()
		} else {
			errors.FailWithMessage("权限不足", context)
			context.Abort()
			return
		}
	}
}
