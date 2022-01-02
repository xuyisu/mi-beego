package models

import (
	"errors"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"mi-beego/pkg/lib"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Activity struct {
	Id         int      `orm:"column(id);auto" description:"主键"`
	CreateTime lib.Time `orm:"column(create_time);type(datetime);auto_now" description:"创建时间"`
	UpdateTime lib.Time `orm:"column(update_time);type(datetime);auto_now" description:"更新时间"`
	CreateUser int64    `orm:"column(create_user)" description:"创建人"`
	UpdateUser int64    `orm:"column(update_user)" description:"更新人"`
	DeleteFlag int8     `orm:"column(delete_flag)" description:"删除标志"`
	ActivityId string   `orm:"column(activity_id)" description:"活动id"`
	Name       string   `orm:"column(name);size(60);null" description:"活动名称"`
	Status     int8     `orm:"column(status)" description:"活动状态"`
	MainImage  string   `orm:"column(main_image);size(100);null" description:"活动图片地址"`
	StartTime  lib.Time `orm:"column(start_time);type(datetime);null" description:"开始时间"`
	EndTime    lib.Time `orm:"column(end_time);type(datetime);null" description:"结束时间"`
}

func (t *Activity) TableName() string {
	return "activity"
}

func init() {
	orm.RegisterModel(new(Activity))
}

// AddActivity insert a new Activity into database and returns
// last inserted Id on success.
func AddActivity(m *Activity) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetActivityById retrieves Activity by Id. Returns error if
// Id doesn't exist
func GetActivityById(id int) (v *Activity, err error) {
	o := orm.NewOrm()
	v = &Activity{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetActivityByTime retrieves Activity by Id. Returns error if
// Id doesn't exist
func GetActivityByTime(date time.Time) (v *Activity, err error) {
	v = &Activity{}
	o := orm.NewOrm()
	qs := o.QueryTable(v)
	err = qs.Filter("start_time__lt", date).Filter("end_time__gt", date).Filter("delete_flag", 0).Filter("status", 1).Limit(1).One(v)
	if err != nil {
		logs.Error(err.Error())
		return nil, nil
	}
	return v, nil
}

// GetAllActivity retrieves all Activity matches certain condition. Returns empty list if
// no records exist
func GetAllActivity(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Activity))
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

	var l []Activity
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

// UpdateActivity updates Activity by Id and returns error if
// the record to be updated doesn't exist
func UpdateActivityById(m *Activity) (err error) {
	o := orm.NewOrm()
	v := Activity{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteActivity deletes Activity by Id and returns error if
// the record to be deleted doesn't exist
func DeleteActivity(id int) (err error) {
	o := orm.NewOrm()
	v := Activity{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Activity{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
