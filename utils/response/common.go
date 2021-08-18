package response

type Response struct {
	ErrorCode    int         `json:"errorCode"`
	Success      bool        `json:"success"`
	ErrorMessage string      `json:"errorMessage"`
	Timestamp    int64       `json:"timestamp"`
	ShowType     int         `json:"showType"`
	Data         interface{} `json:"data"`
	Host         string      `json:"host"`
}

type PageInfo struct {
	Response
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}
