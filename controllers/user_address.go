package controllers

import (
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"mi-beego/models"
	"mi-beego/pkg/lib"
	utils2 "mi-beego/pkg/utils"
	"strconv"
	_ "strconv"
	"strings"
	"time"
)

// UserAddressController operations for UserAddress
type UserAddressController struct {
	beego.Controller
}

// URLMapping ...
func (c *UserAddressController) URLMapping() {
	c.Mapping("Add", c.Add)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("Pages", c.Pages)
	c.Mapping("Update", c.Update)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create UserAddress
// @Param	body		body 	models.UserAddress	true		"body for UserAddress content"
// @Success 201 {int} models.UserAddress
// @Failure 500 body is empty
// @router /add [post]
func (c *UserAddressController) Add() {
	var v models.UserAddress
	loginUser := GetLoginUser(*c.Ctx)
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		//当你分布式的部署你的服务的时候，这个NewWorker的参数记录不同的node配置的值应该不一样
		worker, _ := utils2.NewWorker(1)
		v.AddressId = worker.GetId()
		v.UpdateTime = lib.Time{Time: time.Now()}
		v.UpdateUser = loginUser.Id
		v.CreateTime = lib.Time{Time: time.Now()}
		v.CreateUser = loginUser.Id
		if _, err := models.AddUserAddress(&v); err == nil {
			c.Data["json"] = lib.OkData(v)
		} else {
			logs.Error("查询失败:%v,req=%v", err.Error(), v)
			c.Data["json"] = lib.Err()
		}
	} else {
		logs.Error("查询失败:%v,req=%v", err.Error(), v)
		c.Data["json"] = lib.Err()
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get UserAddress by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.UserAddress
// @Failure 500 :id is empty
// @router /:addressId [get]
func (c *UserAddressController) GetOne() {
	idStr := c.Ctx.Input.Param(":addressId")
	addrId, _ := strconv.ParseInt(idStr, 10, 64)
	v, err := models.GetUserAddressByAddrId(addrId)
	if err != nil {
		logs.Error("查询失败:%v,req=%v", err.Error(), idStr)
		c.Data["json"] = lib.Err()
	} else {
		c.Data["json"] = lib.OkData(v)
	}
	c.ServeJSON()
}

// Pages ...
// @Title Get All
// @Description get UserAddress
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.UserAddress
// @Failure 500
// @router /pages [get]
func (c *UserAddressController) Pages() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64
	loginUser := GetLoginUser(*c.Ctx)
	query["create_user"] = strconv.FormatInt(loginUser.Id, 10)
	query["delete_flag"] = strconv.FormatInt(lib.Zero, 10)
	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("size"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("current"); err == nil {
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
				logs.Error("查询失败:%v,req=%v", "Error: invalid query key/value pair", v)
				c.Data["json"] = lib.ErrMsg("查询失败")
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
	l, count, err := models.GetAllUserAddress(query, fields, sortby, order, offset, limit)
	if err != nil {
		logs.Error("查询失败:%v", err.Error())
		c.Data["json"] = lib.ErrMsg("查询失败")
	} else {
		page := lib.Page{PageNo: offset, PageSize: limit, TotalCount: count, Records: l}
		c.Data["json"] = lib.OkData(page)
	}
	c.ServeJSON()
}

// Update ...
// @Title Update
// @Description update the UserAddress
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.UserAddress	true		"body for UserAddress content"
// @Success 200 {object} models.UserAddress
// @Failure 500 :id is not int
// @router /update [put]
func (c *UserAddressController) Update() {
	var v models.UserAddress
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		//当你分布式的部署你的服务的时候，这个NewWorker的参数记录不同的node配置的值应该不一样
		if err := models.UpdateUserAddressById(&v); err == nil {
			c.Data["json"] = lib.Ok()
		} else {
			logs.Error("更新:%v,req=%v", err.Error(), v)
			c.Data["json"] = lib.Err()
		}
	} else {
		logs.Error("更新失败:%v,req=%v", err.Error(), v)
		c.Data["json"] = lib.Err()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the UserAddress
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 500 id is empty
// @router /:addressId [delete]
func (c *UserAddressController) Delete() {
	idStr := c.Ctx.Input.Param(":addressId")
	addressId, _ := strconv.ParseInt(idStr, 10, 64)
	v, err := models.GetUserAddressByAddrId(addressId)
	if err != nil {
		logs.Error("查询失败:%v,req=%v", err.Error(), idStr)
		c.Data["json"] = lib.ErrMsg("无权限删除")
	}
	if err := models.DeleteUserAddress(v.Id); err == nil {
		c.Data["json"] = lib.Ok()
	} else {
		logs.Error("删除失败:%v,req=%v", err.Error(), idStr)
		c.Data["json"] = lib.Err()
	}
	c.ServeJSON()
}
