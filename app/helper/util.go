package helper

import (
	"errors"
	"reflect"
	"strings"
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

func UrlHttpToHttps(url string) string {
	if strings.Contains(url, "http://") {
		newUrl := strings.ReplaceAll(url, "http://", "https://")
		return newUrl
	}
	return url
}

func GetStructFieldsAndValuesExcept(s interface{}, except []string) ([]string, []string) {
	var fields []string
	var values []string
	tags := GetStructTag(s)
	for _, v := range tags {
		if v == "-" || CheckStringElemInArray(except, v) {
			continue
		}
		fields = append(fields, v)
		values = append(values, ":"+v)
	}
	return fields, values
}

func GetStructTag(v interface{}) []string {
	var keys []string
	keysSlice := keys[:]
	ty := reflect.TypeOf(v)
	for i := 0; i < ty.NumField(); i++ {
		key := ty.Field(i).Tag.Get("db")
		if len(key) == 0 {
			continue
		}
		keysSlice = append(keysSlice, key)
	}
	return keysSlice
}

func CheckStringElemInArray(array []string, value string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}
