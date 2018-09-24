package models

import (
	"github.com/astaxie/beego/orm"
)

type Apiidentity struct {
	Id          int      `orm:"column(id);pk"` // 设置主键
	AppId       string   `json:"appId" orm:"column(appId)"`
	AppSecret   string   `json:"appSecret" orm:"column(appSecret)"`
	AppName     string   `json:"appName" orm:"column(appName)"`
	DealLine    int   `json:"dealline" orm:"column(dealline)"`
	Create_time string   `json:"create_time" orm:"column(create_time)"`
}

type ApiidentityQueryParam struct {
	appId string //模糊查询
	appSecret string //模糊查询
}

func init() {
	orm.RegisterModel(new(Apiidentity))
}

func (a *Apiidentity) TableName() string {
	return ApiidentityTable()
}

//获取 BackendUser 对应的表名称
func ApiidentityTable() string {
	return "sys_api_identity"
}

// 根据appid 与 appsecret
func (a *Apiidentity) GetOneByAppid(appid,appsecret string) (*Apiidentity, error) {
	m := Apiidentity{}
	err := orm.NewOrm().QueryTable(ApiidentityTable()).Filter("appId",appid).Filter("appSecret",appsecret).One(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
