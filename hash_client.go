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

type HashClient struct {
	BaseClient
}

func NewHashClient(drivers ...driver.Hash) *HashClient {
	drs := make([]driver.Common, 0)
	for i := 0; i < len(drivers); i++ {
		drs = append(drs, drivers[i])
	}

	if len(drs) == 0 {
		drs = append(drs, driver.NewMemory())
	}

	return &HashClient{
		BaseClient: BaseClient{drivers: drs},
	}
}

// =======================================================
// ================= HASH ================================
// =======================================================

// HSet 写入hash数据
// 接受以下格式的值：
// HSet("myhash", "key1", "value1", "key2", "value2")
//
// HSet("myhash", []string{"key1", "value1", "key2", "value2"})
//
// HSet("myhash", map[string]interface{}{"key1": "value1", "key2": "value2"})
// 使用“redis”标签播放结构。 type MyHash struct { Key1 string `redis:"key1"`; Key2 int `redis:"key2"` }
//
// HSet("myhash", MyHash{"value1", "value2"}) 警告：redis-server >= 4.0
// 对于struct，可以是结构体指针类型，我们只解析标签为redis的字段。如果你不想读取该字段，可以使用 `redis:"-"` 标志来忽略它，或者不需要设置 redis 标签。对于结构体字段的类型，我们只支持简单的数据类型：string、int/uint(8,16,32,64)、float(32,64)、time.Time(to RFC3339Nano)、time.Duration(to Nanoseconds) ），如果是其他更复杂或者自定义的数据类型，请实现encoding.BinaryMarshaler接口。
func (cli *HashClient) HSet(ctx context.Context, key string, values ...interface{}) driver.IntValuer {
	ctx, _ = cli.preCheck(ctx)

	var value driver.IntValuer
	for i, c := range cli.drivers {
		if v := c.(driver.Hash).HSet(ctx, key, values); i == 0 {
			value = v
		}
	}

	return value
}

// HSetNX 哈希表设置某个字段的值,如果存在的话返回true
func (cli *HashClient) HSetNX(ctx context.Context, key, field string, data interface{}) driver.BoolValuer {
	ctx, _ = cli.preCheck(ctx)
	var value driver.BoolValuer
	for i, c := range cli.drivers {
		if v := c.(driver.Hash).HSetNX(ctx, key, field, data); i == 0 {
			value = v
		}
	}

	return value
}

// HVals 获取Hash表的所有值
func (cli *HashClient) HVals(ctx context.Context, key string) driver.StringSliceValuer {
	ctx, _ = cli.preCheck(ctx)
	var value driver.StringSliceValuer
	for _, c := range cli.drivers {
		if value = c.(driver.Hash).HVals(ctx, key); returnable(value) {
			return value
		}
	}
	return value
}

// HKeys 获取Hash表的所有键
func (cli *HashClient) HKeys(ctx context.Context, key string) driver.StringSliceValuer {
	ctx, _ = cli.preCheck(ctx)

	var value driver.StringSliceValuer
	for _, c := range cli.drivers {
		if value = c.(driver.Hash).HKeys(ctx, key); returnable(value) {
			return value
		}

	}
	return value
}

// HGetAll 获取哈希表中所有的值,包括键/值
func (cli *HashClient) HGetAll(ctx context.Context, key string) driver.MapStringStringValuer {
	ctx, _ = cli.preCheck(ctx)

	var value driver.MapStringStringValuer
	for _, c := range cli.drivers {
		if value = c.(driver.Hash).HGetAll(ctx, key); returnable(value) {
			return value
		}
	}
	return value
}

// HKeysAndScan 获取Hash表的所有键并扫描到dst中
func (cli *HashClient) HKeysAndScan(ctx context.Context, dst interface{}, key string) error {
	return cli.HKeys(ctx, key).ScanSlice(dst)
}

// HLen 获取Hash表的所有键个数
func (cli *HashClient) HLen(ctx context.Context, key string) driver.IntValuer {
	ctx, _ = cli.preCheck(ctx)

	var value driver.IntValuer
	for _, c := range cli.drivers {
		if value = c.(driver.Hash).HLen(ctx, key); returnable(value) {
			return value
		}
	}

	return value
}

// HGet 获取Hash表指定字段的值
func (cli *HashClient) HGet(ctx context.Context, key, field string) driver.StringValuer {
	ctx, _ = cli.preCheck(ctx)

	var value driver.StringValuer
	for _, c := range cli.drivers {
		if value = c.(driver.Hash).HGet(ctx, key, field); returnable(value) {
			return value
		}
	}
	return value
}

// HGetAndScan 获取Hash表指定字段的值
func (cli *HashClient) HGetAndScan(ctx context.Context, dst interface{}, key, field string) error {
	return cli.HGet(ctx, key, field).Scan(dst)
}

// HMGet 获取Hash表指定字段的值
func (cli *HashClient) HMGet(ctx context.Context, key string, fields ...string) driver.SliceValuer {
	ctx, _ = cli.preCheck(ctx)
	var value driver.SliceValuer
	for _, c := range cli.drivers {
		if value = c.(driver.Hash).HMGet(ctx, key, fields...); returnable(value) {
			return value
		}
	}
	return value
}

// HMGetAndScan 获取Hash表指定字段的值并扫描进入到dst中
func (cli *HashClient) HMGetAndScan(ctx context.Context, dst interface{}, key string, fields ...string) error {
	return cli.HMGet(ctx, key, fields...).Scan(dst)
}

// HValsAndScan 获取Hash表的所有值并扫如dst中
func (cli *HashClient) HValsAndScan(ctx context.Context, dst interface{}, key string) error {
	return cli.HVals(ctx, key).ScanSlice(dst)
}

// HDel 删除hash数据
func (cli *HashClient) HDel(ctx context.Context, key string, fields ...string) driver.IntValuer {
	ctx, _ = cli.preCheck(ctx)

	var value driver.IntValuer
	for i, c := range cli.drivers {
		// @TODO 失败重做?
		if v := c.(driver.Hash).HDel(ctx, key, fields...); i == 0 {
			value = v
		}
	}
	return value
}

// HExists 判断哈希表周公某个字段是否存在
func (cli *HashClient) HExists(ctx context.Context, key, field string) driver.BoolValuer {
	ctx, _ = cli.preCheck(ctx)

	var value driver.BoolValuer
	for _, c := range cli.drivers {
		// @TODO 失败重做?
		if value = c.(driver.Hash).HExists(ctx, key, field); returnable(value) {
			return value
		}
	}
	return value
}
