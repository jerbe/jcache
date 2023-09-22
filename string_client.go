package jcache

import (
	"context"
	"time"

	"github.com/jerbe/jcache/v2/driver"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/13 13:40
  @describe :
*/

type StringClient struct {
	baseClient
}

func NewStringClient(drivers ...driver.String) *StringClient {
	drs := make([]driver.Common, len(drivers))
	for i := 0; i < len(drivers); i++ {
		drs[i] = drivers[i]
	}

	if len(drs) == 0 {
		drs = append(drs, driver.NewMemory())
	}

	return &StringClient{
		baseClient: baseClient{drivers: drs},
	}
}

// =======================================================
// ================= STRING ==============================
// =======================================================

// Set 设置数据
func (cli *StringClient) Set(ctx context.Context, key string, data interface{}, expiration time.Duration) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(cli.drivers) == 0 {
		return ErrNoCacheClient
	}

	for _, c := range cli.drivers {
		c.(driver.String).Set(ctx, key, data, expiration)
	}
	return nil
}

// SetNX 设置数据,如果key不存在的话
func (cli *StringClient) SetNX(ctx context.Context, key string, data interface{}, expiration time.Duration) error {
	if len(cli.drivers) == 0 {
		return ErrNoCacheClient
	}

	if ctx == nil {
		ctx = context.Background()
	}

	for _, c := range cli.drivers {
		c.(driver.String).SetNX(ctx, key, data, expiration)
	}
	return nil
}

// Get 获取数据
func (cli *StringClient) Get(ctx context.Context, key string) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val, err := c.(driver.String).Get(ctx, key).Result()
		if err == nil {
			return val, err
		}
	}
	return "", ErrNoRecord
}

// MGet 获取多个Keys的值
func (cli *StringClient) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val, err := c.(driver.String).MGet(ctx, keys...).Result()
		if err == nil {
			return val, err
		}
	}
	return nil, ErrNoRecord
}

// GetAndScan 获取并扫描
func (cli *StringClient) GetAndScan(ctx context.Context, dst interface{}, key string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		err := c.(driver.String).Get(ctx, key).Scan(dst)
		if err == nil {
			return nil
		}
	}

	return ErrNoRecord
}

// MGetAndScan 获取多个Keys的值并扫描进dst中
func (cli *StringClient) MGetAndScan(ctx context.Context, dst interface{}, keys ...string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		err := c.(driver.String).MGet(ctx, keys...).Scan(dst)
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
		val, err := c.(driver.String).Get(ctx, key).Result()
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
func (cli *StringClient) CheckAndScan(ctx context.Context, dst interface{}, key string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val := c.(driver.String).Get(ctx, key)
		if val.Err() == nil && val.Val() == "" {
			return ErrEmpty
		}

		if err := val.Scan(dst); err == nil {
			return nil
		}
	}

	return ErrNoRecord
}
