package apis

import (
	"github.com/gin-gonic/gin"
	"go-admin/models"
	"go-admin/models/page"
	"net/http"
)

func CreateDept(c *gin.Context) {
	var D models.SysDept
	_ = c.BindJSON(&D)
	dep, err := D.Create()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "创建失败", "data": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建成功", "data": dep})
	}
}

func GetAll(c *gin.Context) {
	var pageInfo page.InfoPage
	_ = c.BindQuery(&pageInfo)
	err, list, total := new(models.SysDept).GetList(pageInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "data": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":     200,
			"data":     list,
			"total":    total,
			"page":     pageInfo.Page,
			"pageSize": pageInfo.PageSize,
		})
	}
}

func GetByUUID(c *gin.Context) {
	var D models.SysDept
	uid := c.Param("uuid")
	D.DeptID = uid
	if D, err := D.GetByUUID(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "data": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"dept": D,
			"msg":  "获取部门成功",
		})
	}
}

func GetDepTree(c *gin.Context) {
	dept, err := new(models.SysDept).SetDept()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "data": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"depTree": dept,
			"msg":     "获取组织树成功",
		})
	}
}
