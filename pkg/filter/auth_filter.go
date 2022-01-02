package filter

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"mi-beego/models"
	"mi-beego/pkg/lib"
	"mi-beego/pkg/session"
	"mi-beego/pkg/utils"
	"mi-beego/third_party/redis"
)

// 判断header 投中是否有 Authorization  并且是否是有效的
var AuthUser = func(ctx *context.Context) {

	authorization := ctx.Request.Header.Get(lib.Authorization)
	if authorization == "" {
		ctx.Output.JSON(lib.LoginErr(), false, false)
	} else {
		userRedis := redis.Get(lib.UserLoginToken + authorization)
		if userRedis == "" {
			ctx.Output.JSON(lib.LoginErr(), false, false)
		} else {
			sessionExpire, _ := beego.AppConfig.Int("session_expire")
			redis.Set(lib.UserLoginToken+authorization, userRedis, sessionExpire)
			loginUser := models.LoginUser{}
			utils.JsonToObject(userRedis, &loginUser)
			session.GlobalMap.Store(authorization, loginUser)
		}
	}

}
