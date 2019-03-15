// Package queryize @author K·J Create at 2019-03-04 14:51
package queryize

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"unsafe"
)

// Config config
type Config struct {
	Joiner    string
	Separator string
}

// API api define
type API interface {
	Marshal(v interface{}) (string, error)
	Unmarshal(data map[string][]string, v interface{}) error
}

var ConfigDefault = Config{
	Joiner:    "=",
	Separator: "&",
}.Froze()

const key = "query"

func (c Config) Froze() API {
	return &c
}

func (c Config) Marshal(v interface{}) (string, error) {
	t := reflect.TypeOf(v)
	param := ""
	switch t.Kind() {
	case reflect.Ptr:
		value := reflect.ValueOf(v).Elem()
		for i := 0; i < t.Elem().NumField(); i++ {
			tag := t.Elem().Field(i).Tag.Get(key)
			if "" != tag {
				v := value.Field(i)
				val, err := getValue(v)
				if nil != err {
					return param, err
				}
				if "" != val {
					param += tag + c.Joiner + val + c.Separator
				}
			}
		}
	case reflect.Struct:
		value := reflect.ValueOf(v)
		for i := 0; i < t.NumField(); i++ {
			tag := t.Field(i).Tag.Get(key)
			if "" != tag {
				v := value.Field(i)
				val, err := getValue(v)
				if nil != err {
					return param, err
				}
				if "" != val {
					param += tag + c.Joiner + val + c.Separator
				}
			}
		}
	case reflect.Map:
	case reflect.Slice:
	default:

	}
	if len(param) > 0 {
		return param[:len(param)-1], nil
	}
	return param, nil
}

func (c Config) Unmarshal(data map[string][]string, v interface{}) error {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Ptr {
		return errors.New("queryize: interface must be a pointer")
	}
	value := reflect.ValueOf(v).Elem()
	for i := 0; i < t.Elem().NumField(); i++ {
		tag := t.Elem().Field(i).Tag.Get(key)
		if "" != tag {
			if param, ok := data[tag]; ok {
				if len(param) > 0 && len(param[0]) > 0 {
					err := setValue(value.Field(i), param)
					if nil != err {
						return err
					}
				}
			}
		}
	}
	return nil
}

// setValue set value
func setValue(v reflect.Value, param []string) error {
	if nil == param || len(param) == 0 {
		return errors.New("param invalid")
	}
	if v.CanSet() {
		switch v.Kind() {
		case reflect.Ptr:
			switch v.Type().String() {
			case "*int":
				val, err := strconv.Atoi(param[0])
				if nil != err {
					return err
				}
				p := new(int)
				*p = val
				v.Set(reflect.ValueOf(p))
			case "*int64":
				val, err := strconv.ParseInt(param[0], 10, 64)
				if nil != err {
					return err
				}
				p := new(int64)
				*p = val
				v.Set(reflect.ValueOf(p))
			case "*string":
				p := new(string)
				*p = param[0]
				v.Set(reflect.ValueOf(p))
			default:
				return errors.New("not support type")
			}
		case reflect.String:
			v.SetString(param[0])
		case reflect.Int:
			val, err := strconv.ParseInt(param[0], 10, 32)
			if nil != err {
				return err
			}
			v.SetInt(val)
		case reflect.Int64:
			val, err := strconv.ParseInt(param[0], 10, 64)
			if nil != err {
				return err
			}
			v.SetInt(val)
		default:
			return errors.New("not support type")
		}
	} else {
		if !v.CanAddr() {
			return errors.New("can not get addr")
		}
		addr := v.Addr()
		ptr := unsafe.Pointer(addr.Pointer())
		switch v.Type().String() {
		case "int":
			fallthrough
		case "*int":
			val, err := strconv.Atoi(param[0])
			if nil != err {
				return err
			}
			*(*int)(ptr) = val
		case "int64":
			fallthrough
		case "*int64":
			val, err := strconv.ParseInt(param[0], 10, 64)
			if nil != err {
				return err
			}
			*(*int64)(ptr) = val
		case "*string":
			fallthrough
		case "string":
			*(*string)(ptr) = param[0]
		default:
			return errors.New("not support type")
		}
	}
	return nil
}

func getValue(v reflect.Value) (string, error) {
	if v.CanInterface() {
		switch v.Kind() {
		case reflect.Ptr:
			if !v.IsNil() {
				return fmt.Sprintf("%+v", v.Elem().Interface()), nil
			}
		default:
			return fmt.Sprintf("%v", v.Interface()), nil
		}
	} else {
		if v.CanAddr() {
			addr := v.Addr()
			ptr := unsafe.Pointer(addr.Pointer())
			switch v.Type().String() {
			case "int":
				fallthrough
			case "*int":
				fallthrough
			case "int64":
				fallthrough
			case "*int64":
				return fmt.Sprintf("%d", *(*int)(ptr)), nil
			case "*string":
				fallthrough
			case "string":
				return *(*string)(ptr), nil
			default:
				return "", errors.New("not support type")
			}
		}
	}
	return "", errors.New("not support type")
}
