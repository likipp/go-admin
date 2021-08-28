package apis

import (
	"github.com/gin-gonic/gin"
	"go-admin/internal/gins"
	"go-admin/models"
	"go-admin/utils/response"
)

func CreateKPI(c *gin.Context) {
	var K models.KPI
	//var _ = c.ShouldBind(&K)
	gins.Parse(c, K)
	err, kpi := K.CreateKPI()
	if err != nil {
		response.FailWithMessage("KPI创建失败", c)
		return
	} else {
		response.OkWithData(kpi, c)
	}
}

func GetKPIList(c *gin.Context) {
	var kpiList []models.KPI
	var params models.KPIQueryParam
	gins.ParseQuery(c, &params)
	err, kpiList, total := new(models.KPI).GetKPIList(params)
	if err != nil {
		response.FailWithMessage("KPI查询失败", c)
	} else {
		response.OKWithPageInfo(kpiList, total, params.Current, params.PageSize, c)
	}
}

func GetKPIByUUID(c *gin.Context) {
	var K models.KPI
	K.UUID = c.Param("uuid")
	kpi, err := K.GetKPIByUUID()
	if err != nil {
		response.FailWithMessage("KPI查询失败", c)
		return
	} else {
		response.OkWithData(kpi, c)
	}
}

func UpdateKPIByUUID(c *gin.Context) {
	var newK models.KPI
	gins.ParseJSON(c, &newK)
	// 获取到UUID, 只有uuid有值时才能更新成功
	uid := c.Param("uuid")
	newK.UUID = uid
	_, err := newK.GetKPIByUUID()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	}
	// 获取当前登录用户uuid, 操作后写入数据库
	user := getUserUUID(c)
	newK.UpdateBy = user
	_, err = newK.UpdateKPIByUUID()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	} else {
		response.OkWithMessage("操作成功", c)
	}
}
