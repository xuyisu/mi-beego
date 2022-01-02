package controllers

import (
	"encoding/json"
	"errors"
	"github.com/beego/beego/v2/core/logs"
	"mi-beego/models"
	"mi-beego/pkg/lib"
	"mi-beego/pkg/session"
	"mi-beego/pkg/utils"
	"mi-beego/third_party/redis"
	"strconv"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

// UserController operations for User
type UserController struct {
	beego.Controller
}

// URLMapping ...
func (c *UserController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Login", c.Login)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("GetUser", c.GetUser)
	c.Mapping("LogOut", c.LogOut)
}

// Post ...
// @Title Post
// @Description create User
// @Param	body		body 	models.User	true		"body for User content"
// @Success 201 {int} models.User
// @Failure 500 body is empty
// @router / [post]
func (c *UserController) Post() {
	var v models.User
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		_, err := models.GetUserByUserNameAndPwd(v.UserName, utils.Md5(v.Password))
		if err == nil {
			c.Data["json"] = lib.ErrMsg("当前用户已存在")
			c.ServeJSON()
			return
		}

		if _, err := models.AddUser(&v); err == nil {
			c.Data["json"] = lib.Ok()
		} else {
			c.Data["json"] = lib.Err()
		}
	} else {
		c.Data["json"] = lib.Err()
	}
	c.ServeJSON()
}

// Login ...
// @Title Get One
// @Description login
// @Param	body		body 	models.User	true		"body for User content"
// @Success 200 {object} models.User
// @Failure 500 :id is empty
// @router /login [post]
func (c *UserController) Login() {
	var v models.LoginReq
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {

		if v.UserName == "" || v.Password == "" {
			c.Data["json"] = lib.ErrMsg("用户名密码必填")
			c.ServeJSON()
			return
		}

		user, err := models.GetUserByUserNameAndPwd(v.UserName, utils.Md5(v.Password))
		if err != nil {
			c.Data["json"] = lib.ErrMsg("登录失败,请检查用户名密码")
			c.ServeJSON()
			return
		}
		loginUser := models.LoginUser{
			Id:       user.Id,
			UserName: user.UserName,
			Email:    user.Email,
			Phone:    user.Phone,
		}
		//生成令牌
		token := utils.BuildToken()
		userMap := make(map[string]interface{})
		userMap[lib.Authorization] = token
		userMap["userInfo"] = loginUser
		//缓存令牌
		sessionExpire, _ := beego.AppConfig.Int("session_expire")
		redis.Set(lib.UserLoginToken+token, utils.ObjectToJson(loginUser), sessionExpire)
		c.Data["json"] = lib.OkData(userMap)
	} else {
		c.Data["json"] = lib.Err()
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get User
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.User
// @Failure 500
// @router / [get]
func (c *UserController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}
	if offset > 0 {
		offset--
	}
	l, err := models.GetAllUser(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = lib.Err()
	} else {
		c.Data["json"] = lib.OkData(l)
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the User
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.User	true		"body for User content"
// @Success 200 {object} models.User
// @Failure 500 :id is not int
// @router /:id [put]
func (c *UserController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 64, 10)
	v := models.User{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateUserById(&v); err == nil {
			c.Data["json"] = lib.Ok()
		} else {
			c.Data["json"] = lib.Err()
		}
	} else {
		c.Data["json"] = lib.Err()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the User
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 500 id is empty
// @router /:id [delete]
func (c *UserController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 64, 10)
	if err := models.DeleteUser(id); err == nil {
		c.Data["json"] = lib.Ok()
	} else {
		c.Data["json"] = lib.Err()
	}
	c.ServeJSON()
}

// GetUser ...
// @Title GetUser
// @Description get currentUser
// @Success 200 {string} delete success!
// @Failure 500 id is empty
// @router /getUser [get]
func (c *UserController) GetUser() {
	load, _ := session.GlobalMap.Load(c.Ctx.Request.Header.Get(lib.Authorization))
	c.Data["json"] = lib.OkData(load)
	c.ServeJSON()
}

// LogOut ...
// @Title LogOut
// @Description LogOut
// @Success 200 {string} logout success!
// @Failure 500 id is empty
// @router /logout [post]
func (c *UserController) LogOut() {
	c.Data["json"] = lib.Ok()
	authorization := c.Ctx.Request.Header.Get(lib.Authorization)
	if authorization != "" {
		redis.Delete(lib.UserLoginToken + authorization)
		session.GlobalMap.Delete(authorization)
		logs.Info("退出成功")
	}
	c.ServeJSON()
}
