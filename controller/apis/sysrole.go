package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-admin/models"
	"go-admin/models/page"
	"go-admin/utils"
	"go-admin/utils/response"
	"net/http"
)

func CreateRole(c *gin.Context) {
	var R models.SysRole
	_ = c.ShouldBindJSON(&R)
	role, err := R.CreateRole()
	if err != nil {
		response.Result(http.StatusBadRequest, nil, "创建用户失败", 0, false, c)
	} else {
		response.Result(http.StatusOK, role, "创建用户成功", 0, true, c)
	}
}

func GetRoleList(c *gin.Context) {
	var pageInfo page.InfoPage
	_ = c.BindQuery(&pageInfo)
	fmt.Println(pageInfo, "pageInfo")
	err, list, total := new(models.SysRole).GetList(pageInfo)
	fmt.Println(list, "list")
	if err != nil {
		response.Result(http.StatusBadRequest, nil, "获取角色数据失败", 0, false, c)
	} else {
		response.ResultWithPageInfo(list, "获取数据成功", 0, true, total, pageInfo.Page, pageInfo.PageSize, c)
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
		response.Result(http.StatusBadRequest, nil, "获取角色数据失败", 0, false, c)
	} else {
		response.ResultWithPageInfo(result, "获取数据成功", 0, true, total, rq.Page, rq.PageSize, c)
	}
}
