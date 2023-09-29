package utils

import (
	"errors"
	"reflect"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/8 17:22
  @describe :
*/

// IsNil 检测是否是真nil值
func IsNil(v interface{}) bool {
	if v != nil {
		// 如果目标不是xx类型,则返回
		typ := reflect.TypeOf(v)
		if typ.Kind() != reflect.Pointer {
			return false
		}

		value := reflect.ValueOf(v)
		if !value.IsNil() {
			return false
		}
	}
	return true
}

type CompareFunc func(left, right interface{}) bool

// LT 小于
func LT(left, right interface{}) bool {
	switch v := left.(type) {
	case int:
		return v < right.(int)
	case int8:
		return v < right.(int8)
	case int16:
		return v < right.(int16)
	case int32:
		return v < right.(int32)
	case int64:
		return v < right.(int64)
	case uint:
		return v < right.(uint)
	case uint8:
		return v < right.(uint8)
	case uint16:
		return v < right.(uint16)
	case uint32:
		return v < right.(uint32)
	case uint64:
		return v < right.(uint64)
	case string:
		return v < right.(string)
	case float32:
		return v < right.(float32)
	case float64:
		return v < right.(float64)
	default:
		panic(errors.New("not compare type value"))
	}
}

// LTE 小于等于
func LTE(left, right interface{}) bool {
	switch v := left.(type) {
	case int:
		return v <= right.(int)
	case int8:
		return v <= right.(int8)
	case int16:
		return v <= right.(int16)
	case int32:
		return v <= right.(int32)
	case int64:
		return v <= right.(int64)
	case uint:
		return v <= right.(uint)
	case uint8:
		return v <= right.(uint8)
	case uint16:
		return v <= right.(uint16)
	case uint32:
		return v <= right.(uint32)
	case uint64:
		return v <= right.(uint64)
	case string:
		return v <= right.(string)
	case float32:
		return v <= right.(float32)
	case float64:
		return v <= right.(float64)
	default:
		panic(errors.New("not compare type value"))
	}
}

// GT 大于
func GT(left, right interface{}) bool {
	switch v := left.(type) {
	case int:
		return v > right.(int)
	case int8:
		return v > right.(int8)
	case int16:
		return v > right.(int16)
	case int32:
		return v > right.(int32)
	case int64:
		return v > right.(int64)
	case uint:
		return v > right.(uint)
	case uint8:
		return v > right.(uint8)
	case uint16:
		return v > right.(uint16)
	case uint32:
		return v > right.(uint32)
	case uint64:
		return v > right.(uint64)
	case string:
		return v > right.(string)
	case float32:
		return v > right.(float32)
	case float64:
		return v > right.(float64)
	default:
		panic(errors.New("not compare type value"))
	}
}

// GTE 大于等于
func GTE(left, right interface{}) bool {
	switch v := left.(type) {
	case int:
		return v >= right.(int)
	case int8:
		return v >= right.(int8)
	case int16:
		return v >= right.(int16)
	case int32:
		return v >= right.(int32)
	case int64:
		return v >= right.(int64)
	case uint:
		return v >= right.(uint)
	case uint8:
		return v >= right.(uint8)
	case uint16:
		return v >= right.(uint16)
	case uint32:
		return v >= right.(uint32)
	case uint64:
		return v >= right.(uint64)
	case string:
		return v >= right.(string)
	case float32:
		return v >= right.(float32)
	case float64:
		return v >= right.(float64)
	default:
		panic(errors.New("not compare type value"))
	}
}

// EQ 等于
func EQ(left, right interface{}) bool {
	switch v := left.(type) {
	case int:
		return v == right.(int)
	case int8:
		return v == right.(int8)
	case int16:
		return v == right.(int16)
	case int32:
		return v == right.(int32)
	case int64:
		return v == right.(int64)
	case uint:
		return v == right.(uint)
	case uint8:
		return v == right.(uint8)
	case uint16:
		return v == right.(uint16)
	case uint32:
		return v == right.(uint32)
	case uint64:
		return v == right.(uint64)
	case string:
		return v == right.(string)
	case float32:
		return v == right.(float32)
	case float64:
		return v == right.(float64)
	default:
		panic(errors.New("not compare type value"))
	}
}
