package data_type

type Page struct {
	PageNum int `form:"pageNum" json:"pageNum" binding:"min=1"` // 必填，页面值>=1
	Limit   int `form:"limit" json:"limit" binding:"min=1"`     // 必填，每页条数值>=1
}
