package jcache

import (
	"context"

	"github.com/jerbe/jcache/v2/driver"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/13 13:40
  @describe :
*/

type ListClient struct {
	baseClient
}

func NewListClient(drivers ...driver.List) *ListClient {
	drs := make([]driver.Common, len(drivers))
	for i := 0; i < len(drivers); i++ {
		drs[i] = drivers[i]
	}

	if len(drs) == 0 {
		drs = append(drs, driver.NewMemory())
	}

	return &ListClient{
		baseClient{drivers: drs},
	}
}

// =======================================================
// ================= LIST ================================
// =======================================================

// LTrim 获取列表内的范围数据
func (cli *ListClient) LTrim(ctx context.Context, key string, start, stop int64) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(cli.drivers) == 0 {
		return ErrNoCacheClient
	}

	for _, c := range cli.drivers {
		c.(driver.List).LTrim(ctx, key, start, stop).Result()
	}

	return nil
}

// LPush 推送数据
func (cli *ListClient) LPush(ctx context.Context, key string, data ...interface{}) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(cli.drivers) == 0 {
		return ErrNoCacheClient
	}

	for _, c := range cli.drivers {
		c.(driver.List).LPush(ctx, key, data...)
	}

	return nil
}

// LRang 获取列表内的范围数据
func (cli *ListClient) LRang(ctx context.Context, key string, start, stop int64) ([]string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val, err := c.(driver.List).LRang(ctx, key, start, stop).Result()
		if err == nil {
			return val, nil
		}
	}
	return nil, ErrNoRecord
}

// LRangAndScan 通过扫描方式获取列表内的范围内数据
func (cli *ListClient) LRangAndScan(ctx context.Context, dst interface{}, key string, start, stop int64) error {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val := c.(driver.List).LRang(ctx, key, start, stop)
		if val.Err() == nil && val.ScanSlice(dst) == nil {
			return nil
		}
	}

	return ErrNoRecord
}

// LPop 移除并取出列表内的最后一个元素
func (cli *ListClient) LPop(ctx context.Context, key string) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val, err := c.(driver.List).LPop(ctx, key).Result()
		if err == nil {
			return val, nil
		}
	}

	return "", ErrNoRecord
}

// LPopAndScan 通过扫描方式移除并取出列表内的最后一个元素
func (cli *ListClient) LPopAndScan(ctx context.Context, dst interface{}, key string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val := c.(driver.List).LPop(ctx, key)
		if val.Err() == nil && val.Scan(dst) == nil {
			return nil
		}
	}

	return ErrNoRecord
}

// LShift 移除并取出列表内的第一个元素
func (cli *ListClient) LShift(ctx context.Context, key string) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val, err := c.(driver.List).LShift(ctx, key).Result()
		if err == nil {
			return val, nil
		}
	}

	return "", ErrNoRecord
}

// LShiftAndScan 通过扫描方式移除并取出列表内的第一个元素
func (cli *ListClient) LShiftAndScan(ctx context.Context, dst interface{}, key string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val := c.(driver.List).LShift(ctx, key)
		if val.Err() == nil && val.Scan(dst) == nil {
			return nil
		}
	}

	return ErrNoRecord
}

// LLen 返回列表长度
func (cli *ListClient) LLen(ctx context.Context, key string) (int64, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	for _, c := range cli.drivers {
		val := c.(driver.List).LLen(ctx, key)
		if val.Err() == nil {
			return val.Result()
		}
	}

	return 0, nil
}
