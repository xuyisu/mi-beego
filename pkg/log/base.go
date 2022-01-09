package log

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	//log.SetLogger(logs.AdapterMultiFile, `{"filename":"go-BeeGo.log"}`)
	//log.EnableFuncCallDepth(true)
	//异步输出日志
	log.Async()
	//开启sql 日志
	if beego.BConfig.RunMode == "dev" {
		orm.Debug = true
	}
}
