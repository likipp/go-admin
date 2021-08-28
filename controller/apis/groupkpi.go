package apis

import (
	"github.com/gin-gonic/gin"
	"go-admin/internal/gins"
	"go-admin/models"
	"go-admin/utils/response"
	"net/http"
)

func CreateGroupKPI(c *gin.Context) {
	var G models.GroupKPI
	var _ = c.ShouldBind(&G)
	err, kpi := G.CreateGroupKPI()
	if err != nil {
		response.Result(http.StatusBadRequest, nil, "GroupKPI创建失败", 0, false, c)

		return
	} else {
		response.Result(http.StatusOK, kpi, "创建成功", 0, true, c)
	}
}

func GetGroupKPI(c *gin.Context) {
	err, results := new(models.GroupKPI).GetGroupKPI()
	if err != nil {
		response.Result(http.StatusBadRequest, nil, "获取列表失败", 0, false, c)
		return
	} else {
		response.Result(http.StatusBadRequest, results, "获取列表成功", 0, true, c)
	}
}

func GetGroupKPIDept(c *gin.Context) {
	var params models.KPIDeptQueryParam
	gins.ParseQuery(c, &params)
	err, results := new(models.GroupKPI).GetGroupKPIDept(params)
	if err != nil {
		response.Result(http.StatusBadRequest, nil, "获取列表失败", 0, false, c)
		return
	} else {
		response.Result(http.StatusBadRequest, results, "获取列表失败", 0, false, c)
	}
}
