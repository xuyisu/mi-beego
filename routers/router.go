// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact yisu.martin@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"mi-beego/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/api",

		beego.NSNamespace("/cart",
			beego.NSInclude(
				&controllers.CartController{},
			),
		),

		beego.NSNamespace("/order",
			beego.NSInclude(
				&controllers.OrderInfoController{},
			),
		),

		beego.NSNamespace("/product",
			beego.NSInclude(
				&controllers.ProductController{},
			),
		),

		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),

		beego.NSNamespace("/address",
			beego.NSInclude(
				&controllers.UserAddressController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
