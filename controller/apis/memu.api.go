package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-admin/models"
	"go-admin/utils/response"
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
		response.FailWithMessage(err.Error(), c)
		return
	}
	err, menu := M.CreateBaseMenu()
	if err != nil {
		response.FailWithMessage("创建菜单失败", c)
	} else {
		response.OkWithData(menu, c)
	}
}

func GetMenusTree(c *gin.Context) {
	var M models.BaseMenu
	err, m := M.GetBaseMenu()
	if err != nil {
		return
	}
	response.OkWithData(m, c)
}
