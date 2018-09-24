package controllers
import (
	"encoding/json"
	"github.com/mojocn/base64Captcha"
	"github.com/gtck520/kcapi/enums"
	"github.com/gtck520/kcapi/models"
	"fmt"
	"math/rand"
	"time"
)
// 验证码模块
type CaptchaController struct {
	BaseController
}
//ConfigJsonBody json request body.
type ConfigJsonBody struct {
	Id              string
	CaptchaType     string
	VerifyValue     string
	ConfigAudio     base64Captcha.ConfigAudio
	ConfigCharacter base64Captcha.ConfigCharacter
	ConfigDigit     base64Captcha.ConfigDigit
}
// @Title get captcha
// @Description base64Captcha 创建图片验证码 参数：{"CaptchaType": "验证码类型"} 为空为默认验证码
// @Param	body		body 	ConfigJsonBody	true		"图像对象"
// @Success 200 {json}  ConfigJsonBody
// @Failure 403 fail
// @router /generateCaptcha [post]
func (this *CaptchaController)GenerateCaptchaHandler() {
	//parse request parameters
	//接收客户端发送来的请求参数
	var postParameters ConfigJsonBody
	json.Unmarshal(this.Ctx.Input.RequestBody, &postParameters)

	//config struct for digits
	//数字验证码配置
	var configD = base64Captcha.ConfigDigit{
		Height:     80,
		Width:      240,
		MaxSkew:    0.7,
		DotCount:   80,
		CaptchaLen: 5,
	}
	//config struct for audio
	//声音验证码配置
	var configA = base64Captcha.ConfigAudio{
		CaptchaLen: 6,
		Language:   "zh",
	}
	//config struct for Character
	//字符,公式,验证码配置
	var configC = base64Captcha.ConfigCharacter{
		Height:             47,
		Width:              100,
		//const CaptchaModeNumber:数字,CaptchaModeAlphabet:字母,CaptchaModeArithmetic:算术,CaptchaModeNumberAlphabet:数字字母混合.
		Mode:               base64Captcha.CaptchaModeNumber,
		ComplexOfNoiseText: base64Captcha.CaptchaComplexLower,
		ComplexOfNoiseDot:  base64Captcha.CaptchaComplexLower,
		IsShowHollowLine:   false,
		IsShowNoiseDot:     false,
		IsShowNoiseText:    false,
		IsShowSlimeLine:    false,
		IsShowSineLine:     false,
		CaptchaLen:         6,
	}

	//create base64 encoding captcha
	//创建 base64 图像验证码
	var config interface{}
	switch postParameters.CaptchaType {
	case "audio":
		config = configA
	case "character":
		config = configC
	default:
		config = configD
	}
	captchaId, digitCap := base64Captcha.GenerateCaptcha("", config)
	base64Png := base64Captcha.CaptchaWriteToBase64Encoding(digitCap)

	//or you can do this
	//你也可以是用默认参数 生成图像验证码
	//base64Png := base64Captcha..GenerateCaptchaPngBase64StringDefault(captchaId)
	body := map[string]interface{}{"data": base64Png, "captchaId": captchaId}
	this.jsonResult(enums.JRCodeSucc,"success", body)
}
// @Title get captcha
// @Description base64Captcha 图片验证码验证 参数：{"Id": "关键字","VerifyValue": "输入验证码"}
// @Param	body		body 	ConfigJsonBody	true		"图像对象"
// @Success 200 {json}  ConfigJsonBody
// @Failure 403 fail
// @router /captchaVerify [post]
func (this *CaptchaController)CaptchaVerifyHandle() {

	//parse request parameters
	//接收客户端发送来的请求参数
	var postParameters ConfigJsonBody
	json.Unmarshal(this.Ctx.Input.RequestBody, &postParameters)
	//verify the captcha
	//比较图像验证码
	verifyResult := base64Captcha.VerifyCaptcha(postParameters.Id, postParameters.VerifyValue)
	//set json response
	if verifyResult {
		this.jsonResult(enums.JRCodeSucc,"success", "")
	}else{
		this.jsonResult(enums.JRCodeFailed,"fail", "")
	}
}
// @Title get MobileCode
// @Description base64Captcha 获取手机验证码 参数：{"Mobile": "发送的电话号码"}
// @Param	body		body 	models.MobileLog	true		"图像对象"
// @Success 200 {json}  models.MobileLog.cod
// @Failure 403 fail
// @router /getMobileCode [post]
func (this *CaptchaController)MobileCode() {
	//接收客户端发送来的请求参数
	var postParameters models.MobileLog
	json.Unmarshal(this.Ctx.Input.RequestBody, &postParameters)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	if(postParameters.Mobile==""){
		this.jsonResult(enums.JRCodeFailed,"电话号码不能为空", "")
	}
    onelog,err:=postParameters.GetOnebyMobile(postParameters.Mobile)
	if onelog != nil && err == nil {
		t := time.Now()
		if(onelog.Expires<t.Unix()){
			postParameters.DeleteOnebyMobile(postParameters.Mobile)
			idkey:=this.GetRandomString(10)
			vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
			//将本次发送信息存入数据库
			models.MobileLogAdd(idkey,vcode,postParameters.Mobile)
			body := map[string]interface{}{"vcode": vcode, "key": idkey}
			this.jsonResult(enums.JRCodeSucc,"success", body)
		}else{
			//验证码已发送并且未过期
			body := map[string]interface{}{"vcode": onelog.Code, "key": onelog.Idkey}
			this.jsonResult(enums.JRCodeSucc,"successold", body)
		}
	}else{
		idkey:=this.GetRandomString(10)
		vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
		//将本次发送信息存入数据库
		models.MobileLogAdd(idkey,vcode,postParameters.Mobile)
		body := map[string]interface{}{"vcode": vcode, "key": idkey}
		this.jsonResult(enums.JRCodeSucc,"successfirst", body)
	}
}
// @Title check MobileCode
// @Description base64Captcha 手机验证码验证 参数：{"idkey": "关键字","code":"验证码"}
// @Param	body		body 	models.MobileLog	true		"图像对象"
// @Success 200 {json}  models.MobileLog.cod
// @Failure 403 fail
// @router /getMobileCode [post]
func (this *CaptchaController)CheckMobileCode() {
	//parse request parameters
	//接收客户端发送来的请求参数
	var postParameters models.MobileLog
	json.Unmarshal(this.Ctx.Input.RequestBody, &postParameters)
	onelog,err:=postParameters.GetOnebyMobile(postParameters.Idkey)
	if onelog != nil && err == nil {
		t := time.Now()
		if(onelog.Expires<t.Unix()){
			this.jsonResult(enums.JRCodeFailed,"手机验证码已过期，请重新获取", "")
		}
		if(onelog.Code==postParameters.Code){
			this.jsonResult(enums.JRCodeSucc,"success", onelog.Code)
		}else{
			this.jsonResult(enums.JRCodeFailed,"手机验证码错误", "")

		}
	}else{
		this.jsonResult(enums.JRCodeFailed,"验证码不匹配", "")
	}

}
