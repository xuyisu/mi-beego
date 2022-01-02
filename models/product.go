package models

import (
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"mi-beego/pkg/lib"
	"reflect"
	"strings"
)

type Product struct {
	Id         int      `orm:"column(id);auto" description:"主键" json:"id"`
	CreateTime lib.Time `orm:"column(create_time);type(datetime);auto_now" description:"创建时间" json:"createTime"`
	UpdateTime lib.Time `orm:"column(update_time);type(datetime);auto_now" description:"更新时间" json:"-"`
	CreateUser string   `orm:"column(create_user)" description:"创建人" json:"-"`
	UpdateUser string   `orm:"column(update_user)" description:"更新人" json:"-"`
	DeleteFlag int8     `orm:"column(delete_flag)" description:"删除标志" json:"-"`
	ProductId  string   `orm:"column(product_id);null" description:"商品id" json:"productId"`
	CategoryId string   `orm:"column(category_id);null" description:"品类id" json:"categoryId"`
	Name       string   `orm:"column(name);size(60);null" description:"商品名称" json:"name"`
	SubTitle   string   `orm:"column(sub_title);size(100);null" description:"简要描述" json:"subTitle"`
	MainImage  string   `orm:"column(main_image);size(100);null" description:"商品图片地址" json:"mainImage"`
	SubImages  string   `orm:"column(sub_images);size(100);null" description:"子图片列表" json:"subImages"`
	ActivityId string   `orm:"column(activity_id);null" description:"活动id" json:"activityId"`
	Status     int8     `orm:"column(status)" description:"商品状态" json:"status"`
	Price      float64  `orm:"column(price);digits(20);decimals(2)" description:"商品单价" json:"price"`
	Stock      int      `orm:"column(stock)" description:"库存数" json:"stock"`
}

func (t *Product) TableName() string {
	return "product"
}

func init() {
	orm.RegisterModel(new(Product))
}

// AddProduct insert a new Product into database and returns
// last inserted Id on success.
func AddProduct(m *Product) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetProductById retrieves Product by Id. Returns error if
// Id doesn't exist
func GetProductById(id int) (v *Product, err error) {
	o := orm.NewOrm()
	v = &Product{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetProductById retrieves Product by Id. Returns error if
// Id doesn't exist
func GetProductByProductId(productId string) (v *Product, err error) {
	o := orm.NewOrm()
	v = &Product{ProductId: productId, Status: 1, DeleteFlag: 0}
	if err = o.Read(v, "productId", "status", "deleteFlag"); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllProduct retrieves all Product matches certain condition. Returns empty list if
// no records exist
func GetAllProduct(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, count int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Product))
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
	var l []Product
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

// UpdateProduct updates Product by Id and returns error if
// the record to be updated doesn't exist
func UpdateProductById(m *Product) (err error) {
	o := orm.NewOrm()
	v := Product{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteProduct deletes Product by Id and returns error if
// the record to be deleted doesn't exist
func DeleteProduct(id int) (err error) {
	o := orm.NewOrm()
	v := Product{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Product{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
