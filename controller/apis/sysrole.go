package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-admin/models"
	"go-admin/models/page"
	"go-admin/utils/response"
	"net/http"
)

func CreateRole(c *gin.Context) {
	var R models.SysRole
	_ = c.ShouldBindJSON(&R)
	role, err := R.CreateRole()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "创建失败", "data": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建成功", "data": role})
	}
}

func GetRoleList(c *gin.Context) {
	var pageInfo page.InfoPage

	_ = c.BindQuery(&pageInfo)

	fmt.Println(pageInfo, "pageInfo")

	err, list, total := new(models.SysRole).GetList(pageInfo)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取角色数据失败, %v", err), c)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":     200,
			"data":     list,
			"total":    total,
			"page":     pageInfo.Page,
			"pageSize": pageInfo.PageSize,
		})
	}
	//response.FailWithMessage(fmt.Sprintf("获取角色数据失败, %v", "eeee"), c)
}
