package jcache

import (
	"context"
	"time"

	"github.com/jerbe/jcache/driver"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/8/27 18:08
  @describe :
*/

// =======================================================
// ================= BaseClient ==========================
// =======================================================

type baseClient struct {
	drivers []driver.Common
}

// Exists 判断某个Key是否存在
func (cli *baseClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val := c.Exists(ctx, keys...)
		if val.Err() == nil {
			return val.Result()
		}
	}
	return 0, nil
}

// Del 删除键
func (cli *baseClient) Del(ctx context.Context, keys ...string) error {
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
func (cli *baseClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
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
func (cli *baseClient) ExpireAt(ctx context.Context, key string, at time.Time) error {
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
// ================= Client ==============================
// =======================================================

type Client struct {
	baseClient
	StringClient
	HashClient
	ListClient
}

func NewClient(drivers ...driver.Cache) *Client {
	drvrs := make([]driver.Common, len(drivers))

	for i := 0; i < len(drivers); i++ {
		drvrs[i] = drivers[i]
	}

	bcli := baseClient{drivers: drvrs}

	return &Client{
		baseClient:   bcli,
		StringClient: StringClient{baseClient: bcli},
		HashClient:   HashClient{baseClient: bcli},
		ListClient:   ListClient{baseClient: bcli},
	}
}

/*

// =======================================================
// ================= STRING ==============================
// =======================================================

// Set 设置数据
func (cli *Client) Set(ctx context.Context, key string, data interface{}, expiration time.Duration) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(cli.drivers) == 0 {
		return ErrNoCacheClient
	}

	for _, c := range cli.drivers {
		c.(driver.Cache).Set(ctx, key, data, expiration)
	}
	return nil
}

// SetNX 设置数据,如果key不存在的话
func (cli *Client) SetNX(ctx context.Context, key string, data interface{}, expiration time.Duration) error {
	if len(cli.drivers) == 0 {
		return ErrNoCacheClient
	}

	if ctx == nil {
		ctx = context.Background()
	}

	for _, c := range cli.drivers {
		c.(driver.Cache).SetNX(ctx, key, data, expiration)
	}
	return nil
}

// Get 获取数据
func (cli *Client) Get(ctx context.Context, key string) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val, err := c.(driver.Cache).Get(ctx, key).Result()
		if err == nil {
			return val, err
		}
	}
	return "", ErrNoRecord
}

// MGet 获取多个Keys的值
func (cli *Client) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val, err := c.(driver.Cache).MGet(ctx, keys...).Result()
		if err == nil {
			return val, err
		}
	}
	return nil, ErrNoRecord
}

// GetAndScan 获取并扫描
func (cli *Client) GetAndScan(ctx context.Context, dst interface{}, key string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		err := c.(driver.Cache).Get(ctx, key).Scan(dst)
		if err == nil {
			return nil
		}
	}

	return ErrNoRecord
}

// MGetAndScan 获取多个Keys的值并扫描进dst中
func (cli *Client) MGetAndScan(ctx context.Context, dst interface{}, keys ...string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		err := c.(driver.Cache).MGet(ctx, keys...).Scan(dst)
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
		val, err := c.(driver.Cache).Get(ctx, key).Result()
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
func (cli *Client) CheckAndScan(ctx context.Context, dst interface{}, key string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val := c.(driver.Cache).Get(ctx, key)
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
func (cli *Client) HSet(ctx context.Context, key string, values ...interface{}) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(cli.drivers) == 0 {
		return ErrNoCacheClient
	}

	for _, c := range cli.drivers {
		c.(driver.Cache).HSet(ctx, key, values)
	}

	return nil
}

// HVals 获取Hash表的所有值
func (cli *Client) HVals(ctx context.Context, key string) ([]string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		result, err := c.(driver.Cache).HVals(ctx, key).Result()
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
		result, err := c.(driver.Cache).HKeys(ctx, key).Result()
		if err == nil {
			return result, nil
		}
	}

	return nil, ErrNoRecord
}

// HKeysAndScan 获取Hash表的所有键并扫描到dst中
func (cli *Client) HKeysAndScan(ctx context.Context, dst interface{}, key string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val := c.(driver.Cache).HKeys(ctx, key)
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
		val := c.(driver.Cache).HLen(ctx, key)
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
		val := c.(driver.Cache).HGet(ctx, key, field)
		if val.Err() == nil {
			return val.Val(), nil
		}
	}

	return "", ErrNoRecord
}

// HGetAndScan 获取Hash表指定字段的值
func (cli *Client) HGetAndScan(ctx context.Context, dst interface{}, key, field string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val := c.(driver.Cache).HGet(ctx, key, field)
		if val.Err() == nil && val.Scan(dst) == nil {
			return nil
		}
	}
	return ErrNoRecord
}

// HMGet 获取Hash表指定字段的值
func (cli *Client) HMGet(ctx context.Context, key string, fields ...string) ([]interface{}, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		result, err := c.(driver.Cache).HMGet(ctx, key, fields...).Result()
		if err == nil {
			return result, nil
		}
	}

	return nil, ErrNoRecord
}

// HMGetAndScan 获取Hash表指定字段的值并扫描进入到dst中
func (cli *Client) HMGetAndScan(ctx context.Context, dst interface{}, key string, fields ...string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val := c.(driver.Cache).HMGet(ctx, key, fields...)
		if val.Err() == nil && val.Scan(dst) == nil {
			return nil
		}
	}

	return ErrNoRecord
}

// HValsAndScan 获取Hash表的所有值并扫如dst中
func (cli *Client) HValsAndScan(ctx context.Context, dst interface{}, key string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val := c.(driver.Cache).HVals(ctx, key)
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
		c.(driver.Cache).HDel(ctx, key, fields...)
	}
	return nil
}

// =======================================================
// ================= LIST ================================
// =======================================================

// LTrim 获取列表内的范围数据
func (cli *Client) LTrim(ctx context.Context, key string, start, stop int64) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(cli.drivers) == 0 {
		return ErrNoCacheClient
	}

	for _, c := range cli.drivers {
		c.(driver.Cache).LTrim(ctx, key, start, stop).Result()
	}

	return nil
}

// LPush 推送数据
func (cli *Client) LPush(ctx context.Context, key string, data ...interface{}) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(cli.drivers) == 0 {
		return ErrNoCacheClient
	}

	for _, c := range cli.drivers {
		c.(driver.Cache).LPush(ctx, key, data...)
	}

	return nil
}

// LRang 获取列表内的范围数据
func (cli *Client) LRang(ctx context.Context, key string, start, stop int64) ([]string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val, err := c.(driver.Cache).LRang(ctx, key, start, stop).Result()
		if err == nil {
			return val, nil
		}
	}
	return nil, ErrNoRecord
}

// LRangAndScan 通过扫描方式获取列表内的范围内数据
func (cli *Client) LRangAndScan(ctx context.Context, dst interface{}, key string, start, stop int64) error {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val := c.(driver.Cache).LRang(ctx, key, start, stop)
		if val.Err() == nil && val.ScanSlice(dst) == nil {
			return nil
		}
	}

	return ErrNoRecord
}

// LPop 取出列表内的第一个数据
func (cli *Client) LPop(ctx context.Context, key string) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, c := range cli.drivers {
		val, err := c.(driver.Cache).LPop(ctx, key).Result()
		if err == nil {
			return val, nil
		}
	}

	return "", ErrNoRecord
}

// LPopAndScan 通过扫描方式取出列表内的第一个数据
func (cli *Client) LPopAndScan(ctx context.Context, dst interface{}, key string) error {
	if ctx == nil {
		ctx = context.Background()
	}

	for _, c := range cli.drivers {
		val := c.(driver.Cache).LPop(ctx, key)
		if val.Err() == nil && val.Scan(dst) == nil {
			return nil
		}
	}

	return ErrNoRecord
}

// LLen 返回列表长度
func (cli *Client) LLen(ctx context.Context, key string) (int64, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	for _, c := range cli.drivers {
		val := c.(driver.Cache).LLen(ctx, key)
		if val.Err() == nil {
			return val.Result()
		}
	}

	return 0, nil
}

*/
