package driver

import (
	"context"
	"sync"
	"time"

	"github.com/jerbe/jcache/utils"
)

/*
*

	@author : Jerbe - The porter from Earth
	@time : 2023/9/13 15:09
	@describe :
*/
const (
	// KeepTTL 保持原有的存活时间
	KeepTTL time.Duration = -1

	// ValueMaxTTL 数值最多可存活时长
	ValueMaxTTL = time.Hour * 6
)

// expireable 可以用于过期的
type expireable interface {
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

// IsExpire 是否已经过期
func (ev *expireValue) IsExpire() bool {
	if ev.expireAt == nil {
		// 没有到期时间,强制设定一个到期时间
		e := time.Now().Add(ValueMaxTTL)
		ev.expireAt = &e
		return false
	}
	if ev.expired {
		return true
	}
	ev.expired = ev.expireAt.Before(time.Now())
	return ev.expired
}

// SetExpire 设置可存活时长
func (ev *expireValue) SetExpire(d time.Duration) {
	ev.expired = false
	if d == KeepTTL {
		if ev.expireAt == nil {
			e := time.Now().Add(ValueMaxTTL)
			ev.expireAt = &e
		}
		return
	}

	if d > ValueMaxTTL {
		d = ValueMaxTTL
	}

	t := time.Now().Add(d)
	ev.expireAt = &t
}

// SetExpireAt 设置存活到期时间
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

type baseStoreer interface {
	Del(ctx context.Context, keys ...string) (int64, error)
	Exists(ctx context.Context, keys ...string) (int64, error)
	Expire(ctx context.Context, key string, ttl time.Duration) (bool, error)
	ExpireAt(ctx context.Context, key string, at time.Time) (bool, error)
	Persist(ctx context.Context, key string) (bool, error)
}

// baseStore 基础存储
type baseStore struct {
	values map[string]expireable

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

// keyExists 验证键是否存在
func (s *baseStore) keyExists(key string) bool {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()
	_, ok := s.values[key]
	return ok
}

// Del 删除指定键数量
func (s *baseStore) Del(ctx context.Context, keys ...string) (int64, error) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()
	cnt := int64(0)
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
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
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()
	cnt := int64(0)
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:

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
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
		v, ok := s.values[key]
		if !ok {
			return false, nil
		}
		v.SetExpire(ttl)
		return true, nil
	}
}

// ExpireAt 设置某个key在某个时间后失效
func (s *baseStore) ExpireAt(ctx context.Context, key string, at time.Time) (bool, error) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
		v, ok := s.values[key]
		if !ok {
			return false, nil
		}
		v.SetExpireAt(&at)
		return true, nil
	}
}

// Persist 设置某个key成为持久性的
func (s *baseStore) Persist(ctx context.Context, key string) (bool, error) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
		v, ok := s.values[key]
		if !ok {
			return false, nil
		}

		v.SetExpireAt(nil)
		return true, nil
	}
}
