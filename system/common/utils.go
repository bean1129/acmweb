package common

import (
	"reflect"

	"acmweb/system/text"
)

type ZUtils struct {
}

func NewUtils() *ZUtils {
	return &ZUtils{}
}

// InArray 判断元素是否在数组中
func (c *ZUtils) InArray(value string, array []interface{}) bool {
	for _, v := range array {
		if text.Conv.String(v) == value {
			return true
		}
	}
	return false
}

// Reverse 数组反转
func (c *ZUtils) Reverse(arr *[]string) {
	length := len(*arr)
	var temp string
	for i := 0; i < length/2; i++ {
		temp = (*arr)[i]
		(*arr)[i] = (*arr)[length-1-i]
		(*arr)[length-1-i] = temp
	}
}

func (c *ZUtils) InStringArray(value string, array []string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

func (c *ZUtils) StructCopy(from, to interface{}) {
	fromValue := reflect.ValueOf(from)
	toValue := reflect.ValueOf(to)
	// 必须是指针类型
	if fromValue.Kind() != reflect.Ptr || toValue.Kind() != reflect.Ptr {
		return
	}
	if fromValue.IsNil() || toValue.IsNil() {
		return
	}
	// 获取到来源数据
	fromElem := fromValue.Elem()
	// 需要的数据
	toElem := toValue.Elem()
	for i := 0; i < toElem.NumField(); i++ {
		toField := toElem.Type().Field(i)
		// 看看来源的结构体中是否有这个属性
		fromFieldName, ok := fromElem.Type().FieldByName(toField.Name)
		// 存在相同的属性名称并且类型一致
		// todo 可以根据需要判断是否是空值
		if ok && fromFieldName.Type == toField.Type {
			toElem.Field(i).Set(fromElem.FieldByName(toField.Name))
		}
	}
}
