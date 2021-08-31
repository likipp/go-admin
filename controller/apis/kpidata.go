package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-admin/internal/gins"
	"go-admin/models"
	"go-admin/utils/response"
	"net/http"
)

func CreateKPIData(c *gin.Context) {
	var KD models.KpiData
	var err = c.ShouldBindJSON(&KD)
	if err != nil {
		response.Result(http.StatusBadRequest, nil, "获取前端KPI数据失败", 0, false, c)
		return
	}
	KD.CreateBy = getUserUUID(c)

	err, kpiData := KD.CreateKpiData()
	if err != nil {
		response.Result(http.StatusBadRequest, nil, "录入KPI数据失败", 1, false, c)
		return
	} else {
		response.Result(http.StatusOK, kpiData, "录入成功", 1, true, c)
	}
}

func GetKpiDataList(c *gin.Context) {
	var params models.KpiDataQueryParam
	gins.ParseQuery(c, &params)
	fmt.Println(params, "查询参数")
	err, kpiDataList := new(models.KpiData).GetKpiData(params)
	if err != nil {
		response.Result(http.StatusBadRequest, nil, "获取KPI数据失败", 1, false, c)
		return
	} else {
		response.ResultWithPageInfo(kpiDataList, "获取数据成功", 0, true, int64(len(kpiDataList)), params.Current, params.PageSize, c)
	}
}

func GetKpiDateLine(c *gin.Context) {
	var params models.KpiDataQueryParam
	gins.ParseQuery(c, &params)
	err, kpiDataList := new(models.KpiData).GetKPIDataForLine(params)
	if err != nil {
		response.Result(http.StatusBadRequest, nil, "获取KPI Line数据失败", 0, false, c)
		return
	} else {
		response.ResultWithPageInfo(kpiDataList, "获取数据成功", 0, true, int64(len(kpiDataList)), params.Current, params.PageSize, c)
	}
}
