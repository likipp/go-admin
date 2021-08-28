package gins

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-admin/utils/response"
	"net/http"
)

const (
	UserIDKey = "/user-id"
)

func Parse(c *gin.Context, obj interface{}) {
	if err := c.ShouldBind(&obj); err != nil {
		response.Result(http.StatusBadRequest, nil, "解析请求参数发生错误", 0, false, c)
	}
}

func ParseJSON(c *gin.Context, obj interface{}) {
	if err := c.ShouldBindJSON(obj); err != nil {
		response.Result(http.StatusBadRequest, nil, "解析请求参数发生错误", 0, false, c)
	}
}

func ParseQuery(c *gin.Context, obj interface{}) {
	if err := c.ShouldBindQuery(obj); err != nil {
		response.Result(http.StatusBadRequest, nil, "解析请求参数发生错误", 0, false, c)
	}
	//response.OkWithData(obj, c)
}

func ParseForm(c *gin.Context, obj interface{}) {
	if err := c.ShouldBindWith(obj, binding.Form); err != nil {
		response.Result(http.StatusBadRequest, nil, "解析请求参数发生错误", 0, false, c)
	}
	response.Result(http.StatusOK, nil, "获取数据成功", 0, true, c)
}

func GetUserID(c *gin.Context) string {
	return c.GetString(UserIDKey)
}

func SetUserID(c *gin.Context, userID string) {
	c.Set(UserIDKey, userID)
}
