package controllers

import (
	"encoding/json"
	"mi-beego/models"
	"mi-beego/pkg/lib"
	"strconv"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

// ProductController operations for Product
type ProductController struct {
	beego.Controller
}

// URLMapping ...
func (c *ProductController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("Pages", c.Pages)
	c.Mapping("Put", c.Put)
}

// Post ...
// @Title Post
// @Description create Product
// @Param	body		body 	models.Product	true		"body for Product content"
// @Success 201 {int} models.Product
// @Failure 500 body is empty
// @router / [post]
func (c *ProductController) Post() {
	var v models.Product
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddProduct(&v); err == nil {
			c.Data["json"] = lib.Ok()
		} else {
			c.Data["json"] = lib.Err()
		}
	} else {
		c.Data["json"] = lib.Err()
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Product by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Product
// @Failure 500 :id is empty
// @router /:productId [get]
func (c *ProductController) GetOne() {
	idStr := c.Ctx.Input.Param(":productId")
	//productId, _ := strconv.ParseInt(idStr, 10, 64)
	v, err := models.GetProductByProductId(idStr)
	if err != nil {
		c.Data["json"] = lib.Err()
	} else {
		c.Data["json"] = lib.OkData(v)
	}
	c.ServeJSON()
}

// Pages ...
// @Title Get All
// @Description get Product
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Product
// @Failure 500
// @router /pages [get]
func (c *ProductController) Pages() {
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
				c.Data["json"] = lib.ErrMsg("invalid query key/value pair")
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
	l, count, err := models.GetAllProduct(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = lib.Err()
	} else {
		page := lib.Page{PageNo: offset, PageSize: limit, TotalCount: count, Records: l}
		c.Data["json"] = lib.OkData(page)
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Product
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Product	true		"body for Product content"
// @Success 200 {object} models.Product
// @Failure 500 :id is not int
// @router /:id [put]
func (c *ProductController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Product{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateProductById(&v); err == nil {
			c.Data["json"] = lib.Ok()
		} else {
			c.Data["json"] = lib.Err()
		}
	} else {
		c.Data["json"] = lib.Err()
	}
	c.ServeJSON()
}
