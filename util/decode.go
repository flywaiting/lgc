package util

import (
	"errors"
	"reflect"
	"strconv"
)

func Decode(in, out interface{}) error {
	// 参考 mapstructure
	val := reflect.ValueOf(out)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return errors.New("无效输出")
	}
	val = val.Elem()
	if !val.CanAddr() {
		return errors.New("不可寻址输出，无法进行反射")
	}

	// 输入校验
	if in == nil {
		return nil
	}
	inVal := reflect.ValueOf(in)
	if (inVal.Kind() == reflect.Ptr && inVal.IsNil()) || !inVal.IsValid() {
		return nil
	}

	return decode(in, val)
}

func decode(data interface{}, val reflect.Value) error {
	// data 直接传递 reflect.Value，无法正确获取 Kind, 导致解析失败

	switch val.Kind() {
	case reflect.Int:
		return decodeInt(data, val)
	case reflect.String:
		return decodeStr(data, val)
	case reflect.Slice:
		return decodeSlice(data, val)
	case reflect.Struct:
		return decodeStruct(data, val)
	}
	return nil
}

func decodeStruct(data interface{}, val reflect.Value) (err error) {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	if dataVal.Kind() != reflect.Map {
		return
	}
	if kind := dataVal.Type().Key().Kind(); kind != reflect.String && kind != reflect.Interface {
		return errors.New("目前仅支持 map[string]interface{}")
	}

	valType := val.Type()
	for i := 0; i < valType.NumField(); i++ {
		field := val.Field(i)
		if !field.CanSet() {
			continue
		}

		// 获取字段对应内容
		fieldType := valType.Field(i)
		fieldName := fieldType.Name
		if tag := fieldType.Tag.Get("yaml"); tag != "" {
			fieldName = tag
		}
		mapVal := dataVal.MapIndex(reflect.ValueOf(fieldName))
		if !mapVal.IsValid() {
			continue
		}

		if err = decode(mapVal.Interface(), field); err != nil {
			break
		}
	}
	return
}

func decodeSlice(data interface{}, val reflect.Value) (err error) {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	if dataVal.IsNil() {
		return
	}
	dataKind := dataVal.Kind()
	if dataKind != reflect.Slice && dataKind != reflect.Array {
		return
	}

	valElemType := val.Type().Elem()
	valSlice := val
	// 分配内存
	size := dataVal.Len()
	if valSlice.IsNil() {
		valSlice = reflect.MakeSlice(reflect.SliceOf(valElemType), size, size)
	} else {
		for valSlice.Len() < size {
			valSlice = reflect.Append(valSlice, reflect.Zero(valElemType))
		}
	}

	for i := 0; i < size; i++ {
		decode(dataVal.Index(i).Interface(), valSlice.Index(i))
	}
	val.Set(valSlice)

	return
}
func decodeStr(data interface{}, val reflect.Value) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	switch dataVal.Kind() {
	case reflect.String:
		val.SetString(dataVal.String())
	case reflect.Int:
		val.SetString(strconv.FormatInt(dataVal.Int(), 10))
	}
	return nil
}

func decodeInt(data interface{}, val reflect.Value) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	if dataVal.Kind() == reflect.Int {
		val.SetInt(dataVal.Int())
	}
	return nil
}
