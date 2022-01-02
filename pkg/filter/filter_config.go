package filter

import (
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// 过滤器
	beego.InsertFilter("/api/cart/*", beego.BeforeRouter, AuthUser)
	beego.InsertFilter("/api/order/*", beego.BeforeRouter, AuthUser)
	beego.InsertFilter("/api/address/*", beego.BeforeRouter, AuthUser)
	beego.InsertFilter("/api/user/getUser", beego.BeforeRouter, AuthUser)
}
