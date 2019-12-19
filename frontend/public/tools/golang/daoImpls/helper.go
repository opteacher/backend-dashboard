package utils

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"net/http"
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

func ToObj(mp map[string]interface{}, typ reflect.Type) (interface{}, error) {
	if bytes, err := json.Marshal(mp); err != nil {
		return nil, err
	} else {
		return UnmarshalJSON(bytes, typ)
	}
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Clone(in, out interface{}) error {
	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	dec := gob.NewDecoder(buff)
	if err := enc.Encode(in); err != nil {
		return err
	}
	if err := dec.Decode(out); err != nil {
		return err
	}
	return nil
}

func HttpGetJsonMap(url string) (map[string]interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	retMap := make(map[string]interface{})
	if err := json.Unmarshal(body, &retMap); err != nil {
		return nil, err
	}
	return retMap, nil
}