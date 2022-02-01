package models

import (
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"mi-beego/pkg/lib"
	"reflect"
	"strings"
)

type OrderDetail struct {
	Id                int64    `orm:"column(id);auto" description:"主键" json:"id,string"`
	CreateTime        lib.Time `orm:"column(create_time);type(datetime);auto_now_add" description:"创建时间" json:"createTime"`
	UpdateTime        lib.Time `orm:"column(update_time);type(datetime);auto_now" description:"更新时间" json:"-"`
	CreateUser        int64    `orm:"column(create_user)" description:"创建人" json:"-"`
	UpdateUser        int64    `orm:"column(update_user)" description:"更新人" json:"-"`
	DeleteFlag        int8     `orm:"column(delete_flag)" description:"删除标志" json:"-"`
	OrderNo           string   `orm:"column(order_no);size(60)" description:"订单编号" json:"orderNo"`
	OrderDetailNo     string   `orm:"column(order_detail_no);size(60)" description:"订单明细编号" json:"orderDetailNo"`
	ActivityId        int64    `orm:"column(activity_id);null" description:"活动id" json:"activityId,string"`
	ActivityName      string   `orm:"column(activity_name);size(50);null" description:"活动名称" json:"activityName"`
	ActivityMainImage string   `orm:"column(activity_main_image);size(100);null" description:"活动图片地址" json:"activityMainImage"`
	ProductId         int64    `orm:"column(product_id)" description:"商品id" json:"productId,string"`
	ProductName       string   `orm:"column(product_name);size(50)" description:"商品名称" json:"productName"`
	ProductMainImage  string   `orm:"column(product_main_image);size(100)" description:"商品图片地址" json:"productMainImage"`
	CurrentUnitPrice  float64  `orm:"column(current_unit_price);null;digits(20);decimals(2)" description:"单价" json:"currentUnitPrice"`
	Quantity          int      `orm:"column(quantity);null" description:"数量" json:"quantity"`
	TotalPrice        float64  `orm:"column(total_price);null;digits(20);decimals(2)" description:"总价" json:"totalPrice"`
	UserId            int64    `orm:"column(user_id)" description:"购买人id" json:"userId,string"`
	Status            int8     `orm:"column(status)" description:"订单状态" json:"status"`
	StatusDesc        string   `orm:"column(status_desc);size(20);null" description:"状态描述" json:"statusDesc"`
	CancelTime        lib.Time `orm:"column(cancel_time);type(datetime);null" description:"取消时间" json:"cancelTime"`
	CancelReason      int      `orm:"column(cancel_reason);null" description:"取消原因" json:"cancelReason"`
	SendTime          lib.Time `orm:"column(send_time);type(datetime);null" description:"发货时间" json:"sendTime"`
	ReceiveTime       lib.Time `orm:"column(receive_time);type(datetime);null" description:"签收时间" json:"receiveTime"`
}

func (t *OrderDetail) TableName() string {
	return "order_detail"
}

func init() {
	orm.RegisterModel(new(OrderDetail))
}

// AddOrderDetail insert a new OrderDetail into database and returns
// last inserted Id on success.
func AddOrderDetail(m *OrderDetail) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetOrderDetailById retrieves OrderDetail by Id. Returns error if
// Id doesn't exist
func GetOrderDetailById(id int64) (v *OrderDetail, err error) {
	o := orm.NewOrm()
	v = &OrderDetail{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllOrderDetail retrieves all OrderDetail matches certain condition. Returns empty list if
// no records exist
func GetAllOrderDetail(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, orderDetailList []OrderDetail, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(OrderDetail))
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

	var l []OrderDetail
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

// UpdateOrderDetail updates OrderDetail by Id and returns error if
// the record to be updated doesn't exist
func UpdateOrderDetailById(m *OrderDetail) (err error) {
	o := orm.NewOrm()
	v := OrderDetail{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteOrderDetail deletes OrderDetail by Id and returns error if
// the record to be deleted doesn't exist
func DeleteOrderDetail(id int64) (err error) {
	o := orm.NewOrm()
	v := OrderDetail{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&OrderDetail{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
