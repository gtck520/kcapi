package main

import (
	_ "github.com/gtck520/kcapi/routers"
	_ "github.com/gtck520/kcapi/sysinit"

	"github.com/astaxie/beego"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"//配置可访问静态目录
	}
	beego.Run()
}
