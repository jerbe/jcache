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

type ListClient struct {
	drivers []driver.List
}

func NewListClient(drivers ...driver.List) *ListClient {
	return &ListClient{drivers: drivers}
}

// =======================================================
// ================= COMMON ==============================
// =======================================================

// Exists 判断某个Key是否存在
func (cli *ListClient) Exists(ctx context.Context, keys ...string) (int64, error) {
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
func (cli *ListClient) Del(ctx context.Context, keys ...string) error {
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
func (cli *ListClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
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
func (cli *ListClient) ExpireAt(ctx context.Context, key string, at time.Time) error {
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
// ================= LIST ================================
// =======================================================

// Trim 获取列表内的范围数据
func (cli *ListClient) Trim(ctx context.Context, key string, start, stop int64) error {
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
func (cli *ListClient) Push(ctx context.Context, key string, data ...any) error {
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
func (cli *ListClient) Rang(ctx context.Context, key string, start, stop int64) ([]string, error) {
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
func (cli *ListClient) RangAndScan(ctx context.Context, dst any, key string, start, stop int64) error {
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
func (cli *ListClient) Pop(ctx context.Context, key string) (string, error) {
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
func (cli *ListClient) PopAndScan(ctx context.Context, dst any, key string) error {
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
