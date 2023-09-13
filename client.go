package jcache

import (
	"context"
	"github.com/jerbe/jcache/driver"
	"time"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/8/27 18:08
  @describe :
*/

type Client struct {
	drivers []driver.Cache
}

func NewClient(drivers ...driver.Cache) *Client {
	return &Client{drivers: drivers}
}

// =======================================================
// ================= STRING ==============================
// =======================================================

// Set 设置数据
func (cli *Client) Set(ctx context.Context, key string, data any, expiration time.Duration) error {
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
func (cli *Client) SetNX(ctx context.Context, key string, data any, expiration time.Duration) error {
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
func (cli *Client) Get(ctx context.Context, key string) (string, error) {
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
func (cli *Client) MGet(ctx context.Context, keys ...string) ([]any, error) {
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
func (cli *Client) GetAndScan(ctx context.Context, dst any, key string) error {
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
func (cli *Client) MGetAndScan(ctx context.Context, dst any, keys ...string) error {
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
func (cli *Client) CheckAndGet(ctx context.Context, key string) (string, error) {
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
func (cli *Client) CheckAndScan(ctx context.Context, dst any, key string) error {
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

// =======================================================
// ================= HASH ================================
// =======================================================

// HSet 写入hash数据
func (cli *Client) HSet(ctx context.Context, key string, values ...any) error {
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
func (cli *Client) HVals(ctx context.Context, key string) ([]string, error) {
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
func (cli *Client) HKeys(ctx context.Context, key string) ([]string, error) {
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
func (cli *Client) HKeysAndScan(ctx context.Context, dst any, key string) error {
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
func (cli *Client) HLen(ctx context.Context, key string) (int64, error) {
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
func (cli *Client) HGet(ctx context.Context, key, field string) (string, error) {
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
func (cli *Client) HGetAndScan(ctx context.Context, dst any, key, field string) error {
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
func (cli *Client) HMGet(ctx context.Context, key string, fields ...string) ([]any, error) {
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
func (cli *Client) HMGetAndScan(ctx context.Context, dst any, key string, fields ...string) error {
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
func (cli *Client) HValsAndScan(ctx context.Context, dst any, key string) error {
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
func (cli *Client) HDel(ctx context.Context, key string, fields ...string) error {
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

// =======================================================
// ================= LIST ================================
// =======================================================

// Trim 获取列表内的范围数据
func (cli *Client) Trim(ctx context.Context, key string, start, stop int64) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(cli.drivers) == 0 {
		return ErrNoCacheClient
	}

	for _, c := range cli.drivers {
		c.Trim(ctx, key, start, stop).Result()
	}

	return nil
}

// Push 推送数据
func (cli *Client) Push(ctx context.Context, key string, data ...any) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(cli.drivers) == 0 {
		return ErrNoCacheClient
	}

	for _, c := range cli.drivers {
		c.Push(ctx, key, data...)
	}

	return nil
}

// Rang 获取列表内的范围数据
func (cli *Client) Rang(ctx context.Context, key string, start, stop int64) ([]string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val, err := c.Rang(ctx, key, start, stop).Result()
		if err == nil {
			return val, nil
		}
	}
	return nil, ErrNoRecord
}

// RangAndScan 通过扫描方式获取列表内的范围内数据
func (cli *Client) RangAndScan(ctx context.Context, dst any, key string, start, stop int64) error {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val := c.Rang(ctx, key, start, stop)
		if val.Err() == nil && val.ScanSlice(dst) == nil {
			return nil
		}
	}

	return ErrNoRecord
}

// Pop 取出列表内的第一个数据
func (cli *Client) Pop(ctx context.Context, key string) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val, err := c.Pop(ctx, key).Result()
		if err == nil {
			return val, nil
		}
	}

	return "", ErrNoRecord
}

// PopAndScan 通过扫描方式取出列表内的第一个数据
func (cli *Client) PopAndScan(ctx context.Context, dst any, key string) error {
	if ctx == nil {
		ctx = context.Background()
	}

	for _, c := range cli.drivers {
		val := c.Pop(ctx, key)
		if val.Err() == nil && val.Scan(dst) == nil {
			return nil
		}
	}

	return ErrNoRecord
}

// =======================================================
// ================= COMMON ==============================
// =======================================================

// Exists 判断某个Key是否存在
func (cli *Client) Exists(ctx context.Context, keys ...string) (int64, error) {
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
func (cli *Client) Del(ctx context.Context, keys ...string) error {
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
func (cli *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
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
func (cli *Client) ExpireAt(ctx context.Context, key string, at time.Time) error {
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
