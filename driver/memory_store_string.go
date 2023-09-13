package driver

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
	value string
}

type stringStore struct {
	baseStore
}

func newStringStore() *stringStore {
	ticker := time.NewTicker(time.Second * 10)

	store := &stringStore{
		baseStore: baseStore{
			values:       make(map[string]expirable),
			rwMutex:      sync.RWMutex{},
			expireTicker: ticker,
		},
	}

	go store.checkExpireTick()

	return store
}

// Set 设置数据
func (ss *stringStore) Set(ctx context.Context, key string, data any, expiration time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		value := &stringValue{}
		var err error
		value.value, err = marshalData(data)
		if err != nil {
			return err
		}
		value.SetExpire(expiration)

		ss.rwMutex.Lock()
		defer ss.rwMutex.Unlock()

		ss.values[key] = value
	}

	return nil
}

// SetNX 设置数据,如果key不存在的话
func (ss *stringStore) SetNX(ctx context.Context, key string, data any, expiration time.Duration) (bool, error) {
	ss.rwMutex.RLock()
	defer ss.rwMutex.RUnlock()
	_, ok := ss.values[key]
	if ok {
		return false, nil
	}

	err := ss.Set(ctx, key, data, expiration)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Get 获取数据
func (ss *stringStore) Get(ctx context.Context, key string) (string, error) {
	select {
	case <-ctx.Done():
		return "", nil
	default:
		ss.rwMutex.RLock()
		defer ss.rwMutex.RUnlock()

		d, ok := ss.values[key]
		if !ok {
			return "", MemoryNil
		}

		if d.IsExpire() {
			return "", MemoryNil
		}
		return d.(*stringValue).value, nil
	}
}

// MGet 根据多个Key获取多个值
func (ss *stringStore) MGet(ctx context.Context, keys ...string) ([]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		ss.rwMutex.RLock()
		defer ss.rwMutex.RUnlock()
		var rst = make([]any, 0, len(keys))
		for _, key := range keys {
			if d, ok := ss.values[key]; ok {
				if d.IsExpire() {
					rst = append(rst, nil)
					continue
				}
				rst = append(rst, d.(*stringValue).value)
			} else {
				rst = append(rst, nil)
			}
		}
		return rst, nil
	}
}
