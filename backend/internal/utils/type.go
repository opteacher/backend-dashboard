package utils

import (
	"reflect"
)

func GetAllFields(typ reflect.Type) []reflect.StructField {
	var fields []reflect.StructField
	for j := 0; j < typ.NumField(); j++ {
		fields = append(fields, typ.Field(j))
	}
	return fields
}

func Obj2Map(obj interface{}) map[string]interface{} {
	mp := make(map[string]interface{})
	typ := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)
	for i := 0; i < typ.NumField(); i++ {
		fname := typ.Field(i).Name
		mp[fname] = val.FieldByName(fname).Interface()
	}
	return mp
}
