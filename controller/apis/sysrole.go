package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-admin/models"
	"go-admin/models/page"
	"go-admin/utils"
	"go-admin/utils/errors"
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
	err, list, total := new(models.SysRole).GetList(pageInfo)
	if err != nil {
		errors.FailWithMessage(fmt.Sprintf("获取角色数据失败, %v", err), c)
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

func GetRoleByQuery(c *gin.Context) {
	var r models.SysRole
	var rq models.RoleQuery
	r.ID = utils.StringConvUint(c.Param("id"))
	_ = c.BindQuery(&rq)
	err, result, total := r.GetRoleByQuery(rq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "查询失败", "data": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":     200,
			"data":     result,
			"total":    total,
			"page":     rq.Page,
			"pageSize": rq.PageSize,
		})
	}
}
