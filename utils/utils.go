package utils

import "reflect"

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/8 17:22
  @describe :
*/

// IsNil 检测是否是真nil值
func IsNil(v any) bool {
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
