package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-admin/models"
	"go-admin/utils/errors"
)

func Create(c *gin.Context) {
	var M models.BaseMenu
	errs := c.ShouldBindBodyWith(&M, binding.JSON).Error()
	if errs != "" {
		errors.FailWithMessage(errs, c)
		return
	}
	err, menu := M.CreateMenu()
	if err != nil {
		errors.FailWithMessage("创建菜单失败", c)
	} else {
		errors.OkWithData(menu, c)
	}
}
