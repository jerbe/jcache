package driver

import (
	"context"
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

type hashStore struct {
	baseStore
}

func newHashStore() *hashStore {
	ticker := time.NewTicker(time.Second * 10)
	store := &hashStore{
		baseStore: baseStore{
			values:       make(map[string]expirable),
			rwMutex:      sync.RWMutex{},
			expireTicker: ticker,
		},
	}

	go store.checkExpireTick()
	return store
}

func (s *hashStore) HSet(ctx context.Context, key string, data ...any) (int64, error) {
	return 0, nil
}

// HDel 哈希表删除指定字段(fields)
func (s *hashStore) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	return 0, nil
}

// HGet 哈希表获取一个数据
func (s *hashStore) HGet(ctx context.Context, key string, field string) (string, error) {
	return "", nil
}

// HMGet 哈希表获取多个数据
func (s *hashStore) HMGet(ctx context.Context, key string, fields ...string) ([]any, error) {
	return nil, nil
}

// HKeys 哈希表获取某个Key的所有字段(field)
func (s *hashStore) HKeys(ctx context.Context, key string) ([]string, error) {
	return nil, nil
}

// HVals 哈希表获取所有值
func (s *hashStore) HVals(ctx context.Context, key string) ([]string, error) {
	return nil, nil
}

// HLen 哈希表所有字段的数量
func (s *hashStore) HLen(ctx context.Context, key string) (int64, error) {
	return 0, nil
}
