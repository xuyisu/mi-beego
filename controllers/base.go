package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web/context"
	"mi-beego/models"
	"mi-beego/pkg/lib"
	"mi-beego/pkg/session"
)

func GetLoginUser(ctx context.Context) models.LoginUser {
	authorization := ctx.Request.Header.Get(lib.Authorization)
	loginUser, _ := session.GlobalMap.Load(authorization)
	if loginUser != nil {
		return loginUser.(models.LoginUser)
	} else {
		logs.Error("未获取到用户")
		return models.LoginUser{}
	}
}
