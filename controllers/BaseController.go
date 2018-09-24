package controllers

import (
	"fmt"
	"strings"

	"github.com/gtck520/kcapi/enums"
	"github.com/gtck520/kcapi/models"
	"github.com/gtck520/kcapi/utils"

	"math/rand"
	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

type BaseController struct {
	beego.Controller
	controllerName string             //当前控制名称
	actionName     string             //当前action名称
	curUser        models.BackendUser //当前用户信息
}
//检查token是否有效
func (this *BaseController) checkToken(tokenString string,typea int)(string,error) {
	mysigningkey := beego.AppConfig.String("jvt::mysigningkey")
	if(typea == 1){
		//issuer := beego.AppConfig.String("jvt::issuera")
	}else{
		//issuer := beego.AppConfig.String("jvt::issueru")
	}
	hmacSampleSecret := []byte(mysigningkey)
	// sample token string taken from the New example
	//tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.u1riaD1rW97opCoAuRCTy4w58Br-Zk-bh7vLiRIsrpU"
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error)
	 token,err:= jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//fmt.Println(claims["foo"], claims["nbf"])

		beego.Info(fmt.Sprintf("调试变量输出 spend: %v", claims))
		return claims["sub"].(string), nil
	} else {
		//fmt.Println(err)

		beego.Info(fmt.Sprintf("token检查出错了 spend: %v", err))
		return "token检查出错了", err
	}
	return "呀，token检查出错了", nil
}

func (this *BaseController) Prepare() {
	//附值
	this.controllerName, this.actionName = this.GetControllerAndAction()

	this.Data["siteApp"] = beego.AppConfig.String("site.app")
	this.Data["siteName"] = beego.AppConfig.String("site.name")
	this.Data["siteVersion"] = beego.AppConfig.String("site.version")

	//从Session里获取数据 设置用户信息
	this.adapterUserInfo()
}

//根据用户的app信息生成token并存入数据库
func (this *BaseController) getToken(appid, appscret string) ( string , error) {
	api := models.Apiidentity{}
	apidata, err := api.GetOneByAppid(appid, appscret)
	if apidata != nil && err == nil {

		mysigningkey := beego.AppConfig.String("jvt::mysigningkey")
		issuer := beego.AppConfig.String("jvt::issuera")
		expiresats := beego.AppConfig.String("jvt::expiresat")
		expiresat,_  := strconv.Atoi(expiresats)
		mySigningKey := []byte(mysigningkey)
		// Create the Claims
		claims := &jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000),
			ExpiresAt: int64(time.Now().Unix() + int64(expiresat)),
			Issuer:    issuer,
			Subject:   strings.TrimSpace(apidata.AppId),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		ss, err := token.SignedString(mySigningKey)
		//beego.Info(fmt.Sprintf("调试变量输出 spend: %d s", ss))

		if err != nil {
			return "信息检验不正确",err
		}else{
     		return ss,nil
		}
	}else{
		return "用户名或密码不正确",err
	}
}

//从session里取用户信息
func (this *BaseController) adapterUserInfo() {
	a := this.GetSession("backenduser")
	if a != nil {
		this.curUser = a.(models.BackendUser)
		this.Data["backenduser"] = a
	}
}

// checkLogin判断用户是否登录，未登录则跳转至登录页面
// 一定要在BaseController.Prepare()后执行
func (this *BaseController) checkLogin() {
	if this.curUser.Id == 0 {
		//登录页面地址
		urlstr := this.URLFor("HomeController.Login") + "?url="

		//登录成功后返回的址为当前
		returnURL := this.Ctx.Request.URL.Path

		//如果ajax请求则返回相应的错码和跳转的地址
		if this.Ctx.Input.IsAjax() {
			//由于是ajax请求，因此地址是header里的Referer
			returnURL := this.Ctx.Input.Refer()
			this.jsonResult(enums.JRCode302, "请登录", urlstr+returnURL)
		}
		this.Redirect(urlstr+returnURL, 302)
		this.StopRun()
	}
}

// 判断某 Controller.Action 当前用户是否有权访问
func (this *BaseController) checkActionAuthor(ctrlName, ActName string) bool {
	if this.curUser.Id == 0 {
		return false
	}

	//从session获取用户信息
	user := this.GetSession("backenduser")

	//类型断言
	v, ok := user.(models.BackendUser)
	if ok {
		//如果是超级管理员，则直接通过
		if v.IsSuper == true {
			return true
		}

		//遍历用户所负责的资源列表
		for i, _ := range v.ResourceUrlForList {
			urlfor := strings.TrimSpace(v.ResourceUrlForList[i])
			if len(urlfor) == 0 {
				continue
			}
			// TestController.Get,:last,xie,:first,asta
			strs := strings.Split(urlfor, ",")
			if len(strs) > 0 && strs[0] == (ctrlName+"."+ActName) {
				return true
			}
		}
	}
	return false
}

// checkLogin判断用户是否有权访问某地址，无权则会跳转到错误页面
//一定要在BaseController.Prepare()后执行
// 会调用checkLogin
// 传入的参数为忽略权限控制的Action
func (this *BaseController) checkAuthor(ignores ...string) {
	//先判断是否登录
	this.checkLogin()

	//如果Action在忽略列表里，则直接通用
	for _, ignore := range ignores {
		if ignore == this.actionName {
			return
		}
	}

	hasAuthor := this.checkActionAuthor(this.controllerName, this.actionName)
	if !hasAuthor {
		utils.LogDebug(fmt.Sprintf("author control: path=%s.%s userid=%v  无权访问", this.controllerName, this.actionName, this.curUser.Id))

		//如果没有权限
		if !hasAuthor {
			if this.Ctx.Input.IsAjax() {
				this.jsonResult(enums.JRCode401, "无权访问", "")
			} else {

			}
		}
	}
}

//SetBackendUser2Session 获取用户信息（包括资源UrlFor）保存至Session
//被 HomeController.DoLogin 调用
func (this *BaseController) setBackendUser2Session(userId int) error {
	m, err := models.BackendUserOne(userId)
	if err != nil {
		return err
	}

	//获取这个用户能获取到的所有资源列表
	resourceList := models.ResourceTreeGridByUserId(userId, 1000)
	for _, item := range resourceList {
		m.ResourceUrlForList = append(m.ResourceUrlForList, strings.TrimSpace(item.UrlFor))
	}
	this.SetSession("backenduser", *m)
	return nil
}

// 设置模板
func (this *BaseController) jsonResult(code enums.JsonResultCode, msg string, obj interface{}) {
	res := &models.JsonResult{code, msg, obj}
	this.Data["json"] = res
	this.ServeJSON()
	this.StopRun()
}
//获取随机字符串
func (this *BaseController) GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

