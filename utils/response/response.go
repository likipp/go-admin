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

func ResultWithPageInfo(data interface{}, msg string, showType int, success bool, total int64, page, size int, c *gin.Context) {
	c.JSON(http.StatusOK, &PageInfo{
		Response: Response{
			ErrorCode:    http.StatusOK,
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

//func Ok() {
//	Result(http.StatusOK, map[string]interface{}{}, "操作成功", 0, true)
//}
//
//func OkWithMessage(message string) {
//	Result(http.StatusOK, map[string]interface{}{}, message, 0, true)
//}
//
//func OkWithData(data interface{}) {
//	Result(http.StatusOK, data, "操作成功", 0, true)
//}
//
//func OKWithPageInfo(data interface{}, total int64, page, size int) {
//	ResultWithPageInfo(http.StatusOK, data, "操作成功", 0, true, total, page, size)
//}
//
//func OkDetailed(data interface{}, message string) {
//	Result(http.StatusOK, data, message, 0, true)
//}
//
//func Fail() {
//	Result(http.StatusBadRequest, map[string]interface{}{}, "操作失败", 2, false)
//}
//
//func FailWithMessage(message string) {
//	Result(http.StatusInternalServerError, map[string]interface{}{}, message, 1, false)
//}
//
//func FailWithDetailed(data interface{}, message string) {
//	Result(http.StatusBadRequest, data, message, 2, false)
//}
