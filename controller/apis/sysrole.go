package apis

import (
	"github.com/gin-gonic/gin"
	"go-admin/models"
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

}
