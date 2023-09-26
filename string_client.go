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
	BaseClient
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
		BaseClient: BaseClient{drivers: drs},
	}
}

// =======================================================
// ================= STRING ==============================
// =======================================================

// Set 设置数据
func (cli *StringClient) Set(ctx context.Context, key string, data interface{}, expiration time.Duration) driver.StatusValuer {
	ctx, _ = cli.preCheck(ctx)

	var value driver.StatusValuer
	for i, c := range cli.drivers {
		if v := c.(driver.String).Set(ctx, key, data, expiration); i == 0 {
			value = v
		}
	}
	return value
}

// SetNX 设置数据,如果key不存在的话
func (cli *StringClient) SetNX(ctx context.Context, key string, data interface{}, expiration time.Duration) driver.BoolValuer {
	ctx, _ = cli.preCheck(ctx)

	var value driver.BoolValuer
	for i, c := range cli.drivers {
		if v := c.(driver.String).SetNX(ctx, key, data, expiration); i == 0 {
			value = v
		}
	}
	return value
}

// Get 获取数据
func (cli *StringClient) Get(ctx context.Context, key string) driver.StringValuer {
	ctx, _ = cli.preCheck(ctx)
	var value driver.StringValuer
	for _, c := range cli.drivers {
		if value = c.(driver.String).Get(ctx, key); returnable(value) {
			return value
		}
	}
	return value
}

// MGet 获取多个Keys的值
func (cli *StringClient) MGet(ctx context.Context, keys ...string) driver.SliceValuer {
	ctx, _ = cli.preCheck(ctx)

	var value driver.SliceValuer
	for _, c := range cli.drivers {
		if value = c.(driver.String).MGet(ctx, keys...); returnable(value) {
			return value
		}
	}
	return value
}

// GetAndScan 获取并扫描
func (cli *StringClient) GetAndScan(ctx context.Context, dst interface{}, key string) error {
	return cli.Get(ctx, key).Scan(dst)
}

// MGetAndScan 获取多个Keys的值并扫描进dst中
func (cli *StringClient) MGetAndScan(ctx context.Context, dst interface{}, keys ...string) error {
	return cli.MGet(ctx, keys...).Scan(dst)
}
