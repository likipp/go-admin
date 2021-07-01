package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-admin/models"
	"go-admin/utils/errors"
)

func CreateBaseMenu(c *gin.Context) {
	var M models.BaseMenu
	errs := c.ShouldBindBodyWith(&M, binding.JSON).Error()
	if errs != "" {
		errors.FailWithMessage("获取前段信息失败", c)
		return
	}
	err, menu := M.CreateBaseMenu()
	if err != nil {
		errors.FailWithMessage("创建菜单失败", c)
	} else {
		errors.OkWithData(menu, c)
	}
}
