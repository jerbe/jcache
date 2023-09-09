package memory

import (
	"context"
	"encoding"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"time"

	"github.com/jerbe/jcache/utils"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/8 15:58
  @describe :
*/

// expirable 可以用于过期的
type expirable interface {
	// IsExpire 判断是否已经过期了
	IsExpire() bool

	// SetExpire 设置过期时间间隔
	SetExpire(time.Duration)

	// SetExpireAt 设置时间
	SetExpireAt(*time.Time)
}

type expireValue struct {
	// expireAt 到期时间
	expireAt *time.Time
}

func (ev *expireValue) IsExpire() bool {
	if ev.expireAt == nil {
		return false
	}
	return ev.expireAt.Before(time.Now())
}

func (ev *expireValue) SetExpire(d time.Duration) {
	if d <= 0 {
		ev.SetExpireAt(nil)
		return
	}

	t := time.Now()
	t.Add(d)

	ev.SetExpireAt(&t)
}

func (ev *expireValue) SetExpireAt(t *time.Time) {
	// 如果是空值,直接设置成空
	if utils.IsNil(t) {
		ev.expireAt = nil
		return
	}

	ev.expireAt = t
}

type StoreType string

const (
	StoreTypeString    StoreType = "String"
	StoreTypeHash      StoreType = "Hash"
	StoreTypeList      StoreType = "List"
	StoreTypeSet       StoreType = "Set"
	StoreTypeSortedSet StoreType = "SortedSet"
)

type Cache struct {
	strStore stringStore
}

// checkKeyType 检测Key的类型是否正确
func (mc *Cache) checkKeyType(key string, useFor StoreType) (bool, error) {
	mc.strStore.valRWMutex.RLock()
	_, ok := mc.strStore.values[key]
	mc.strStore.valRWMutex.RUnlock()
	if ok {
		if useFor != StoreTypeString {
			return false, fmt.Errorf("该Key已经是String类型,不可用设置成%s类型", useFor)
		}
		return true, nil
	}
	return false, nil
}

// Exists 判断某个Key是否存在
func (mc *Cache) Exists(ctx context.Context, key string) (bool, error) {
	// 遍历字符串存储
	mc.strStore.valRWMutex.RLock()
	_, ok := mc.strStore.values[key]
	mc.strStore.valRWMutex.RUnlock()
	if ok {
		return true, nil
	}

	// 遍历Hash存储

	// 遍历List存储
	return false, nil
}

// Set 设置数据
func (mc *Cache) Set(ctx context.Context, key string, data any, expiration time.Duration) error {
	// 检测该Key是否被其他类型用了
	if _, err := mc.checkKeyType(key, StoreTypeString); err != nil {
		return err
	}

	return mc.strStore.Set(ctx, key, data, expiration)
}

// SetNX 设置数据,如果key不存在的话
func (mc *Cache) SetNX(ctx context.Context, key string, data any, expiration time.Duration) error {
	// 检测该Key是否被其他类型用了
	if _, err := mc.checkKeyType(key, StoreTypeString); err != nil {
		return err
	}

	return mc.strStore.SetNX(ctx, key, data, expiration)
}

// Get 获取数据
func (mc *Cache) Get(ctx context.Context, key string) (string, error) {
	if _, err := mc.checkKeyType(key, StoreTypeString); err != nil {
		return "", err
	}

	d, err := mc.strStore.Get(ctx, key)
	if err != nil {
		return "", err
	}
	return d, nil
}

// MGet 获取数据
func (mc *Cache) MGet(ctx context.Context, keys ...string) ([]any, error) {
	// 模拟redis
	return mc.strStore.MGet(ctx, keys...)
}

// MGet 获取数据
func (mc *Cache) MGetAndScan(ctx context.Context, dst any, keys ...string) error {
	// 模拟redis
	return mc.strStore.MGetAndScan(ctx, dst, keys...)
}

// marshalVal 解析数据
func marshalVal(data any) (string, error) {
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
