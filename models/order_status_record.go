package models

import (
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"mi-beego/pkg/lib"
	"reflect"
	"strings"
)

type OrderStatusRecord struct {
	Id            int      `orm:"column(id);auto" description:"主键"`
	CreateTime    lib.Time `orm:"column(create_time);type(datetime)" description:"创建时间"`
	OrderNo       string   `orm:"column(order_no);size(60)" description:"订单编号"`
	OrderDetailNo string   `orm:"column(order_detail_no);size(60)" description:"订单明细编号"`
	ProductId     string   `orm:"column(product_id)" description:"商品id"`
	ProductName   string   `orm:"column(product_name);size(60);null" description:"商品名称"`
	Status        int8     `orm:"column(status)" description:"订单状态"`
	StatusDesc    string   `orm:"column(status_desc);size(60);null" description:"状态描述"`
}

func (t *OrderStatusRecord) TableName() string {
	return "order_status_record"
}

func init() {
	orm.RegisterModel(new(OrderStatusRecord))
}

// AddOrderStatusRecord insert a new OrderStatusRecord into database and returns
// last inserted Id on success.
func AddOrderStatusRecord(m *OrderStatusRecord) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetOrderStatusRecordById retrieves OrderStatusRecord by Id. Returns error if
// Id doesn't exist
func GetOrderStatusRecordById(id int) (v *OrderStatusRecord, err error) {
	o := orm.NewOrm()
	v = &OrderStatusRecord{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllOrderStatusRecord retrieves all OrderStatusRecord matches certain condition. Returns empty list if
// no records exist
func GetAllOrderStatusRecord(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(OrderStatusRecord))
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
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
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
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []OrderStatusRecord
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
		return ml, nil
	}
	return nil, err
}

// UpdateOrderStatusRecord updates OrderStatusRecord by Id and returns error if
// the record to be updated doesn't exist
func UpdateOrderStatusRecordById(m *OrderStatusRecord) (err error) {
	o := orm.NewOrm()
	v := OrderStatusRecord{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteOrderStatusRecord deletes OrderStatusRecord by Id and returns error if
// the record to be deleted doesn't exist
func DeleteOrderStatusRecord(id int) (err error) {
	o := orm.NewOrm()
	v := OrderStatusRecord{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&OrderStatusRecord{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
