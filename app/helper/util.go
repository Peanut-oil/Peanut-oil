package helper

import (
	"errors"
	"reflect"
)

// Struct2Map 转换struct为map
func Struct2Map(obj interface{}, tagName string) (map[string]interface{}, error) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	var data = make(map[string]interface{})
	if t.Kind() != reflect.Struct {
		return data, errors.New("type is not struct")
	}
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Tag.Get(tagName)
		if name == "-" {
			continue
		}
		if name == "" {
			name = t.Field(i).Name
		}
		data[name] = v.Field(i).Interface()
	}
	return data, nil
}
