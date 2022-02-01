package models

import (
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"mi-beego/pkg/lib"
	"reflect"
	"strings"
)

type UserAddress struct {
	Id           int64    `orm:"column(id);auto" description:"主键" json:"id,string"`
	AddressId    int64    `orm:"column(address_id)" description:"地址id" json:"addressId,string"`
	CreateTime   lib.Time `orm:"column(create_time);type(datetime);auto_now" description:"创建时间" json:"-"`
	UpdateTime   lib.Time `orm:"column(update_time);type(datetime);auto_now" description:"更新时间" json:"-"`
	CreateUser   int64    `orm:"column(create_user)" description:"创建人" json:"-"`
	UpdateUser   int64    `orm:"column(update_user)" description:"更新人" json:"-"`
	DeleteFlag   int8     `orm:"column(delete_flag)" description:"删除标志" json:"-"`
	DefaultFlag  int8     `orm:"column(default_flag);null" description:"默认标志" json:"defaultFlag"`
	ReceiveName  string   `orm:"column(receive_name);size(60)" description:"收货人" json:"receiveName"`
	ReceivePhone string   `orm:"column(receive_phone);size(20)" description:"联系号码" json:"receivePhone"`
	Province     string   `orm:"column(province);size(20)" description:"省份" json:"province"`
	ProvinceCode string   `orm:"column(province_code);size(10)" description:"省份编码" json:"provinceCode"`
	City         string   `orm:"column(city);size(20)" description:"城市" json:"city"`
	CityCode     string   `orm:"column(city_code);size(10)" description:"城市编码" json:"cityCode"`
	Area         string   `orm:"column(area);size(20)" description:"区" json:"area"`
	AreaCode     string   `orm:"column(area_code);size(10)" description:"区编码" json:"areaCode"`
	Street       string   `orm:"column(street);size(100);null" description:"详细地址" json:"street"`
	PostalCode   string   `orm:"column(postal_code);size(10);null" description:"邮编" json:"postalCode"`
	AddressLabel int8     `orm:"column(address_label);null" description:"地址标签" json:"addressLabel"`
}

func (t *UserAddress) TableName() string {
	return "user_address"
}

func init() {
	orm.RegisterModel(new(UserAddress))
}

// AddUserAddress insert a new UserAddress into database and returns
// last inserted Id on success.
func AddUserAddress(m *UserAddress) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetUserAddressById retrieves UserAddress by Id. Returns error if
// Id doesn't exist
func GetUserAddressById(id int64) (v *UserAddress, err error) {
	o := orm.NewOrm()
	v = &UserAddress{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetUserAddressByAddrId retrieves UserAddress by Id. Returns error if
// Id doesn't exist
func GetUserAddressByAddrId(addrId int64) (v *UserAddress, err error) {
	o := orm.NewOrm()
	v = &UserAddress{AddressId: addrId}
	if err = o.Read(v, "AddressId"); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllUserAddress retrieves all UserAddress matches certain condition. Returns empty list if
// no records exist
func GetAllUserAddress(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, count int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(UserAddress))
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
					return nil, 0, errors.New("Error: Invalid order. Must be either [asc|desc]")
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
					return nil, 0, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, 0, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, 0, errors.New("Error: unused 'order' fields")
		}
	}
	count, _ = qs.Count()
	var l []UserAddress
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
		return ml, count, nil
	}
	return nil, 0, err
}

// UpdateUserAddress updates UserAddress by Id and returns error if
// the record to be updated doesn't exist
func UpdateUserAddressById(m *UserAddress) (err error) {
	o := orm.NewOrm()
	v := UserAddress{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		//选择性更新字段内容
		if num, err = o.Update(m, "UpdateTime", "UpdateUser", "DeleteFlag", "DefaultFlag", "ReceiveName", "ReceivePhone", "Province", "ProvinceCode", "City", "CityCode", "Area", "AreaCode", "Street", "PostalCode", "AddressLabel"); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUserAddress deletes UserAddress by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUserAddress(id int64) (err error) {
	o := orm.NewOrm()
	v := UserAddress{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&UserAddress{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
