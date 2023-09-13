package driver

import (
	"context"
	"github.com/jerbe/jcache/utils"
	"sync"
	"time"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/13 15:09
  @describe :
*/

// ValueMaxTTL 数值最多可存活时长
const ValueMaxTTL = time.Hour * 6

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

	expired bool
}

func (ev *expireValue) IsExpire() bool {
	if ev.expireAt == nil {
		return false
	}
	if ev.expired {
		return true
	}
	ev.expired = ev.expireAt.Before(time.Now())
	return ev.expired
}

func (ev *expireValue) SetExpire(d time.Duration) {
	ev.expired = false
	if d <= 0 {
		ev.SetExpireAt(nil)
		return
	}

	t := time.Now().Add(d)

	ev.SetExpireAt(&t)
}

func (ev *expireValue) SetExpireAt(t *time.Time) {
	ev.expired = false

	// 限定一个value最多只能存活 ValueMaxTTL 时

	// 如果是空值,直接设置成空
	if utils.IsNil(t) {
		e := time.Now().Add(ValueMaxTTL)
		ev.expireAt = &e
		return
	}

	maxTTL := time.Now().Add(ValueMaxTTL)
	if t.After(maxTTL) {
		t = &maxTTL
	}
	ev.expireAt = t
}

// baseStore 基础存储
type baseStore struct {
	values map[string]expirable

	rwMutex sync.RWMutex

	expireTicker *time.Ticker
}

// deleteExpiredKeys 删除过期的键
func (s *baseStore) deleteExpiredKeys() {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	// 移除已经被标记成过期的key
	for k, v := range s.values {
		if v.IsExpire() {
			delete(s.values, k)
		}
	}
}

// checkExpireTick 检测到期的tick
func (s *baseStore) checkExpireTick() {
	defer func() {
		if obj := recover(); obj != nil {
			go s.checkExpireTick()
		}
	}()
	for {
		select {
		case <-s.expireTicker.C:
			s.deleteExpiredKeys()
		}
	}
}

// Del 删除指定键数量
func (s *baseStore) Del(ctx context.Context, keys ...string) (int64, error) {
	cnt := int64(0)
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		s.rwMutex.Lock()
		defer s.rwMutex.Unlock()
		for _, key := range keys {
			if _, ok := s.values[key]; ok {
				delete(s.values, key)
				cnt++
			}
		}
		return cnt, nil
	}
}

// Exists 判断键是否存在
func (s *baseStore) Exists(ctx context.Context, keys ...string) (int64, error) {
	cnt := int64(0)
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		s.rwMutex.RLock()
		defer s.rwMutex.RUnlock()
		for _, key := range keys {
			if _, ok := s.values[key]; ok {
				cnt++
			}
		}
		return cnt, nil
	}
}

// Expire 设置某个key的存活时间
func (s *baseStore) Expire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
		s.rwMutex.Lock()
		defer s.rwMutex.Unlock()

		v, ok := s.values[key]
		if !ok {
			return false, MemoryNil
		}
		v.SetExpire(ttl)
		return true, nil
	}
}

// ExpireAt 设置某个key在某个时间后失效
func (s *baseStore) ExpireAt(ctx context.Context, key string, at time.Time) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
		s.rwMutex.Lock()
		defer s.rwMutex.Unlock()

		v, ok := s.values[key]
		if !ok {
			return false, MemoryNil
		}

		v.SetExpireAt(&at)
		return true, nil
	}
}
