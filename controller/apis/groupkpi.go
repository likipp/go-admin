package apis

import (
	"github.com/gin-gonic/gin"
	"go-admin/models"
	"go-admin/utils/errors"
)

func CreateGroupKPI(c *gin.Context) {
	var G models.GroupKPI
	var _ = c.ShouldBind(&G)
	err, kpi := G.CreateGroupKPI()
	if err != nil {
		errors.FailWithMessage("GroupKPI创建失败", c)
		return
	} else {
		errors.OkWithData(kpi, c)
	}
}

func GetGroupKPI(c *gin.Context) {
	err, results := new(models.GroupKPI).GetGroupKPI()
	if err != nil {
		errors.FailWithMessage("成功获取列表", c)
		return
	} else {
		errors.OkWithData(results, c)
	}
}