package apis

import (
	"github.com/gin-gonic/gin"
	"go-admin/internal/gins"
	"go-admin/models"
	"go-admin/utils/response"
)

func CreateGroupKPI(c *gin.Context) {
	var G models.GroupKPI
	var _ = c.ShouldBind(&G)
	err, kpi := G.CreateGroupKPI()
	if err != nil {
		response.FailWithMessage("GroupKPI创建失败", c)
		return
	} else {
		response.OkWithData(kpi, c)
	}
}

func GetGroupKPI(c *gin.Context) {
	err, results := new(models.GroupKPI).GetGroupKPI()
	if err != nil {
		response.FailWithMessage("获取列表成功", c)
		return
	} else {
		response.OkWithData(results, c)
	}
}

func GetGroupKPIDept(c *gin.Context) {
	var params models.KPIDeptQueryParam
	gins.ParseQuery(c, &params)
	err, results := new(models.GroupKPI).GetGroupKPIDept(params)
	if err != nil {
		response.FailWithMessage("获取列表失败", c)
		return
	} else {
		response.OkWithData(results, c)
	}
}
