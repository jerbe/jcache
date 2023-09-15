package driver

import (
	"context"
	"errors"
	"sync"
	"time"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/13 14:55
  @describe :
*/

// hashValue 字符串值
type hashValue struct {
	expireValue

	// value 值
	value map[string]string
}

func newHashValue() *hashValue {
	defaultExpireAt := time.Now().Add(ValueMaxTTL)
	return &hashValue{
		expireValue: expireValue{
			expireAt: &defaultExpireAt,
			expired:  false,
		},
		value: make(map[string]string),
	}
}

type hashStore struct {
	baseStore
}

func newHashStore() *hashStore {
	ticker := time.NewTicker(time.Second * 10)
	store := &hashStore{
		baseStore: baseStore{
			values:       make(map[string]expireable),
			rwMutex:      sync.RWMutex{},
			expireTicker: ticker,
		},
	}

	go store.checkExpireTick()
	return store
}

// HSet 写入hash数据
// 接受以下格式的值：
// HSet("myhash", "key1", "value1", "key2", "value2")
//
// HSet("myhash", []string{"key1", "value1", "key2", "value2"})
//
// HSet("myhash", map[string]interface{}{"key1": "value1", "key2": "value2"})
// 使用“redis”标签播放结构。 type MyHash struct { Key1 string `redis:"key1"`; Key2 int `redis:"key2"` }
//
// HSet("myhash", MyHash{"value1", "value2"}) 警告：redis-server >= 4.0
// 对于struct，可以是结构体指针类型，我们只解析标签为redis的字段。如果你不想读取该字段，可以使用 `redis:"-"` 标志来忽略它，或者不需要设置 redis 标签。对于结构体字段的类型，我们只支持简单的数据类型：string、int/uint(8,16,32,64)、float(32,64)、time.Time(to RFC3339Nano)、time.Duration(to Nanoseconds) ），如果是其他更复杂或者自定义的数据类型，请实现encoding.BinaryMarshaler接口。
func (s *hashStore) HSet(ctx context.Context, key string, data ...interface{}) (int64, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:

		val, ok := s.values[key].(*hashValue)
		if !ok {
			val = newHashValue()
		}
		result := sliceArgs(data)
		if len(result)%2 != 0 {
			return 0, errors.New("the number of parameters is incorrect")
		}

		newCnt := int64(0)
		for i := 0; i < len(result); i += 2 {
			field, err := marshalData(result[i])
			if err != nil {
				return 0, err
			}

			if _, ok := val.value[field]; !ok {
				newCnt++
			}

			value, err := marshalData(result[i+1])
			if err != nil {
				return 0, err
			}
			val.value[field] = value
		}

		s.values[key] = val

		return newCnt, nil
	}
}

// HDel 哈希表删除指定字段(fields)
func (s *hashStore) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		val, ok := s.values[key].(*hashValue)
		if !ok {
			return 0, nil
		}

		affectsCnt := int64(0)
		for _, field := range fields {
			if _, ok := val.value[field]; ok {
				affectsCnt++
				delete(val.value, field)
			}
		}

		return affectsCnt, nil
	}
}

// HGet 哈希表获取一个数据
func (s *hashStore) HGet(ctx context.Context, key string, field string) (string, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		val, ok := s.values[key].(*hashValue)
		if !ok {
			return "", MemoryNil
		}

		str, ok := val.value[field]
		if !ok {
			return "", MemoryNil
		}

		return str, nil
	}

}

// HMGet 哈希表获取多个数据
func (s *hashStore) HMGet(ctx context.Context, key string, fields ...string) ([]interface{}, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		rest := make([]interface{}, 0)

		val, ok := s.values[key].(*hashValue)
		if !ok {
			return rest, MemoryNil
		}

		// 不能使用for rang, 因为它是无序的
		for i := 0; i < len(fields); i++ {
			field := fields[i]
			str, ok := val.value[field]
			if ok {
				rest = append(rest, str)
			} else {
				rest = append(rest, nil)
			}
		}

		return rest, nil
	}
}

// HKeys 哈希表获取某个Key的所有字段(field)
func (s *hashStore) HKeys(ctx context.Context, key string) ([]string, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		val, ok := s.values[key].(*hashValue)
		if !ok {
			return []string{}, MemoryNil
		}

		rest := make([]string, len(val.value))
		var i int
		for field, _ := range val.value {
			rest[i] = field
			i++
		}
		return rest, nil
	}
}

// HVals 哈希表获取所有值
func (s *hashStore) HVals(ctx context.Context, key string) ([]string, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	select {
	case <-ctx.Done():
		return []string{}, ctx.Err()
	default:
		val, ok := s.values[key].(*hashValue)
		if !ok {
			return []string{}, MemoryNil
		}

		rest := make([]string, len(val.value))
		var i int
		for _, value := range val.value {
			rest[i] = value
			i++
		}
		return rest, nil
	}
}

// HLen 哈希表所有字段的数量
func (s *hashStore) HLen(ctx context.Context, key string) (int64, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		val, ok := s.values[key].(*hashValue)
		if !ok {
			return 0, MemoryNil
		}

		return int64(len(val.value)), nil
	}
}
