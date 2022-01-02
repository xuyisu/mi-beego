package controllers

import (
	"encoding/json"
	"errors"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/shopspring/decimal"
	"mi-beego/models"
	"mi-beego/pkg/lib"
	"mi-beego/pkg/utils"
	"strconv"
	"strings"
	"time"
)

// OrderInfoController operations for OrderInfo
type OrderInfoController struct {
	beego.Controller
}

// URLMapping ...
func (c *OrderInfoController) URLMapping() {
	c.Mapping("Create", c.Create)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("Pages", c.Pages)
	c.Mapping("Pay", c.Pay)
}

// Create ...
// @Title create
// @Description create OrderInfo
// @Param	body		body 	models.OrderInfo	true		"body for OrderInfo content"
// @Success 201 {int} models.OrderInfo
// @Failure 500 body is error
// @router /create [post]
func (c *OrderInfoController) Create() {
	var v models.OrderInfo
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		address, err := models.GetUserAddressByAddrId(v.AddressId)
		if err != nil || address == nil {
			c.Data["json"] = lib.ErrMsg("当前地址已不存在，请重新添加地址")
			c.ServeJSON()
		}
		//根据需要改造
		var orderNo = strconv.FormatInt(time.Now().UnixNano()/1000, 10)
		totalOrderPrice, _ := decimal.NewFromString(".0")
		//查询购物车
		query := make(map[string]string)
		loginUser := GetLoginUser(*c.Ctx)
		query["userId"] = strconv.FormatInt(loginUser.Id, 10)
		//这里只处理100个，可修改
		_, cartList, _ := models.GetAllCart(query, nil, nil, nil, 0, 100)
		if len(cartList) == 0 {
			c.Data["json"] = lib.ErrMsg("恭喜您的购物车已经被清空了，再加一车吧")
			c.ServeJSON()
		}
		o := orm.NewOrm()
		begin, err := o.Begin()
		//遍历购物车
		for _, cart := range cartList {
			//查询商品
			product, _ := models.GetProductByProductId(cart.ProductId)
			if product.Stock <= 0 {
				c.Data["json"] = lib.ErrMsg("商品:" + cart.ProductName + " 已售尽,请选择其它产品")
				begin.Rollback()
				c.ServeJSON()
			}
			orderDetail := models.OrderDetail{CurrentUnitPrice: product.Price}

			//判断活动
			activity, _ := models.GetActivityByTime(time.Now())
			if activity != nil {
				orderDetail.ActivityId = activity.ActivityId
				orderDetail.ActivityName = activity.Name
				orderDetail.ActivityMainImage = activity.MainImage
			}
			orderDetail.OrderDetailNo = strconv.FormatInt(time.Now().UnixNano(), 10)
			orderDetail.OrderNo = orderNo
			orderDetail.ProductId = product.ProductId
			orderDetail.ProductMainImage = product.MainImage
			orderDetail.ProductName = product.Name
			orderDetail.Quantity = cart.Quantity
			orderDetail.Status = lib.PaymentStatueUnPay
			orderDetail.StatusDesc = lib.PaymentStatueUnPayDesc
			orderDetail.TotalPrice = cart.ProductTotalPrice
			orderDetail.UserId = loginUser.Id
			orderDetail.CreateUser = loginUser.Id
			totalOrderPrice = totalOrderPrice.Add(decimal.NewFromFloat(orderDetail.TotalPrice))
			models.AddOrderDetail(&orderDetail)
			cart.DeleteFlag = lib.One
			cart.UpdateTime = lib.Time{Time: time.Now()}
			cart.UpdateUser = loginUser.Id
			models.UpdateCartById(&cart)
			//设置订单状态记录
			statusRecord := models.OrderStatusRecord{
				OrderNo:       orderNo,
				OrderDetailNo: orderDetail.OrderDetailNo,
				ProductId:     orderDetail.ProductId,
				ProductName:   orderDetail.ProductName,
				Status:        orderDetail.Status,
				StatusDesc:    orderDetail.StatusDesc,
				CreateTime:    lib.Time{Time: time.Now()},
			}
			models.AddOrderStatusRecord(&statusRecord)
		}
		//订单主表
		orderInfo := models.OrderInfo{
			OrderNo:         orderNo,
			AddressId:       address.AddressId,
			Area:            address.Area,
			City:            address.City,
			Payment:         totalOrderPrice.InexactFloat64(),
			PaymentType:     lib.PaymentTypeOnline,
			PaymentTypeDesc: lib.PaymentTypeOnlineDesc,
			PostalCode:      address.PostalCode,
			Province:        address.Province,
			ReceiveName:     address.ReceiveName,
			ReceivePhone:    address.ReceivePhone,
			Street:          address.Street,
			Status:          lib.PaymentStatueUnPay,
			StatusDesc:      lib.PaymentStatueUnPayDesc,
			CreateUser:      loginUser.Id,
			UserId:          loginUser.Id,
		}
		if _, err := models.AddOrderInfo(&orderInfo); err == nil {
			c.Data["json"] = lib.OkData(orderNo)
			begin.Commit()
		} else {
			c.Data["json"] = lib.Err()
			begin.Rollback()
		}
	} else {
		logs.Error("解析异常:%v", err)
		c.Data["json"] = lib.Err()
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get OrderInfo by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.OrderInfo
// @Failure 500 :id is error
// @router /:orderNo [get]
func (c *OrderInfoController) GetOne() {
	orderNo := c.Ctx.Input.Param(":orderNo")
	v, err := models.GetOrderInfoByOrderNo(orderNo)
	if err != nil {
		logs.Error("查询失败,req=%v,res=%v", orderNo, err)
		c.Data["json"] = lib.Err()
	} else {
		loginUser := GetLoginUser(*c.Ctx)
		orderInfoVo := models.OrderInfoVo{}
		orderInfoVo.OrderInfo = *v
		query := make(map[string]string)
		query["user_id"] = strconv.FormatInt(loginUser.Id, 10)
		query["delete_flag"] = strconv.FormatInt(lib.Zero, 10)
		query["order_no"] = v.OrderNo
		//这里只处理100个，可修改
		_, orderDetailList, _ := models.GetAllOrderDetail(query, nil, []string{"id"}, []string{"desc"}, lib.Zero, lib.DefaultLimit)
		orderInfoVo.Details = orderDetailList
		c.Data["json"] = lib.OkData(orderInfoVo)
	}
	c.ServeJSON()
}

// Pages ...
// @Title Get All
// @Description get OrderInfo
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.OrderInfo
// @Failure 500 error
// @router /pages [get]
func (c *OrderInfoController) Pages() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	//只获取当前用户的
	loginUser := GetLoginUser(*c.Ctx)
	query["user_id"] = strconv.FormatInt(loginUser.Id, 10)
	query["delete_flag"] = strconv.FormatInt(lib.Zero, 10)
	sortby = append(sortby, "update_time")
	order = append(order, "desc")
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
	_, orderList, count, err := models.GetAllOrderInfo(query, fields, sortby, order, offset, limit)
	if err != nil {
		logs.Error("订单列表查询异常,req=%v", utils.ObjectToJson(query))
		c.Data["json"] = lib.ErrMsg("查询异常")
	} else {
		var orderInfoVos []models.OrderInfoVo
		if orderList != nil {
			orderInfoVo := models.OrderInfoVo{}
			for _, orderRes := range orderList {
				orderInfoVo.OrderInfo = orderRes
				query := make(map[string]string)
				query["user_id"] = strconv.FormatInt(loginUser.Id, 10)
				query["delete_flag"] = strconv.FormatInt(lib.Zero, 10)
				query["orderNo"] = orderRes.OrderNo
				//这里只处理100个，可修改
				_, orderDetailList, _ := models.GetAllOrderDetail(query, nil, []string{"id"}, []string{"desc"}, lib.Zero, lib.DefaultLimit)
				orderInfoVo.Details = orderDetailList
				orderInfoVos = append(orderInfoVos, orderInfoVo)
			}
		}
		if orderInfoVos == nil {
			orderInfoVos = make([]models.OrderInfoVo, 0, 0)
		}
		page := lib.Page{PageNo: offset, PageSize: limit, TotalCount: count, Records: orderInfoVos}
		c.Data["json"] = lib.OkData(page)
	}
	c.ServeJSON()
}

