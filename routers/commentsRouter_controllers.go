package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

	beego.GlobalControllerRouter["mi-beego/controllers:CartController"] = append(beego.GlobalControllerRouter["mi-beego/controllers:CartController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           "/:productId",
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mi-beego/controllers:CartController"] = append(beego.GlobalControllerRouter["mi-beego/controllers:CartController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           "/:productId",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mi-beego/controllers:CartController"] = append(beego.GlobalControllerRouter["mi-beego/controllers:CartController"],
		beego.ControllerComments{
			Method:           "Add",
			Router:           "/add",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mi-beego/controllers:CartController"] = append(beego.GlobalControllerRouter["mi-beego/controllers:CartController"],
		beego.ControllerComments{
			Method:           "List",
			Router:           "/list",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mi-beego/controllers:CartController"] = append(beego.GlobalControllerRouter["mi-beego/controllers:CartController"],
		beego.ControllerComments{
			Method:           "SelectAll",
			Router:           "/selectAll",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mi-beego/controllers:CartController"] = append(beego.GlobalControllerRouter["mi-beego/controllers:CartController"],
		beego.ControllerComments{
			Method:           "GetCount",
			Router:           "/sum",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mi-beego/controllers:CartController"] = append(beego.GlobalControllerRouter["mi-beego/controllers:CartController"],
		beego.ControllerComments{
			Method:           "UnSelectAll",
			Router:           "/unSelectAll",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mi-beego/controllers:OrderInfoController"] = append(beego.GlobalControllerRouter["mi-beego/controllers:OrderInfoController"],
		beego.ControllerComments{
			Method:           "GetOne",
			Router:           "/:orderNo",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mi-beego/controllers:OrderInfoController"] = append(beego.GlobalControllerRouter["mi-beego/controllers:OrderInfoController"],
		beego.ControllerComments{
			Method:           "Create",
			Router:           "/create",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mi-beego/controllers:OrderInfoController"] = append(beego.GlobalControllerRouter["mi-beego/controllers:OrderInfoController"],
		beego.ControllerComments{
			Method:           "Pages",
			Router:           "/pages",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mi-beego/controllers:OrderInfoController"] = append(beego.GlobalControllerRouter["mi-beego/controllers:OrderInfoController"],
		beego.ControllerComments{
			Method:           "Pay",
			Router:           "/pay",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mi-beego/controllers:ProductController"] = append(beego.GlobalControllerRouter["mi-beego/controllers:ProductController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           "/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mi-beego/controllers:ProductController"] = append(beego.GlobalControllerRouter["mi-beego/controllers:ProductController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           "/:id",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mi-beego/controllers:ProductController"] = append(beego.GlobalControllerRouter["mi-beego/controllers:ProductController"],
		beego.ControllerComments{
			Method:           "GetOne",
			Router:           "/:productId",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mi-beego/controllers:ProductController"] = append(beego.GlobalControllerRouter["mi-beego/controllers:ProductController"],
		beego.ControllerComments{
			Method:           "Pages",
			Router:           "/pages",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mi-beego/controllers:UserAddressController"] = append(beego.GlobalControllerRouter["mi-beego/controllers:UserAddressController"],
		beego.ControllerComments{
			Method:           "GetOne",
			Router:           "/:addressId",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mi-beego/controllers:UserAddressController"] = append(beego.GlobalControllerRouter["mi-beego/controllers:UserAddressController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           "/:addressId",
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mi-beego/controllers:UserAddressController"] = append(beego.GlobalControllerRouter["mi-beego/controllers:UserAddressController"],
		beego.ControllerComments{
			Method:           "Add",
			Router:           "/add",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mi-beego/controllers:UserAddressController"] = append(beego.GlobalControllerRouter["mi-beego/controllers:UserAddressController"],
		beego.ControllerComments{
			Method:           "Pages",
			Router:           "/pages",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mi-beego/controllers:UserAddressController"] = append(beego.GlobalControllerRouter["mi-beego/controllers:UserAddressController"],
		beego.ControllerComments{
			Method:           "Update",
			Router:           "/update",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mi-beego/controllers:UserController"] = append(beego.GlobalControllerRouter["mi-beego/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           "/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mi-beego/controllers:UserController"] = append(beego.GlobalControllerRouter["mi-beego/controllers:UserController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           "/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mi-beego/controllers:UserController"] = append(beego.GlobalControllerRouter["mi-beego/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           "/:id",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mi-beego/controllers:UserController"] = append(beego.GlobalControllerRouter["mi-beego/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           "/:id",
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mi-beego/controllers:UserController"] = append(beego.GlobalControllerRouter["mi-beego/controllers:UserController"],
		beego.ControllerComments{
			Method:           "GetUser",
			Router:           "/getUser",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mi-beego/controllers:UserController"] = append(beego.GlobalControllerRouter["mi-beego/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Login",
			Router:           "/login",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mi-beego/controllers:UserController"] = append(beego.GlobalControllerRouter["mi-beego/controllers:UserController"],
		beego.ControllerComments{
			Method:           "LogOut",
			Router:           "/logout",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
