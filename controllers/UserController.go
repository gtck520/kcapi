package controllers

import (
	"strings"
	"strconv"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/gtck520/kcapi/enums"
	"github.com/gtck520/kcapi/models"
	"github.com/gtck520/kcapi/utils"
	"time"
	"encoding/json"
)

type UserController struct {
	BaseController
}

// @Title do login
// @Description 根据用户名密码登录 {"Username":"string","Userpass":"string","Idkey":"string","Code":"string"}
// @Param	body  body 	LoginParam {"Username":"string","Userpass":"string","Idkey":"string","Code":"string"}	 true "The object content"
// @Success 200 {string} models.Object.userinfo
// @Failure 403 username or password is empty
// @router /doLogin [post]
func (this *UserController) DoLogin() {
	remoteAddr := this.Ctx.Request.RemoteAddr
	addrs := strings.Split(remoteAddr, "::1")
	if len(addrs) > 1 {
		remoteAddr = "localhost"
	}

	var postParameters models.LoginParam
	json.Unmarshal(this.Ctx.Input.RequestBody, &postParameters)
	username := strings.TrimSpace(postParameters.Username)
	userpwd := strings.TrimSpace(postParameters.Userpass)
	if err := models.LoginTraceAdd(username, remoteAddr, time.Now()); err != nil {
		beego.Error("LoginTraceAdd error.")
	}
	beego.Info(fmt.Sprintf("login: %s IP: %s", username, remoteAddr))

	//先验证手机验证码
	msg,status:=this.checkMobileCode(postParameters.Idkey,postParameters.Code)
   if(status==false){
	   this.jsonResult(enums.JRCodeFailed, msg, "")
   }
	if len(username) == 0 || len(userpwd) == 0 {
		this.jsonResult(enums.JRCodeFailed, "用户名和密码不正确", "")
	}

	userpwd = utils.String2md5(userpwd)
	user, err := models.BackendUserOneByUserName(username, userpwd)
	if user != nil && err == nil {
		if user.Status == enums.Disabled {
			this.jsonResult(enums.JRCodeFailed, "用户被禁用，请联系管理员", "")
		}
		//保存用户信息到session
		this.setBackendUser2Session(user.Id)
		utoken,_:=this.getUserToken(strconv.Itoa(user.Id))
		//获取用户信息
		userinfo := map[string]interface{}{"utoken": utoken, "adminid": user.Id,"username":user.UserName}
		this.jsonResult(enums.JRCodeSucc, "登录成功", userinfo)
	} else {
		this.jsonResult(enums.JRCodeFailed, "用户名或者密码错误", "")
	}
}
// @Title do login
// @Description 根据用户id获取信息 {"Userid":"int"}
// @Param	body  body 	userParam {"Userid":"int"}	 true "The object content"
// @Success 200 {string} models.Object.userinfo
// @Failure 403 username or password is empty
// @router /getUserInfo [post]
func (this *UserController) GetUserInfo(){
	var postParameters struct {
		Userid string
	}
	json.Unmarshal(this.Ctx.Input.RequestBody, &postParameters)
	beego.Info(fmt.Sprintf("asdfasdfsdf: %s", postParameters.Userid))
	userid,_:=strconv.Atoi(postParameters.Userid)
	userinfo,err := models.BackendUserOne(userid)
	if err != nil {
		this.jsonResult(enums.JRCodeFailed, "获取用户信息失败", err)
	}
	this.jsonResult(enums.JRCodeSucc, "获取用户信息成功", userinfo)
}

//退出
func (this *UserController) Logout() {
	user := models.BackendUser{}
	this.SetSession("backenduser", user)
	//this.pageLogin()
}
