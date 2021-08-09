package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-admin/models"
	"go-admin/service"
	"go-admin/utils/response"
)

func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get("claims")
		waitUse, ok := claims.(*models.CustomClaims)
		if !ok {
			response.FailWithMessage("获取用户失败", c)
			c.Abort()
			return
		}
		obj := c.Request.URL.RequestURI()
		act := c.Request.Method
		e, err := service.Casbin()
		if err != nil {
			c.Abort()
			response.FailWithMessage("Casbin失败", c)
		}
		// 先循环用户所拥有的角色是否拥有权限, 如果用户属于管理员组，则默认拥有所有权
		for _, v := range waitUse.Roles {
			sub := v.RoleName
			ok, _ = e.Enforce(sub, obj, act)
			// 如果拥有administrators权限，直接通过
			if ok {
				c.Next()
				return
			}
		}
		//service.HasPermissions("359681968171909121", "324851701305573377")
		fmt.Println("开始校验用户权限")
		// 先查看用户是否拥有权限, 如果已经拥有了权限, 则不查看所属是否拥有权限
		sub := waitUse.UUID
		ok, err = e.Enforce(sub, obj, act)
		if ok {
			c.Next()
		}
		//service.AddRolesForUser(sub, "default")
		response.FailWithMessage("权限不足", c)
		c.Abort()
		return
	}
}
