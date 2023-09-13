package jcache

import (
	"context"
	"time"

	"github.com/jerbe/jcache/driver"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/13 13:40
  @describe :
*/

type StringClient struct {
	drivers []driver.String
}

func NewStringClient(drivers ...driver.String) *StringClient {
	return &StringClient{drivers: drivers}
}

// =======================================================
// ================= COMMON ==============================
// =======================================================

// Exists 判断某个Key是否存在
func (cli *StringClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val := c.Exists(ctx, keys...)
		if val.Err() == nil {
			return val.Result()
		}
	}
	return 0, ErrNoRecord
}

// Del 删除键
func (cli *StringClient) Del(ctx context.Context, keys ...string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(cli.drivers) == 0 {
		return ErrNoCacheClient
	}

	for _, c := range cli.drivers {
		c.Del(ctx, keys...)
	}

	return nil
}

// Expire 设置某个Key的TTL时长
func (cli *StringClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(cli.drivers) == 0 {
		return ErrNoCacheClient
	}

	for _, c := range cli.drivers {
		c.Expire(ctx, key, expiration)
	}

	// @TODO 其他缓存方法
	return nil
}

// ExpireAt 设置某个key在指定时间内到期
func (cli *StringClient) ExpireAt(ctx context.Context, key string, at time.Time) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(cli.drivers) == 0 {
		return ErrNoCacheClient
	}

	for _, c := range cli.drivers {
		c.ExpireAt(ctx, key, &at)
	}

	// @TODO 其他缓存方法
	return nil
}

// =======================================================
// ================= STRING ==============================
// =======================================================

// Set 设置数据
func (cli *StringClient) Set(ctx context.Context, key string, data any, expiration time.Duration) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(cli.drivers) == 0 {
		return ErrNoCacheClient
	}

	for _, c := range cli.drivers {
		c.Set(ctx, key, data, expiration)
	}
	return nil
}

// SetNX 设置数据,如果key不存在的话
func (cli *StringClient) SetNX(ctx context.Context, key string, data any, expiration time.Duration) error {
	if len(cli.drivers) == 0 {
		return ErrNoCacheClient
	}

	if ctx == nil {
		ctx = context.Background()
	}

	for _, c := range cli.drivers {
		c.SetNX(ctx, key, data, expiration)
	}
	return nil
}

// Get 获取数据
func (cli *StringClient) Get(ctx context.Context, key string) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val, err := c.Get(ctx, key).Result()
		if err == nil {
			return val, err
		}
	}
	return "", ErrNoRecord
}

// MGet 获取多个Keys的值
func (cli *StringClient) MGet(ctx context.Context, keys ...string) ([]any, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val, err := c.MGet(ctx, keys...).Result()
		if err == nil {
			return val, err
		}
	}
	return nil, ErrNoRecord
}

// GetAndScan 获取并扫描
func (cli *StringClient) GetAndScan(ctx context.Context, dst any, key string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		err := c.Get(ctx, key).Scan(dst)
		if err == nil {
			return nil
		}
	}

	return ErrNoRecord
}

// MGetAndScan 获取多个Keys的值并扫描进dst中
func (cli *StringClient) MGetAndScan(ctx context.Context, dst any, keys ...string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		err := c.MGet(ctx, keys...).Scan(dst)
		if err == nil {
			return nil
		}
	}
	return ErrNoRecord
}

// CheckAndGet 检测并获取数据
func (cli *StringClient) CheckAndGet(ctx context.Context, key string) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val, err := c.Get(ctx, key).Result()
		if err == nil && val == "" {
			return "", ErrEmpty
		}
		if err == nil {
			return val, nil
		}
	}
	return "", ErrNoRecord
}

// CheckAndScan 获取数据
func (cli *StringClient) CheckAndScan(ctx context.Context, dst any, key string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val := c.Get(ctx, key)
		if val.Err() == nil && val.Val() == "" {
			return ErrEmpty
		}

		if err := val.Scan(dst); err == nil {
			return nil
		}
	}

	return ErrNoRecord
}
