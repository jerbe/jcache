package memory

import (
	"context"
	"sync"
	"time"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/8 17:50
  @describe :
*/

// stringValue 字符串值
type stringValue struct {
	expireValue

	// value 值
	value []byte
}

type stringStore struct {
	rwLock sync.RWMutex

	values map[string]*stringValue
}

func newStringStore() *stringStore {
	return &stringStore{
		rwLock: sync.RWMutex{},
		values: make(map[string]*stringValue),
	}
}

// Set 设置数据
func (ss *stringStore) Set(ctx context.Context, key string, data any, expiration time.Duration) error {
	ss.rwLock.Lock()
	defer ss.rwLock.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		value := &stringValue{}
		var err error
		value.value, err = marshalVal(data)
		if err != nil {
			return err
		}
		value.SetExpire(expiration)
		ss.values[key] = value
	}
	return nil
}

// SetNX 设置数据,如果key不存在的话
func (ss *stringStore) SetNX(ctx context.Context, key string, data any, expiration time.Duration) error {
	ss.rwLock.RLock()
	_, ok := ss.values[key]
	ss.rwLock.RUnlock()
	if ok {
		return nil
	}

	return ss.Set(ctx, key, data, expiration)
}

// Get 获取数据
func (ss *stringStore) Get(ctx context.Context, key string) ([]byte, error) {
	ss.rwLock.RLock()
	defer ss.rwLock.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		if d, ok := ss.values[key]; ok {
			return d.value, nil
		}
	}
	return nil, Nil
}

// MGet 根据多个Key获取多个值
func (ss *stringStore) MGet(ctx context.Context, keys ...string) ([][]byte, error) {
	ss.rwLock.RLock()
	defer ss.rwLock.RUnlock()

	var result = make([][]byte, 0, len(keys))
	for _, key := range keys {
		if d, ok := ss.values[key]; ok {
			result = append(result, d.value)
		}
	}

	if len(result) == 0 {
		return nil, Nil
	}

	return result, nil
}
