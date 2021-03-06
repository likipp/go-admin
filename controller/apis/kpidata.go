package apis

import (
	"github.com/gin-gonic/gin"
	"go-admin/internal/gins"
	"go-admin/models"
	"go-admin/utils/response"
	"net/http"
)

func CreateKPIData(c *gin.Context) {
	var KD models.KpiData
	var _ = c.ShouldBind(&KD)
	err, kpiData := KD.CreateKpiData()
	if err != nil {
		response.FailWithMessage("创建KPI数据失败", c)
		return
	} else {
		response.OkWithData(kpiData, c)
	}
}

func GetKpiDataList(c *gin.Context) {
	var params models.KpiDataQueryParam
	gins.ParseQuery(c, &params)
	err, kpiDataList := new(models.KpiData).GetKpiData(params)
	if err != nil {
		response.FailWithMessage("获取KPI数据失败", c)
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":     200,
			"data":     kpiDataList,
			"total":    len(kpiDataList),
			"page":     params.Current,
			"pageSize": params.PageSize,
		})
	}
}

func GetKpiDateLine(c *gin.Context) {
	var params models.KpiDataQueryParam
	gins.ParseQuery(c, &params)
	err, kpiDataList := new(models.KpiData).GetKPIDataForLine(params)
	if err != nil {
		response.FailWithMessage("获取KPI Line数据失败", c)
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":     200,
			"data":     kpiDataList,
			"total":    len(kpiDataList),
			"page":     params.Current,
			"pageSize": params.PageSize,
		})
	}
}
