package page

type InfoPage struct {
	// antd中需要current与pageSize字段 所以此处后期改变了
	Page     int `form:"current"`
	PageSize int `form:"pageSize"`
}

type Paging interface {
	GetList(InfoPage) (err error, list interface{}, total int)
}
