package swagger

import beego "github.com/beego/beego/v2/server/web"

func init() {
	//尽在dev  环境下启动swagger
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
}
