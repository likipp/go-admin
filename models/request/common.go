package request

type PageInfo struct {
	Page     int `json:"page" form:"page"`         // 页码
	PageSize int `json:"pageSize" form:"pageSize"` // 每页大小
}

type GetById struct {
	ID float64 `json:"id" form:"id"`
}

type GetRoleId struct {
	RoleId string
}
