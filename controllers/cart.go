package controllers

import (
	"encoding/json"
	"errors"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/shopspring/decimal"
	"mi-beego/models"
	"mi-beego/pkg/lib"
	"strconv"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

// CartController operations for Cart
type CartController struct {
	beego.Controller
}

// URLMapping ...
func (c *CartController) URLMapping() {
	c.Mapping("Add", c.Add)
	c.Mapping("GetCount", c.GetCount)
	c.Mapping("List", c.List)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("SelectAll", c.SelectAll)
	c.Mapping("UnSelectAll", c.UnSelectAll)
}

// Add ...
// @Title Add
// @Description create Cart
// @Param	body		body 	models.Cart	true		"body for Cart content"
// @Success 201 {int} models.Cart
// @Failure 500 body is empty
// @router /add [post]
func (c *CartController) Add() {
	loginUser := GetLoginUser(*c.Ctx)
	var v models.Cart
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		product, _ := models.GetProductByProductId(v.ProductId)
		if product == nil {
			c.Data["json"] = lib.ErrMsg("当前商品已下架或删除")
			c.ServeJSON()
			return
		}
		cart, _ := models.GetCartByProductId(v.ProductId)
		o := orm.NewOrm()
		begin, err := o.Begin()
		if cart == nil {
			v.ProductName = product.Name
			v.ProductSubtitle = product.SubTitle
			v.ProductUnitPrice = product.Price
			v.ProductMainImage = product.MainImage
			v.Quantity = lib.One
			totalPrice := decimal.NewFromInt(int64(v.Quantity))
			productPrice := decimal.NewFromFloat(product.Price)
			v.ProductTotalPrice = totalPrice.Mul(productPrice).InexactFloat64()
			activity, _ := models.GetActivityByTime(time.Now())
			if activity != nil {
				v.ActivityName = activity.Name
				v.ActivityId = activity.ActivityId
			}
			v.Selected = lib.One
			v.CreateUser = loginUser.Id
			v.UserId = loginUser.Id
			v.UpdateUser = loginUser.Id
			if _, err := models.AddCart(&v); err == nil {
				count, _ := models.GetCartCount(loginUser.Id)
				c.Data["json"] = lib.OkData(count)
				begin.Commit()
			} else {
				logs.Error(err)
				c.Data["json"] = lib.ErrMsg("添加失败")
				begin.Rollback()
				return
			}
		} else {
			cart.UpdateTime = lib.Time{Time: time.Now()}
			cart.Quantity = cart.Quantity + lib.One
			totalPrice := decimal.NewFromInt(int64(cart.Quantity))
			productPrice := decimal.NewFromFloat(product.Price)
			cart.ProductTotalPrice = totalPrice.Mul(productPrice).InexactFloat64()
			if err = models.UpdateCartById(cart); err == nil {
				begin.Commit()
				count, _ := models.GetCartCount(loginUser.Id)
				c.Data["json"] = lib.OkData(count)
			} else {
				logs.Error(err)
				c.Data["json"] = lib.ErrMsg("更新失败")
				begin.Rollback()
				return
			}
		}

	} else {
		c.Data["json"] = lib.ErrMsg("解析失败")
	}
	c.ServeJSON()
}

// GetCount ...
// @Title Get Count
// @Description 查询购物车数量
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Cart
// @Failure 500 :id is empty
// @router /sum [get]
func (c *CartController) GetCount() {
	loginUser := GetLoginUser(*c.Ctx)
	v, err := models.GetCartCount(loginUser.Id)
	if err != nil {
		c.Data["json"] = lib.ErrMsg("查询失败")
	} else {
		c.Data["json"] = lib.OkData(v)
	}
	c.ServeJSON()
}

// @router /selectAll [put]
func (c *CartController) SelectAll() {

	query := make(map[string]string)
	loginUser := GetLoginUser(*c.Ctx)
	query["userId"] = strconv.FormatInt(loginUser.Id, 10)
	//这里只处理100个，可修改
	_, cartList, _ := models.GetAllCart(query, nil, []string{"id"}, []string{"desc"}, lib.Zero, lib.DefaultLimit)
	if cartList != nil {
		for _, v := range cartList {
			if v.Selected == lib.Zero {
				v.Selected = lib.One
				v.UpdateTime = lib.Time{Time: time.Now()}
				v.UpdateUser = lib.Zero
			}
			models.UpdateCartById(&v)
		}
		c.Data["json"] = lib.OkMsg("操作成功")
	} else {
		c.Data["json"] = lib.Err()
	}
	c.ServeJSON()
}

