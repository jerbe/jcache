package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jerbe/jcache/internal/hscan"
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
	valRWMutex sync.RWMutex

	values map[string]*stringValue

	ekRWMutex sync.RWMutex

	expiredKeys map[string]any

	expireTicker *time.Ticker
}

func newStringStore() *stringStore {
	ticker := time.NewTicker(time.Second * 10)
	store := &stringStore{
		valRWMutex:   sync.RWMutex{},
		values:       make(map[string]*stringValue),
		ekRWMutex:    sync.RWMutex{},
		expiredKeys:  make(map[string]any),
		expireTicker: ticker,
	}
	go store.checkExpireTick()

	return store
}

// deleteExpiredKeys 删除过期的键
func (ss *stringStore) deleteExpiredKeys() {
	fmt.Println("ss.deleteExpiredKeys")
	ss.valRWMutex.Lock()
	ss.ekRWMutex.Lock()
	defer ss.ekRWMutex.Unlock()
	defer ss.valRWMutex.Unlock()

	for key, _ := range ss.expiredKeys {
		delete(ss.values, key)
		delete(ss.expiredKeys, key)
	}

	for key, value := range ss.values {
		if value.IsExpire() {
			delete(ss.values, key)
		}
	}
}

// checkExpireTick 检测到期的tick
func (ss *stringStore) checkExpireTick() {
	defer func() {
		if obj := recover(); obj != nil {
			go ss.checkExpireTick()
		}
	}()
	for {
		select {
		case <-ss.expireTicker.C:
			ss.deleteExpiredKeys()
		}
	}
}

func (ss *stringStore) setExpiredKey(key string) {
	ss.ekRWMutex.Lock()
	defer ss.ekRWMutex.Unlock()
	ss.expiredKeys[key] = nil
}

// Set 设置数据
func (ss *stringStore) Set(ctx context.Context, key string, data any, expiration time.Duration) error {
	ss.valRWMutex.Lock()
	defer ss.valRWMutex.Unlock()

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

		delete(ss.expiredKeys, key)
		ss.values[key] = value
	}
	return nil
}

// SetNX 设置数据,如果key不存在的话
func (ss *stringStore) SetNX(ctx context.Context, key string, data any, expiration time.Duration) error {
	ss.valRWMutex.RLock()
	_, ok := ss.values[key]
	ss.valRWMutex.RUnlock()
	if ok {
		return nil
	}

	return ss.Set(ctx, key, data, expiration)
}

// Get 获取数据
func (ss *stringStore) Get(ctx context.Context, key string) (string, error) {
	ss.valRWMutex.RLock()
	defer ss.valRWMutex.RUnlock()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		if d, ok := ss.values[key]; ok {
			if d.IsExpire() {
				ss.setExpiredKey(key)
				return "", Nil
			}
			return d.value, nil
		}
	}
	return "", Nil
}

// GetAndScan 获取数据并将输入扫入dst
func (ss *stringStore) GetAndScan(ctx context.Context, dst any, key string) error {
	ss.valRWMutex.RLock()
	defer ss.valRWMutex.RUnlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		bytes, err := ss.Get(ctx, key)
		if err != nil {
			return err
		}
		err = Scan([]byte(bytes), dst)
		if err != nil {
			return err
		}
		return nil
	}
}

// MGet 根据多个Key获取多个值
func (ss *stringStore) MGet(ctx context.Context, keys ...string) ([]any, error) {
	ss.valRWMutex.RLock()
	defer ss.valRWMutex.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var result = make([]any, 0, len(keys))
		for _, key := range keys {
			if d, ok := ss.values[key]; ok {
				if d.IsExpire() {
					result = append(result, nil)
					ss.setExpiredKey(key)
					continue
				}

				result = append(result, d.value)
			} else {
				result = append(result, nil)
			}
		}
		return result, nil
	}
}

// MGetAndScan 根据多个Key获取多个值
func (ss *stringStore) MGetAndScan(ctx context.Context, dst any, keys ...string) error {
	args := make([]any, len(keys))
	for i := 0; i < len(keys); i++ {
		args[i] = keys[i]
	}
	val, err := ss.MGet(ctx, keys...)
	if err != nil {
		return err
	}
	err = hscan.Scan(dst, args, val)
	if err != nil {
		return err
	}
	return nil
}
