package jcache

import (
	"context"
	"github.com/jerbe/jcache/driver"
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
	drvrs := make([]driver.Common, len(drivers))
	for i := 0; i < len(drivers); i++ {
		drvrs[i] = drivers[i]
	}

	return &ListClient{
		baseClient{drivers: drvrs},
	}
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
		c.(driver.List).Trim(ctx, key, start, stop).Result()
	}

	return nil
}

// Push 推送数据
func (cli *ListClient) Push(ctx context.Context, key string, data ...interface{}) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(cli.drivers) == 0 {
		return ErrNoCacheClient
	}

	for _, c := range cli.drivers {
		c.(driver.List).Push(ctx, key, data...)
	}

	return nil
}

// Rang 获取列表内的范围数据
func (cli *ListClient) Rang(ctx context.Context, key string, start, stop int64) ([]string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val, err := c.(driver.List).Rang(ctx, key, start, stop).Result()
		if err == nil {
			return val, nil
		}
	}
	return nil, ErrNoRecord
}

// RangAndScan 通过扫描方式获取列表内的范围内数据
func (cli *ListClient) RangAndScan(ctx context.Context, dst interface{}, key string, start, stop int64) error {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val := c.(driver.List).Rang(ctx, key, start, stop)
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
		val, err := c.(driver.List).Pop(ctx, key).Result()
		if err == nil {
			return val, nil
		}
	}

	return "", ErrNoRecord
}

// PopAndScan 通过扫描方式取出列表内的第一个数据
func (cli *ListClient) PopAndScan(ctx context.Context, dst interface{}, key string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val := c.(driver.List).Pop(ctx, key)
		if val.Err() == nil && val.Scan(dst) == nil {
			return nil
		}
	}

	return ErrNoRecord
}
