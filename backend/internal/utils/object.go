package utils

import (
	"fmt"
	"reflect"
	"encoding/json"
)

// TODO: 逻辑巨乱，找时间加注释！！！！
func FillWithMap(obj interface{}, mp map[string]interface{}) interface{} {
	typ := reflect.TypeOf(obj)
	ele := reflect.ValueOf(obj)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		ele = ele.Elem()
	}
	fmt.Println(typ)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		filler := ele.Field(i)
		mname := CamelToPascal(field.Name)
		if _, exs := mp[mname]; !exs {
			continue
		}
		mfield := mp[mname]
		mvalue := reflect.ValueOf(mfield)
		fkind := field.Type.Kind()
		if fkind == reflect.Slice || fkind == reflect.Array {
			filler.Set(reflect.MakeSlice(field.Type, mvalue.Len(), mvalue.Len()))
			etype := field.Type.Elem()
			if etype.Kind() == reflect.Ptr {
				etype = etype.Elem()
			}
			for j := 0; j < mvalue.Len(); j++ {
				item := mvalue.Index(j)
				fitem := filler.Index(j)
				subMap := item.Interface().(map[string]interface{})
				fvalue, exs := subMap[ToSingular(mname)]
				if fitem.Type().Kind() != reflect.Ptr && exs {
					fitem.Set(reflect.ValueOf(fvalue))
				} else {
					fitem.Set(reflect.New(etype))
					FillWithMap(fitem.Interface(), subMap)
				}
			}
		} else {
			mvalue = mvalue.Convert(field.Type)
			filler.Set(mvalue)
		}
	}
	return obj
}

// NOTE: 返回的是src类型的指针
func Copy(src interface{}) (interface{}, error) {
	if reflect.TypeOf(src).Kind() == reflect.Ptr {
		src = reflect.ValueOf(src).Elem().Interface()
	}
	if bytes, err := json.Marshal(src); err != nil {
		return nil, err
	} else if tgt, err := UnmarshalJSON(bytes, reflect.TypeOf(src)); err != nil {
		return nil, err
	} else {
		return tgt, nil
	}
}