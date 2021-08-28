package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-admin/models"
	"go-admin/utils/response"
	"net/http"
)

func CreateBaseMenu(c *gin.Context) {
	var M models.BaseMenu
	//if err := c.ShouldBind(&M); err != nil {
	//	response.FailWithMessage("获取前段数据失败", c)
	//	return
	//}
	err := c.ShouldBindBodyWith(&M, binding.JSON)
	M.CreateBy = getUserUUID(c)
	if err != nil {
		response.Result(http.StatusBadRequest, nil, "获取数据失败", 0, false, c)
		return
	}
	err, menu := M.CreateBaseMenu()
	if err != nil {
		response.Result(http.StatusBadRequest, nil, "创建菜单失败", 0, false, c)
	} else {
		response.Result(http.StatusOK, menu, "创建菜单成功", 0, true, c)
	}
}

func GetMenusTree(c *gin.Context) {
	var M models.BaseMenu
	err, m := M.GetBaseMenu()
	if err != nil {
		return
	}
	response.Result(http.StatusOK, m, "获取列表成功", 0, true, c)
}
