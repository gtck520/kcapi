package controllers

import (
	"strings"

	"fmt"
	"github.com/astaxie/beego"
	"github.com/gtck520/kcapi/enums"
	"github.com/gtck520/kcapi/models"
	"github.com/gtck520/kcapi/utils"
	"time"
)

type LoginController struct {
	BaseController
}

// @Title do login
// @Description 根据用户名密码登录
// @Param	body		body 	models.BackendUser	true		"The object content"
// @Success 200 {string} models.Object.userinfo
// @Failure 403 username or password is empty
// @router /doLogin [post]
func (this *LoginController) doLogin() {
	remoteAddr := this.Ctx.Request.RemoteAddr
	addrs := strings.Split(remoteAddr, "::1")
	if len(addrs) > 1 {
		remoteAddr = "localhost"
	}

	username := strings.TrimSpace(this.GetString("UserName"))
	userpwd := strings.TrimSpace(this.GetString("UserPwd"))

	if err := models.LoginTraceAdd(username, remoteAddr, time.Now()); err != nil {
		beego.Error("LoginTraceAdd error.")
	}
	beego.Info(fmt.Sprintf("login: %s IP: %s", username, remoteAddr))

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

		//获取用户信息
		this.jsonResult(enums.JRCodeSucc, "登录成功", "")
	} else {
		this.jsonResult(enums.JRCodeFailed, "用户名或者密码错误", "")
	}
}

//退出
func (this *LoginController) Logout() {
	user := models.BackendUser{}
	this.SetSession("backenduser", user)
	//this.pageLogin()
}
