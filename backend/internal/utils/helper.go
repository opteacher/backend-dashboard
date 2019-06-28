package utils

import (
	"reflect"
	"encoding/json"
)

func GetTypeOfPtr(ptr interface{}) reflect.Type {
	return reflect.TypeOf(ptr).Elem()
}

func WrapJsonUnmarshal(bdata []byte, outTyp reflect.Type) (interface{}, error) {
	ret := reflect.New(outTyp)
	err := json.Unmarshal(bdata, ret)
	return ret, err
}