package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

//const (
//	ERROR   = 7
//	SUCCESS = 0
//)

type Response struct {
	ErrorCode    int         `json:"errorCode"`
	Success      bool        `json:"success"`
	ErrorMessage string      `json:"errorMessage"`
	Timestamp    int64       `json:"timestamp"`
	ShowType     int         `json:"showType"`
	Data         interface{} `json:"data"`
	Host         string      `json:"host"`
}

func (r *Response) Error() string {
	return r.ErrorMessage
}

func Result(code int, data interface{}, msg string, success bool, c *gin.Context) {
	c.JSON(code, &Response{
		ErrorCode:    code,
		Success:      success,
		ErrorMessage: msg,
		Timestamp:    time.Now().Unix(),
		Data:         data,
		Host:         c.ClientIP(),
	})
}

func Ok(c *gin.Context) {
	Result(http.StatusOK, map[string]interface{}{}, "操作成功", true, c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(http.StatusOK, map[string]interface{}{}, message, true, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(http.StatusOK, data, "操作成功", true, c)
}

func OkDetailed(data interface{}, message string, c *gin.Context) {
	Result(http.StatusOK, data, message, true, c)
}

func Fail(c *gin.Context) {
	Result(http.StatusBadRequest, map[string]interface{}{}, "操作失败", false, c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(http.StatusBadRequest, map[string]interface{}{}, message, false, c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(http.StatusBadRequest, data, message, false, c)
}
