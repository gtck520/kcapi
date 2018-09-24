package controllers

import (
	"strings"
	"github.com/gtck520/kcapi/enums"
	"github.com/astaxie/beego"
	"encoding/json"
	"github.com/gtck520/kcapi/models"
	"fmt"
)

// 接口验证模块
type AccessController struct {
	BaseController
}
type GetTokens struct{
	TokenString string
}
// @Title get token
// @Description 根据客户端appid与秘钥获取token
// @Param	body		body 	models.Apiidentity	true		"The object content"
// @Success 200 {string} models.Object.app_token
// @Failure 403 appid or appsecret is empty
// @router /getToken [post,get]
func (this *AccessController) GetToken() {
	remoteAddr := this.Ctx.Request.RemoteAddr
	addrs := strings.Split(remoteAddr, "::1")
	if len(addrs) > 1 {
		remoteAddr = "localhost"
	}

	var access models.Apiidentity
	json.Unmarshal(this.Ctx.Input.RequestBody, &access)
	beego.Info(fmt.Sprintf("调试变量输出 spend: %d s", access))
	appid := strings.TrimSpace(access.AppId)
	appsecret := strings.TrimSpace(access.AppSecret)
	tokendata,err:=this.getToken(appid,appsecret)
	if err != nil {
		this.jsonResult(enums.JRCodeFailed, tokendata,err)
	}
	this.jsonResult(enums.JRCodeSucc, "成功",tokendata)
}
// @Title check token
// @Description 根据客户端token验证 参数：{TokenString："需要验证的token"}
// @Param	token  body  GetTokens  true "The object content"
// @Success 200 {string} models.Object.app_token
// @Failure 403 appid or appsecret is empty
// @router /checkToken [post,get]
func (this *AccessController) CheckToken() {
	var post GetTokens
	json.Unmarshal(this.Ctx.Input.RequestBody, &post)
	tokenstring:=strings.TrimSpace(post.TokenString)
	appinfo,err:=this.checkToken(tokenstring,1)

	if err != nil {
		this.jsonResult(enums.JRCodeFailed, "出错啦",err)
	}
	this.jsonResult(enums.JRCodeSucc, "成功",appinfo)
}



