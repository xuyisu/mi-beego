package mysql

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	sqlConn, err := beego.AppConfig.String("sqlConn")
	if err != nil {
		logs.Error("异常:%v", err)
	}
	orm.RegisterDataBase("default", "mysql", sqlConn)
}
