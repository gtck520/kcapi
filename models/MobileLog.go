package models

import (

	"github.com/astaxie/beego/orm"
	"time"
)

type MobileLog  struct {
	Id          int  `json:"id,omitempty" orm:"column(id);pk"` // 设置主键
	Idkey  string `json:"idk,omitempty"`
	Mobile  string `json:"phone,omitempty"`
	Code   string  `json:"cod,omitempty"`
	CreatTime    int64 `json:"ctime,omitempty"`
	Expires int64  `json:"exp,omitempty"`
}


func init() {
	orm.RegisterModel(new(MobileLog))
}

func (a *MobileLog) TableName() string {
	return MobileLogTable()
}

//获取 BackendUser 对应的表名称
func MobileLogTable() string {
	return "sys_mobilelog"
}
//添加记录
func MobileLogAdd(_key string, _code string, _mobile string) error {
	t := time.Now()
	addTime:=t.Unix()
	//过期时间为30分钟
	expTime:=addTime+1800
	m := MobileLog{Idkey: _key, Code: _code, Mobile: _mobile,CreatTime:addTime,Expires:expTime}
	o := orm.NewOrm()
	if _, err := o.Insert(&m); err == nil {
		return nil
	} else {
		return err
	}
}
//根据电话号码读取一条数据
func (a *MobileLog)GetOnebyMobile(mobile string) (*MobileLog, error) {
	o := orm.NewOrm()
	m := MobileLog{}
	err := o.QueryTable(MobileLogTable()).Filter("mobile", mobile).One(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
//根据关机字取一条数据
func (a *MobileLog)GetOnebyIdkey(idkey string) (*MobileLog, error) {
	o := orm.NewOrm()
	m := MobileLog{}
	err := o.QueryTable(MobileLogTable()).Filter("idkey", idkey).One(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
//删除指定号码的存储
func (a *MobileLog)DeleteOnebyMobile(mobile string) error {
	query := orm.NewOrm().QueryTable(MobileLogTable())
	_, err := query.Filter("mobile", mobile).Delete()
	return  err
}


