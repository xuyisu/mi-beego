package main

import (
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
	_ "mi-beego/pkg/filter"
	_ "mi-beego/pkg/log"
	_ "mi-beego/pkg/session"
	_ "mi-beego/routers"
	_ "mi-beego/third_party/mysql"
	_ "mi-beego/third_party/redis"
	_ "mi-beego/third_party/swagger"
)

func main() {
	beego.Run()
}
