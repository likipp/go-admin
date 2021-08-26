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
		response.FailWithMessage("获取前端KPI数据失败", c)
		return
	}
	fmt.Println(KD, "未获取值前")
	KD.CreateBy = getUserUUID(c)
	fmt.Println(KD.CreateBy, KD.GroupKPI, "前端数据")
	//ok := service.HasPermissions(KD.CreateBy, KD.GroupKPI， "POST")
	//if ok {
	//	fmt.Println(ok, "是否有权限")
	//} else {
	//	c.Abort()
	//	return
	//}

	err, kpiData := KD.CreateKpiData()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	} else {
		response.OkWithData(kpiData, c)
	}
}

func GetKpiDataList(c *gin.Context) {
	var params models.KpiDataQueryParam
	gins.ParseQuery(c, &params)
	fmt.Println(params, "查询参数")
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
