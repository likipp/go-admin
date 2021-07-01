package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-admin/models"
	"go-admin/utils/errors"
)

func CreateBaseMenu(c *gin.Context) {
	var M models.BaseMenu
	//if err := c.ShouldBind(&M); err != nil {
	//	errors.FailWithMessage("获取前段数据失败", c)
	//	return
	//}
	err := c.ShouldBindBodyWith(&M, binding.JSON)
	if err != nil {
		errors.FailWithMessage(err.Error(), c)
		return
	}
	err, menu := M.CreateBaseMenu()
	if err != nil {
		errors.FailWithMessage("创建菜单失败", c)
	} else {
		errors.OkWithData(menu, c)
	}
}
