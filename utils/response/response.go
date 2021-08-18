package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (r *Response) Error() string {
	return r.ErrorMessage
}

func Result(code int, data interface{}, msg string, showType int, success bool, c *gin.Context) {
	c.JSON(code, &Response{
		ErrorCode:    code,
		Success:      success,
		ErrorMessage: msg,
		ShowType:     showType,
		Timestamp:    time.Now().Unix(),
		Data:         data,
		Host:         c.ClientIP(),
	})
}

func ResultWithPageInfo(code int, data interface{}, msg string, showType int, success bool, total int64, page, size int, c *gin.Context) {
	c.JSON(code, &PageInfo{
		Response: Response{
			ErrorCode:    code,
			Success:      success,
			ErrorMessage: msg,
			ShowType:     showType,
			Timestamp:    time.Now().Unix(),
			Data:         data,
			Host:         c.ClientIP(),
		},
		Total:    total,
		Page:     page,
		PageSize: size,
	})
}

func Ok(c *gin.Context) {
	Result(http.StatusOK, map[string]interface{}{}, "操作成功", 0, true, c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(http.StatusOK, map[string]interface{}{}, message, 0, true, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(http.StatusOK, data, "操作成功", 0, true, c)
}

func OKWithPageInfo(data interface{}, total int64, page, size int, c *gin.Context) {
	ResultWithPageInfo(http.StatusOK, data, "操作成功", 0, true, total, page, size, c)
}

func OkDetailed(data interface{}, message string, c *gin.Context) {
	Result(http.StatusOK, data, message, 0, true, c)
}

func Fail(c *gin.Context) {
	Result(http.StatusBadRequest, map[string]interface{}{}, "操作失败", 2, false, c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(http.StatusInternalServerError, map[string]interface{}{}, message, 1, false, c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(http.StatusBadRequest, data, message, 2, false, c)
}
