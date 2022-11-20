package auth

import "goskeleton/app/model"

func CreateCasbinRuleFactory(sqlType string) *CasbinRule {
	return &CasbinRule{BaseModel: model.BaseModel{DB: model.UseDbConn(sqlType)}}
}

type CasbinRule struct {
	model.BaseModel `json:"-"`
	Ptype           string ` json:"ptype"`
	V0              string `json:"v0"`
	V1              string `json:"v1"`
	V2              string `json:"v2"`
	V3              string `json:"v3"`
	V4              string `json:"v4"`
	V5              string `json:"v5"`
}

// TableName 表名
func (c *CasbinRule) TableName() string {
	return "tb_auth_casbin_rule"
}
