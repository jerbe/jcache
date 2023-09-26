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
	BaseClient
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
		BaseClient{drivers: drs},
	}
}

// =======================================================
// ================= LIST ================================
// =======================================================

// LTrim 获取列表内的范围数据
func (cli *ListClient) LTrim(ctx context.Context, key string, start, stop int64) driver.StatusValuer {
	ctx, _ = cli.preCheck(ctx)

	var value driver.StatusValuer
	for i, c := range cli.drivers {
		if v := c.(driver.List).LTrim(ctx, key, start, stop); i == 0 {
			value = v
		}
	}

	return value
}

// LPush 推送数据
func (cli *ListClient) LPush(ctx context.Context, key string, data ...interface{}) driver.IntValuer {
	ctx, _ = cli.preCheck(ctx)

	var value driver.IntValuer
	for i, c := range cli.drivers {
		if v := c.(driver.List).LPush(ctx, key, data...); i == 0 {
			value = v
		}
	}
	return value
}

// LRang 获取列表内的范围数据
func (cli *ListClient) LRang(ctx context.Context, key string, start, stop int64) driver.StringSliceValuer {
	ctx, _ = cli.preCheck(ctx)

	var value driver.StringSliceValuer
	for _, c := range cli.drivers {
		if value = c.(driver.List).LRang(ctx, key, start, stop); returnable(value) {
			return value
		}
	}
	return value
}

// LRangAndScan 通过扫描方式获取列表内的范围内数据
func (cli *ListClient) LRangAndScan(ctx context.Context, dst interface{}, key string, start, stop int64) error {
	return cli.LRang(ctx, key, start, stop).ScanSlice(dst)
}

// LPop 移除并取出列表内的最后一个元素
func (cli *ListClient) LPop(ctx context.Context, key string) driver.StringValuer {
	ctx, _ = cli.preCheck(ctx)

	var value driver.StringValuer
	for _, c := range cli.drivers {
		if value = c.(driver.List).LPop(ctx, key); returnable(value) {
			return value
		}
	}
	return value
}

// LPopAndScan 通过扫描方式移除并取出列表内的最后一个元素
func (cli *ListClient) LPopAndScan(ctx context.Context, dst interface{}, key string) error {
	return cli.LPop(ctx, key).Scan(dst)
}

// LShift 移除并取出列表内的第一个元素
func (cli *ListClient) LShift(ctx context.Context, key string) driver.StringValuer {
	ctx, _ = cli.preCheck(ctx)

	var value driver.StringValuer
	for _, c := range cli.drivers {
		if value = c.(driver.List).LShift(ctx, key); returnable(value) {
			return value
		}
	}

	return value
}

// LShiftAndScan 通过扫描方式移除并取出列表内的第一个元素
func (cli *ListClient) LShiftAndScan(ctx context.Context, dst interface{}, key string) error {
	return cli.LShift(ctx, key).Scan(dst)
}

// LLen 返回列表长度
func (cli *ListClient) LLen(ctx context.Context, key string) driver.IntValuer {
	ctx, _ = cli.preCheck(ctx)
	var value driver.IntValuer
	for _, c := range cli.drivers {
		if value = c.(driver.List).LLen(ctx, key); returnable(value) {
			return value
		}
	}
	return value
}
