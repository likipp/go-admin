package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-admin/models"
	"go-admin/utils/errors"
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
