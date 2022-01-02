package models

import (
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"mi-beego/pkg/lib"
	"reflect"
	"strings"
)

type OrderInfo struct {
	Id              int      `orm:"column(id);auto" description:"主键" json:"id"`
	CreateTime      lib.Time `orm:"column(create_time);type(datetime);auto_now_add" description:"创建时间" json:"createTime"`
	UpdateTime      lib.Time `orm:"column(update_time);type(datetime);auto_now" description:"更新时间" json:"-"`
	CreateUser      int64    `orm:"column(create_user)" description:"创建人" json:"-"`
	UpdateUser      int64    `orm:"column(update_user)" description:"更新人" json:"-"`
	DeleteFlag      int8     `orm:"column(delete_flag)" description:"删除标志" json:"-"`
	OrderNo         string   `orm:"column(order_no);size(60)" description:"订单编号" json:"orderNo"`
	Payment         float64  `orm:"column(payment);null;digits(20);decimals(2)" description:"支付金额" json:"payment"`
	PaymentType     int8     `orm:"column(payment_type);null" description:"支付类型" json:"paymentType"`
	PaymentTypeDesc string   `orm:"column(payment_type_desc);size(20);null" description:"支付类型描述" json:"paymentTypeDesc"`
	Postage         float64  `orm:"column(postage);null;digits(20);decimals(2)" description:"邮费" json:"postage"`
	Status          int8     `orm:"column(status)" description:"订单状态" json:"status"`
	StatusDesc      string   `orm:"column(status_desc);size(20)" description:"状态描述" json:"statusDesc"`
	PaymentTime     lib.Time `orm:"column(payment_time);type(datetime);null" description:"支付时间" json:"paymentTime"`
	AddressId       string   `orm:"column(address_id);null" description:"地址id" json:"addressId"`
	ReceiveName     string   `orm:"column(receive_name);size(50);null" description:"收货人" json:"receiveName"`
	ReceivePhone    string   `orm:"column(receive_phone);size(20);null" description:"联系号码" json:"receivePhone"`
	Province        string   `orm:"column(province);size(20);null" description:"省份" json:"province"`
	City            string   `orm:"column(city);size(20);null" description:"城市" json:"city"`
	Area            string   `orm:"column(area);size(20);null" description:"区" json:"area"`
	Street          string   `orm:"column(street);size(50);null" description:"详细地址" json:"street"`
	PostalCode      string   `orm:"column(postal_code);size(255);null" description:"邮编" json:"postalCode"`
	UserId          int64    `orm:"column(user_id)" description:"购买人id" json:"userId"`
}

//订单发那会列表新对象
type OrderInfoVo struct {
	OrderInfo
	Details []OrderDetail `json:"details"`
}

//支付请求
type PayReq struct {
	OrderNo string `json:"orderNo"`
	PayTool int    `json:"payTool"`
}

func (t *OrderInfo) TableName() string {
	return "order_info"
}

func init() {
	orm.RegisterModel(new(OrderInfo))
}

// AddOrderInfo insert a new OrderInfo into database and returns
// last inserted Id on success.
func AddOrderInfo(m *OrderInfo) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetOrderInfoById retrieves OrderInfo by Id. Returns error if
// Id doesn't exist
func GetOrderInfoByOrderNo(orderNo string) (v *OrderInfo, err error) {
	o := orm.NewOrm()
	v = &OrderInfo{OrderNo: orderNo, DeleteFlag: 0}
	if err = o.Read(v, "orderNo", "deleteFlag"); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllOrderInfo retrieves all OrderInfo matches certain condition. Returns empty list if
// no records exist
func GetAllOrderInfo(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, orderRes []OrderInfo, count int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(OrderInfo))
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
					return nil, nil, 0, errors.New("Error: Invalid order. Must be either [asc|desc]")
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
					return nil, nil, 0, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, nil, 0, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, nil, 0, errors.New("Error: unused 'order' fields")
		}
	}
	count, _ = qs.Count()
	var l []OrderInfo
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
		return ml, l, count, nil
	}
	return nil, nil, 0, err
}

// UpdateOrderInfo updates OrderInfo by Id and returns error if
// the record to be updated doesn't exist
func UpdateOrderInfoById(m *OrderInfo) (err error) {
	o := orm.NewOrm()
	v := OrderInfo{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteOrderInfo deletes OrderInfo by Id and returns error if
// the record to be deleted doesn't exist
func DeleteOrderInfo(id int) (err error) {
	o := orm.NewOrm()
	v := OrderInfo{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&OrderInfo{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
