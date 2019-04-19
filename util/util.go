package util

import (
	"fmt"
	"time"
	"errors"
	"strconv"
	"strings"
	"reflect"

	"anole/common"
)

// Convert2Object 将params 转化为idl
func Convert2Object(params common.Params, obj interface{}) error {
	for _, param := range params {
		err := fillField(obj, param.Key, param.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

// fillField ...
func fillField(obj interface{}, name string, value interface{}) error {

	var FormatName = func(name string) string {
		var field string

		var UcFirst = func(str string) string {
			if str != "" && str[0] != ' ' {
				return strings.ToUpper(string(str[0])) + str[1:]
			}
			return str
		}

		if !strings.Contains(name, "_") {
			return UcFirst(name)
		}
		for _, name := range strings.Split(name, "_") {
			field += UcFirst(name)
		}
		return field
	}

	field := FormatName(name)
	structValue := reflect.ValueOf(obj).Elem()         //结构体属性值
	structFieldValue := structValue.FieldByName(field) //结构体单个属性值

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type() //结构体的类型
	val := reflect.ValueOf(value)              //map值的反射值

	var err error
	if structFieldType != val.Type() {
		val, err = TypeConversion(fmt.Sprintf("%v", value), structFieldValue.Type().Name()) //类型转换
		if err != nil {
			return err
		}
	}

	structFieldValue.Set(val)
	return nil
}

// TypeConversion 类型转换
func TypeConversion(value string, ntype string) (reflect.Value, error) {
	if ntype == "string" {
		return reflect.ValueOf(value), nil
	} else if ntype == "time.Time" {
		t, err := time.ParseInLocation(common.TimeFormat, value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "Time" {
		t, err := time.ParseInLocation(common.TimeFormat, value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "int" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	} else if ntype == "int8" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	} else if ntype == "int32" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int64(i)), err
	} else if ntype == "int64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "float32" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	} else if ntype == "float64" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	}

	// TODO 增加其他一些类型的转换

	return reflect.ValueOf(value), errors.New("未知的类型：" + ntype)
}