// @router /unSelectAll [put]
func (c *CartController) UnSelectAll() {
	query := make(map[string]string)
	loginUser := GetLoginUser(*c.Ctx)
	query["userId"] = strconv.FormatInt(loginUser.Id, 10)
	//这里只处理100个，可修改
	_, cartList, _ := models.GetAllCart(query, nil, []string{"id"}, []string{"desc"}, lib.Zero, lib.DefaultLimit)
	if cartList != nil {
		for _, v := range cartList {
			if v.Selected == lib.One {
				v.Selected = lib.Zero
				v.UpdateTime = lib.Time{Time: time.Now()}
				v.UpdateUser = lib.Zero
			}
			models.UpdateCartById(&v)
		}
		c.Data["json"] = lib.Ok()
	} else {
		c.Data["json"] = lib.Err()
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Cart
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Cart
// @Failure 500
// @router /list [get]
func (c *CartController) List() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 100
	var offset int64
	loginUser := GetLoginUser(*c.Ctx)
	query["user_id"] = strconv.FormatInt(loginUser.Id, 10)
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

	m1, cartRes, err := models.GetAllCart(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		totalQuantity := lib.Zero
		totalPrice, _ := decimal.NewFromString(".0")
		selectAll := true
		for _, v := range cartRes {
			//计算价格
			if v.Selected < lib.One {
				selectAll = false
			}
			totalQuantity += v.Quantity
			totalPrice = totalPrice.Add(decimal.NewFromFloat(v.ProductTotalPrice))
		}
		resp := models.CartResp{CartTotalPrice: totalPrice.String(), CartTotalQuantity: totalQuantity, SelectedAll: selectAll, CartProductList: m1}
		c.Data["json"] = lib.OkData(resp)
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Cart
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Cart	true		"body for Cart content"
// @Success 200 {object} models.Cart
// @Failure 500 :id is not int
// @router /:productId [put]
func (c *CartController) Put() {
	idStr := c.Ctx.Input.Param(":productId")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	product, _ := models.GetProductByProductId(id)
	if product == nil {
		c.Data["json"] = lib.ErrMsg("当前商品已下架或删除")
	}
	query := make(map[string]string)
	query["productId"] = idStr
	cartInterface, cartList, _ := models.GetAllCart(query, nil, []string{"id"}, []string{"desc"}, lib.Zero, lib.One)
	if cartInterface == nil {
		c.Data["json"] = lib.ErrMsg("该商品不在购物车")
	} else {

		v := cartList[lib.Zero]
		cartReq := models.CartReq{}
		if err := json.Unmarshal(c.Ctx.Input.RequestBody, &cartReq); err == nil {
			if cartReq.Type != nil {
				if int8(cartReq.Type.(float64)) == lib.One {
					v.Quantity = v.Quantity + lib.One
				} else {
					if v.Quantity <= lib.One {
						c.Data["json"] = lib.ErrMsg("不能再减了,要减没了")
						c.ServeJSON()
						return
					}
					v.Quantity = v.Quantity - lib.One
				}
			}
			if cartReq.Selected != nil {
				v.Selected = int8(cartReq.Selected.(float64))
			}
			totalPrice := decimal.NewFromInt(int64(v.Quantity))
			productPrice := decimal.NewFromFloat(product.Price)
			v.ProductTotalPrice = totalPrice.Mul(productPrice).InexactFloat64()
			if err := models.UpdateCartById(&v); err == nil {
				c.Data["json"] = lib.Ok()
			} else {
				c.Data["json"] = lib.ErrMsg("操作失败")
			}
		} else {
			c.Data["json"] = lib.ErrMsg("操作失败")
		}
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Cart
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 500 id is empty
// @router /:productId [delete]
func (c *CartController) Delete() {
	idStr := c.Ctx.Input.Param(":productId")
	var query map[string]string
	query = make(map[string]string)
	query["productId"] = idStr
	cartInterface, cartList, _ := models.GetAllCart(query, nil, []string{"id"}, []string{"desc"}, lib.Zero, lib.One)
	if cartInterface == nil {
		c.Data["json"] = lib.ErrMsg("该商品不在购物车")
	} else {
		loginUser := GetLoginUser(*c.Ctx)
		cart := cartList[lib.Zero]
		cart.DeleteFlag = lib.One
		cart.UpdateUser = loginUser.Id
		cart.UpdateTime = lib.Time{Time: time.Now()}
		if err := models.UpdateCartById(&cart); err == nil {
			c.Data["json"] = lib.Ok()
		} else {
			c.Data["json"] = lib.ErrMsg("删除失败")
		}
	}
	c.ServeJSON()
}
