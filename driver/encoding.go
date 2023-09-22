package driver

import (
	"encoding"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/13 11:42
  @describe :
*/

// marshalData 编码数据
func marshalData(data interface{}) (string, error) {
	// 先判断是否是指针类型
	val := make([]byte, 0)
	switch d := data.(type) {
	case nil:
	case string:
		val = []byte(d)
	case []byte:
		// 复制数据,不能直接设置
		val = make([]byte, len(d))
		copy(val, d)
	case int:
		val = strconv.AppendInt(val, int64(d), 10)
	case int8:
		val = strconv.AppendInt(val, int64(d), 10)
	case int16:
		val = strconv.AppendInt(val, int64(d), 10)
	case int32:
		val = strconv.AppendInt(val, int64(d), 10)
	case int64:
		val = strconv.AppendInt(val, d, 10)
	case uint:
		val = strconv.AppendUint(val, uint64(d), 10)
	case uint8:
		val = strconv.AppendUint(val, uint64(d), 10)
	case uint16:
		val = strconv.AppendUint(val, uint64(d), 10)
	case uint32:
		val = strconv.AppendUint(val, uint64(d), 10)
	case uint64:
		val = strconv.AppendUint(val, d, 10)
	case float32:
		val = strconv.AppendFloat(val, float64(d), 'f', -1, 64)
	case float64:
		val = strconv.AppendFloat(val, d, 'f', -1, 64)
	case bool:
		if d {
			val = strconv.AppendInt(val, 1, 10)
			break
		}
		val = strconv.AppendInt(val, 0, 10)
	case time.Time:
		val = d.AppendFormat(val, time.RFC3339Nano)
	case *time.Time:
		val = d.AppendFormat(val, time.RFC3339Nano)
	case time.Duration:
		val = strconv.AppendInt(val, d.Nanoseconds(), 10)
	case encoding.BinaryMarshaler:
		b, err := d.MarshalBinary()
		if err != nil {
			return "", err
		}
		val = b
	case net.IP:
		val = d
	default:
		return "", fmt.Errorf(
			"memory cache: can't marshal %T (implement encoding.BinaryMarshaler)", d)
	}
	return string(val), nil
}

func sliceArgs(args []interface{}) []interface{} {
	dst := make([]interface{}, 0, len(args))
	if len(args) == 0 {
		return nil
	}
	if len(args) == 1 {
		return sliceArg(dst, args[0])
	}

	for _, arg := range args {
		dst = sliceArg(dst, arg)
	}

	return dst
}

func sliceArg(dst []interface{}, arg interface{}) []interface{} {
	switch arg := arg.(type) {
	case []string:
		for _, s := range arg {
			dst = append(dst, s)
		}
		return dst
	case []interface{}:
		dst = append(dst, arg...)
		return dst
	case map[string]interface{}:
		for k, v := range arg {
			dst = append(dst, k, v)
		}
		return dst
	case map[string]string:
		for k, v := range arg {
			dst = append(dst, k, v)
		}
		return dst
	case time.Time, time.Duration, encoding.BinaryMarshaler, net.IP:
		return append(dst, arg)
	default:
		// scan struct field
		v := reflect.ValueOf(arg)
		if v.Type().Kind() == reflect.Ptr {
			if v.IsNil() {
				// error: arg is not a valid object
				return dst
			}
			v = v.Elem()
		}

		if v.Type().Kind() == reflect.Struct {
			return appendStructField(dst, v)
		}

		return append(dst, arg)
	}
}

func getInTag(tag reflect.StructTag, names ...string) string {
	for _, name := range names {
		val := tag.Get(name)
		if val != "" && val != "-" {
			return val
		}
	}
	return ""
}

// appendStructField appends the field and value held by the structure v to dst, and returns the appended dst.
func appendStructField(dst []interface{}, v reflect.Value) []interface{} {
	typ := v.Type()
	for i := 0; i < typ.NumField(); i++ {
		tag := getInTag(typ.Field(i).Tag, "redis", "memory")
		if tag == "" {
			continue
		}

		name, opt, _ := strings.Cut(tag, ",")
		if name == "" {
			continue
		}

		field := v.Field(i)

		// miss field
		if omitEmpty(opt) && isEmptyValue(field) {
			continue
		}

		if field.CanInterface() {
			dst = append(dst, name, field.Interface())
		}
	}

	return dst
}

func omitEmpty(opt string) bool {
	for opt != "" {
		var name string
		name, opt, _ = strings.Cut(opt, ",")
		if name == "omitempty" {
			return true
		}
	}
	return false
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Pointer:
		return v.IsNil()
	}
	return false
}
