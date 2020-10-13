package gorm

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-admin/utils/errors"
)

func ParseJSON(c *gin.Context, obj interface{}) {
	if err := c.ShouldBindJSON(obj); err != nil {
		errors.FailWithMessage("解析请求参数发生错误", c)
	}
	errors.OkWithMessage("获取数据成功", c)
}

func ParseQuery(c *gin.Context, obj interface{}) {
	if err := c.ShouldBindQuery(obj); err != nil {
		errors.FailWithMessage("解析请求参数发生错误", c)
	}
	errors.OkWithMessage("获取数据成功", c)
}

func ParseForm(c *gin.Context, obj interface{}) {
	if err := c.ShouldBindWith(obj, binding.Form); err != nil {
		errors.FailWithMessage("解析请求参数发生错误", c)
	}
	errors.OkWithMessage("获取数据成功", c)
}
