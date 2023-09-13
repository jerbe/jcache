package driver

import (
	"encoding"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"time"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/13 11:42
  @describe :
*/

// marshalData 编码数据
func marshalData(data any) (string, error) {
	// 先判断是否是指针类型
	if data != nil {
		if _, ok := data.(encoding.BinaryMarshaler); !ok && reflect.TypeOf(data).Kind() == reflect.Ptr {
			value := reflect.Indirect(reflect.ValueOf(data))
			data = value.Interface()
		}
	}

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
