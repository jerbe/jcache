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

func newStringValue() *stringValue {
	defaultExpireAt := time.Now().Add(ValueMaxTTL)
	return &stringValue{
		expireValue: expireValue{
			expireAt: &defaultExpireAt,
			expired:  false,
		},
	}
}

type stringStore struct {
	baseStore
}

func newStringStore() *stringStore {
	ticker := time.NewTicker(time.Second * 10)

	store := &stringStore{
		baseStore: baseStore{
			values:       make(map[string]expireable),
			rwMutex:      sync.RWMutex{},
			expireTicker: ticker,
		},
	}

	// 定时检测到期key
	go store.checkExpireTick()

	return store
}

// Set 设置数据
func (ss *stringStore) Set(ctx context.Context, key string, data interface{}, expiration time.Duration) error {
	ss.rwMutex.Lock()
	defer ss.rwMutex.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		val, ok := ss.values[key].(*stringValue)
		if !ok {
			val = newStringValue()
		}

		var err error
		val.value, err = marshalData(data)
		if err != nil {
			return err
		}

		if ok && expiration != KeepTTL {
			val.SetExpire(expiration)
		}

		ss.values[key] = val
	}

	return nil
}

// SetNX 设置数据,如果key不存在的话
func (ss *stringStore) SetNX(ctx context.Context, key string, data interface{}, expiration time.Duration) (bool, error) {
	ss.rwMutex.RLock()
	defer ss.rwMutex.RUnlock()

	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
		val, ok := ss.values[key].(*stringValue)
		if ok {
			return false, nil
		}

		val = newStringValue()
		var err error
		val.value, err = marshalData(data)
		if err != nil {
			return false, err
		}

		if ok && expiration != KeepTTL {
			val.SetExpire(expiration)
		}
		ss.values[key] = val
		return true, nil
	}
}

// Get 获取数据
func (ss *stringStore) Get(ctx context.Context, key string) (string, error) {
	ss.rwMutex.RLock()
	defer ss.rwMutex.RUnlock()

	select {
	case <-ctx.Done():
		return "", nil
	default:
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
func (ss *stringStore) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	ss.rwMutex.RLock()
	defer ss.rwMutex.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var rst = make([]interface{}, len(keys), len(keys))
		for i, key := range keys {
			if d, ok := ss.values[key]; ok {
				if d.IsExpire() {
					continue
				}
				rst[i] = d.(*stringValue).value
			}
		}
		return rst, nil
	}
}