// Pay ...
// @Title Pay
// @Description update the OrderInfo
// @Param	body		body 	models.OrderInfo	true		"body for OrderInfo content"
// @Success 200 {object} models.OrderInfo
// @Failure 500  error
// @router /pay [post]
func (c *OrderInfoController) Pay() {
	param := models.PayReq{}
	loginUser := GetLoginUser(*c.Ctx)
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &param); err == nil {
		if order, _ := models.GetOrderInfoByOrderNo(param.OrderNo); err == nil {
			//待修改用户
			if order.UserId != loginUser.Id {
				c.Data["json"] = lib.ErrMsg("您无权查询他人订单")
				c.ServeJSON()
			}
			if order.Status != lib.PaymentStatueUnPay {
				c.Data["json"] = lib.ErrMsg("您没有待支付的订单")
				c.ServeJSON()
			}
			//开启事务
			o := orm.NewOrm()
			begin, _ := o.Begin()
			order.Status = lib.PaymentStatuePay
			order.StatusDesc = lib.PaymentStatuePayDesc
			date := lib.Time{Time: time.Now()}
			order.PaymentTime = date
			order.UpdateTime = date
			err := models.UpdateOrderInfoById(order)
			if err != nil {
				begin.Rollback()
			}
			//查询订单明细并更新
			loginUser := GetLoginUser(*c.Ctx)
			query := make(map[string]string)
			query["user_id"] = strconv.FormatInt(loginUser.Id, 10)
			query["delete_flag"] = strconv.FormatInt(lib.Zero, 10)
			query["order_no"] = order.OrderNo
			//这里只处理100个，可修改
			_, orderDetailList, _ := models.GetAllOrderDetail(query, nil, []string{"id"}, []string{"desc"}, lib.Zero, lib.DefaultLimit)

			for _, detail := range orderDetailList {
				detail.Status = lib.PaymentStatuePay
				detail.StatusDesc = lib.PaymentStatuePayDesc
				detail.UpdateTime = lib.Time{Time: time.Now()}
				detail.UpdateUser = loginUser.Id
				err := models.UpdateOrderDetailById(&detail)
				if err != nil {
					begin.Rollback()
				}
				//设置订单状态记录
				statusRecord := models.OrderStatusRecord{
					OrderNo:       detail.OrderNo,
					OrderDetailNo: detail.OrderDetailNo,
					ProductId:     detail.ProductId,
					ProductName:   detail.ProductName,
					Status:        detail.Status,
					StatusDesc:    detail.StatusDesc,
					CreateTime:    lib.Time{Time: time.Now()},
				}
				_, err = models.AddOrderStatusRecord(&statusRecord)
				if err != nil {
					begin.Rollback()
				}
				begin.Commit()
			}
			c.Data["json"] = lib.OkMsg("付款成功")
		} else {
			c.Data["json"] = lib.Err()

		}
	} else {
		c.Data["json"] = lib.Err()
	}
	c.ServeJSON()
}
