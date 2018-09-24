package models

import (

	"github.com/astaxie/beego/orm"
)

type Claims  struct {
	Id          int  `json:"jti,omitempty" orm:"column(id);pk"` // 设置主键
	Audience  string `json:"aud,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	Issuer    string `json:"iss,omitempty"`
	NotBefore int64  `json:"nbf,omitempty"`
	Subject   string `json:"sub,omitempty"`
	Userid   string `json:"usr,omitempty"`
}


func init() {
	orm.RegisterModel(new(Claims))
}

func (a *Claims) TableName() string {
	return ClaimsTable()
}

//获取 BackendUser 对应的表名称
func ClaimsTable() string {
	return "sys_claims"
}


