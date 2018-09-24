// @APIVersion 1.0.1233333
// @Title konger basic API
// @Description 一个基础的后台框架api
// @Contact 496317580@qq.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/gtck520/kcapi/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/access",
			beego.NSInclude(
				&controllers.AccessController{},
			),
		),
		beego.NSNamespace("/captcha",
			beego.NSInclude(
				&controllers.CaptchaController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
