package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-admin/internal/gins"
	"go-admin/models"
	"go-admin/utils/errors"
	"net/http"
)

func CreateKPI(c *gin.Context) {
	var K models.KPI
	var _ = c.ShouldBind(&K)
	fmt.Println(K, "K")
	err, kpi := K.CreateKPI()
	if err != nil {
		errors.FailWithMessage("KPI创建失败", c)
		return
	} else {
		errors.OkWithData(kpi, c)
	}
}

func GetKPIList(c *gin.Context) {
	var kpiList []models.KPI
	var params models.KPIQueryParam
	gins.ParseQuery(c, &params)
	err, kpiList := new(models.KPI).GetKPIList(params)
	if err != nil {
		errors.FailWithMessage("KPI查询失败", c)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":     200,
			"data":     kpiList,
			"total":    len(kpiList),
			"page":     params.Current,
			"pageSize": params.PageSize,
		})
	}
}
