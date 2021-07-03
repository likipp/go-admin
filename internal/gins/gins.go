package gins

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-admin/utils/response"
)

const (
	UserIDKey = "/user-id"
)

func ParseJSON(c *gin.Context, obj interface{}) {
	if err := c.ShouldBindJSON(obj); err != nil {
		response.FailWithMessage("解析请求参数发生错误", c)
		return
	}
	//response.OkWithMessage("获取数据成功", c)
}

func ParseQuery(c *gin.Context, obj interface{}) {
	if err := c.ShouldBindQuery(obj); err != nil {
		response.FailWithMessage("解析请求参数发生错误", c)
	}
	//response.OkWithData(obj, c)
}

func ParseForm(c *gin.Context, obj interface{}) {
	if err := c.ShouldBindWith(obj, binding.Form); err != nil {
		response.FailWithMessage("解析请求参数发生错误", c)
	}
	response.OkWithMessage("获取数据成功", c)
}

func GetUserID(c *gin.Context) string {
	return c.GetString(UserIDKey)
}

func SetUserID(c *gin.Context, userID string) {
	c.Set(UserIDKey, userID)
}
