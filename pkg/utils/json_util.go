package utils

import (
	"encoding/json"
	"fmt"
)

//利用json  实现不同对象间的属性复制
func CopyToObject(objSource interface{}, objTarget interface{}) error {
	toJson := ObjectToJson(objSource)
	return JsonToObject(toJson, objTarget)
}

//json  转对象
func JsonToObject(data string, obj interface{}) error {
	return json.Unmarshal([]byte(data), obj)
}

//对象转json 字符串
func ObjectToJson(data interface{}) string {
	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("json 转换错误,", err)
		return ""
	} else {
		return string(bytes)
	}
}
