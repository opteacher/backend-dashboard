package utils

import (
	"reflect"
	"encoding/json"
)

func UnmarshalJSON(bdata []byte, outTyp reflect.Type) (interface{}, error) {
	ret := reflect.New(outTyp).Interface()
	err := json.Unmarshal(bdata, ret)
	return ret, err
}