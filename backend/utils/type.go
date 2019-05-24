package utils

import "reflect"

func GetAllFields(typ reflect.Type) []reflect.StructField {
	var fields []reflect.StructField
	for j := 0; j < typ.NumField(); j++ {
		fields = append(fields, typ.Field(j))
	}
	return fields
}
