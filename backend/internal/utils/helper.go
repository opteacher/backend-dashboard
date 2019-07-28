package utils

import (
	"encoding/json"
	"reflect"
)

func UnmarshalJSON(bdata []byte, outTyp reflect.Type) (interface{}, error) {
	ret := reflect.New(outTyp).Interface()
	err := json.Unmarshal(bdata, ret)
	return ret, err
}

func ToMap(entry interface{}) (map[string]interface{}, error) {
	if bytes, err := json.Marshal(entry); err != nil {
		return nil, err
	} else if mp, err := UnmarshalJSON(bytes, reflect.TypeOf((*map[string]interface{})(nil)).Elem()); err != nil {
		return nil, err
	} else {
		return *(mp.(*map[string]interface{})), nil
	}
}