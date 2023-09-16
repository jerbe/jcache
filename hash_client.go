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

type HashClient struct {
	baseClient
}

func NewHashClient(drivers ...driver.Hash) *HashClient {
	drvrs := make([]driver.Common, len(drivers))
	for i := 0; i < len(drivers); i++ {
		drvrs[i] = drivers[i]
	}

	return &HashClient{
		baseClient: baseClient{drivers: drvrs},
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
func (cli *HashClient) HSet(ctx context.Context, key string, values ...interface{}) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(cli.drivers) == 0 {
		return ErrNoCacheClient
	}

	for _, c := range cli.drivers {
		c.(driver.Hash).HSet(ctx, key, values)
	}

	return nil
}

// HVals 获取Hash表的所有值
func (cli *HashClient) HVals(ctx context.Context, key string) ([]string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		result, err := c.(driver.Hash).HVals(ctx, key).Result()
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
		result, err := c.(driver.Hash).HKeys(ctx, key).Result()
		if err == nil {
			return result, nil
		}
	}

	return nil, ErrNoRecord
}

// HKeysAndScan 获取Hash表的所有键并扫描到dst中
func (cli *HashClient) HKeysAndScan(ctx context.Context, dst interface{}, key string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val := c.(driver.Hash).HKeys(ctx, key)
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
		val := c.(driver.Hash).HLen(ctx, key)
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
		val := c.(driver.Hash).HGet(ctx, key, field)
		if val.Err() == nil {
			return val.Val(), nil
		}
	}

	return "", ErrNoRecord
}

// HGetAndScan 获取Hash表指定字段的值
func (cli *HashClient) HGetAndScan(ctx context.Context, dst interface{}, key, field string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val := c.(driver.Hash).HGet(ctx, key, field)
		if val.Err() == nil && val.Scan(dst) == nil {
			return nil
		}
	}
	return ErrNoRecord
}

// HMGet 获取Hash表指定字段的值
func (cli *HashClient) HMGet(ctx context.Context, key string, fields ...string) ([]interface{}, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		result, err := c.(driver.Hash).HMGet(ctx, key, fields...).Result()
		if err == nil {
			return result, nil
		}
	}

	return nil, ErrNoRecord
}

// HMGetAndScan 获取Hash表指定字段的值并扫描进入到dst中
func (cli *HashClient) HMGetAndScan(ctx context.Context, dst interface{}, key string, fields ...string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val := c.(driver.Hash).HMGet(ctx, key, fields...)
		if val.Err() == nil && val.Scan(dst) == nil {
			return nil
		}
	}

	return ErrNoRecord
}

// HValsAndScan 获取Hash表的所有值并扫如dst中
func (cli *HashClient) HValsAndScan(ctx context.Context, dst interface{}, key string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val := c.(driver.Hash).HVals(ctx, key)
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
		c.(driver.Hash).HDel(ctx, key, fields...)
	}
	return nil
}
