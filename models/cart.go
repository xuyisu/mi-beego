package models

import (
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"mi-beego/pkg/lib"
	"reflect"
	"strings"
)

type Cart struct {
	Id                int64    `orm:"column(id);auto" description:"主键" json:"id,string"`
	CreateTime        lib.Time `orm:"column(create_time);type(datetime);auto_now_add" description:"创建时间" json:"-"`
	UpdateTime        lib.Time `orm:"column(update_time);type(datetime);auto_now" description:"更新时间" json:"-"`
	CreateUser        int64    `orm:"column(create_user)" description:"创建人" json:"-"`
	UpdateUser        int64    `orm:"column(update_user)" description:"更新人" json:"-"`
	DeleteFlag        int8     `orm:"column(delete_flag)" description:"删除标志" json:"-"`
	UserId            int64    `orm:"column(user_id);null" description:"用户id" json:"userId,string"`
	ActivityId        int64    `orm:"column(activity_id);null" description:"活动id" json:"activityId,string"`
	ActivityName      string   `orm:"column(activity_name);size(255);null" description:"活动名称" json:"activityName"`
	ProductId         int64    `orm:"column(product_id)" description:"商品id" json:"productId,string"`
	ProductName       string   `orm:"column(product_name);size(255)" description:"商品名称" json:"productName"`
	ProductSubtitle   string   `orm:"column(product_subtitle);size(255);null" description:"商品简要描述" json:"productSubtitle"`
	ProductMainImage  string   `orm:"column(product_main_image);size(255);null" description:"商品图片地址" json:"productMainImage"`
	Quantity          int      `orm:"column(quantity)" description:"数量" json:"quantity"`
	ProductUnitPrice  float64  `orm:"column(product_unit_price);digits(20);decimals(2)" description:"单价" json:"productUnitPrice"`
	Selected          int8     `orm:"column(selected)" description:"是否已选择 1是 0 否" json:"selected"`
	ProductTotalPrice float64  `orm:"column(product_total_price);digits(20);decimals(2)" description:"总价格" json:"productTotalPrice"`
}

type CartResp struct {
	//购物车总价
	CartTotalPrice string `json:"cartTotalPrice"`
	//总数量
	CartTotalQuantity int `json:"cartTotalQuantity"`
	//是否全选
	SelectedAll bool `json:"selectedAll"`

	//购物车列表
	CartProductList interface{} `json:"cartProductList"`
}

type CartReq struct {
	//数量
	//Quantity int `json:"quantity"`
	//是否选中
	Selected interface{} `json:"selected"`
	//类型
	Type interface{} `json:"type"`
}

func (t *Cart) TableName() string {
	return "cart"
}

func init() {
	orm.RegisterModel(new(Cart))
}

// AddCart insert a new Cart into database and returns
// last inserted Id on success.
func AddCart(m *Cart) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetCartById retrieves Cart by Id. Returns error if
// Id doesn't exist
func GetCartCount() (count int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Cart))
	count, err = qs.Filter("delete_flag", 0).Filter("user_id", 0).Count()
	if err != nil {
		logs.Error(err.Error())
	}
	return count, nil
}

// GetCartById retrieves Cart by Id. Returns error if
// Id doesn't exist
func GetCartByProductId(productId int64) (v *Cart, err error) {
	o := orm.NewOrm()
	v = &Cart{ProductId: productId, DeleteFlag: lib.Zero}
	if err = o.Read(v, "productId", "deleteFlag"); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllCart retrieves all Cart matches certain condition. Returns empty list if
// no records exist
func GetAllCart(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, cartRes []Cart, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Cart))
	qs = qs.Filter("delete_flag", 0)
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Cart
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, l, nil
	}
	return nil, nil, err
}

// UpdateCart updates Cart by Id and returns error if
// the record to be updated doesn't exist
func UpdateCartById(m *Cart) (err error) {
	o := orm.NewOrm()
	v := Cart{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteCart deletes Cart by Id and returns error if
// the record to be deleted doesn't exist
func DeleteCart(id int64) (err error) {
	o := orm.NewOrm()
	v := Cart{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Cart{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
