package jcache

import (
	"context"
	"github.com/jerbe/jcache/driver"
	"time"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/13 13:40
  @describe :
*/

type HashClient struct {
	drivers []driver.Hash
}

func NewHashClient(drivers ...driver.Hash) *HashClient {
	return &HashClient{drivers: drivers}
}

// =======================================================
// ================= COMMON ==============================
// =======================================================

// Exists 判断某个Key是否存在
func (cli *HashClient) Exists(ctx context.Context, keys ...string) (int64, error) {
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
func (cli *HashClient) Del(ctx context.Context, keys ...string) error {
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
func (cli *HashClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
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
func (cli *HashClient) ExpireAt(ctx context.Context, key string, at time.Time) error {
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
// ================= HASH ================================
// =======================================================

// HSet 写入hash数据
func (cli *HashClient) HSet(ctx context.Context, key string, values ...any) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(cli.drivers) == 0 {
		return ErrNoCacheClient
	}

	for _, c := range cli.drivers {
		c.HSet(ctx, key, values)
	}

	return nil
}

// HVals 获取Hash表的所有值
func (cli *HashClient) HVals(ctx context.Context, key string) ([]string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		result, err := c.HVals(ctx, key).Result()
		if err == nil {
			return result, nil
		}

	}

	return nil, ErrNoRecord
}

// HKeys 获取Hash表的所有键
func (cli *HashClient) HKeys(ctx context.Context, key string) ([]string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		result, err := c.HKeys(ctx, key).Result()
		if err == nil {
			return result, nil
		}
	}

	return nil, ErrNoRecord
}

// HKeysAndScan 获取Hash表的所有键并扫描到dst中
func (cli *HashClient) HKeysAndScan(ctx context.Context, dst any, key string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val := c.HKeys(ctx, key)
		if val.Err() == nil && val.ScanSlice(dst) == nil {
			return nil
		}
	}

	return ErrNoRecord
}

// HLen 获取Hash表的所有键个数
func (cli *HashClient) HLen(ctx context.Context, key string) (int64, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val := c.HLen(ctx, key)
		if val.Err() == nil {
			return val.Val(), nil
		}
	}

	return 0, ErrNoRecord
}

// HGet 获取Hash表指定字段的值
func (cli *HashClient) HGet(ctx context.Context, key, field string) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val := c.HGet(ctx, key, field)
		if val.Err() == nil {
			return val.Val(), nil
		}
	}

	return "", ErrNoRecord
}

// HGetAndScan 获取Hash表指定字段的值
func (cli *HashClient) HGetAndScan(ctx context.Context, dst any, key, field string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val := c.HGet(ctx, key, field)
		if val.Err() == nil && val.Scan(dst) == nil {
			return nil
		}
	}
	return ErrNoRecord
}

// HMGet 获取Hash表指定字段的值
func (cli *HashClient) HMGet(ctx context.Context, key string, fields ...string) ([]any, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		result, err := c.HMGet(ctx, key, fields...).Result()
		if err == nil {
			return result, nil
		}
	}

	return nil, ErrNoRecord
}

// HMGetAndScan 获取Hash表指定字段的值并扫描进入到dst中
func (cli *HashClient) HMGetAndScan(ctx context.Context, dst any, key string, fields ...string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val := c.HMGet(ctx, key, fields...)
		if val.Err() == nil && val.Scan(dst) == nil {
			return nil
		}
	}

	return ErrNoRecord
}

// HValsAndScan 获取Hash表的所有值并扫如dst中
func (cli *HashClient) HValsAndScan(ctx context.Context, dst any, key string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val := c.HVals(ctx, key)
		if val.Err() == nil && val.ScanSlice(dst) == nil {
			return nil
		}
	}

	return ErrNoRecord
}

// HDel 删除hash数据
func (cli *HashClient) HDel(ctx context.Context, key string, fields ...string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(cli.drivers) == 0 {
		return ErrNoCacheClient
	}

	for _, c := range cli.drivers {
		// @TODO 失败重做?
		c.HDel(ctx, key, fields...)
	}
	return nil
}
