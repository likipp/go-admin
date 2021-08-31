package apis

import (
	"github.com/gin-gonic/gin"
	"go-admin/internal/gins"
	"go-admin/models"
	"go-admin/utils/response"
	"net/http"
)

func CreateKPI(c *gin.Context) {
	var K models.KPI
	//var _ = c.ShouldBind(&K)
	gins.Parse(c, K)
	err, kpi := K.CreateKPI()
	if err != nil {
		response.Result(http.StatusBadRequest, nil, "KPI创建失败", 1, false, c)
		return
	} else {
		response.Result(http.StatusOK, kpi, "创建成功", 1, true, c)
	}
}

func GetKPIList(c *gin.Context) {
	var kpiList []models.KPI
	var params models.KPIQueryParam
	gins.ParseQuery(c, &params)
	err, kpiList, total := new(models.KPI).GetKPIList(params)
	if err != nil {
		response.Result(http.StatusBadRequest, nil, "KPI查询失败", 1, false, c)
	} else {
		response.ResultWithPageInfo(kpiList, "获取列表成功", 1, true, total, params.Current, params.PageSize, c)
	}
}

func GetKPIByUUID(c *gin.Context) {
	var K models.KPI
	K.UUID = c.Param("uuid")
	kpi, err := K.GetKPIByUUID()
	if err != nil {
		response.Result(http.StatusBadRequest, nil, "KPI查询失败", 1, false, c)
		return
	} else {
		response.Result(http.StatusOK, kpi, "查询成功", 1, true, c)
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
		response.Result(http.StatusBadRequest, nil, "更新失败", 1, false, c)
	}
	// 获取当前登录用户uuid, 操作后写入数据库
	user := getUserUUID(c)
	newK.UpdateBy = user
	_, err = newK.UpdateKPIByUUID()
	if err != nil {
		response.Result(http.StatusBadRequest, nil, "更新失败", 1, false, c)
		return
	} else {
		response.Result(http.StatusBadRequest, nil, "更新成功", 1, true, c)
	}
}
