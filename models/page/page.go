package page

type InfoPage struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

type Paging interface {
	GetList(InfoPage) (err error, list interface{}, total int)
}
