package jcache

import (
	"context"
	"time"

	"github.com/jerbe/jcache/v2/driver"
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
		c.ExpireAt(ctx, key, at)
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
	drs := make([]driver.Common, len(drivers))

	for i := 0; i < len(drivers); i++ {
		drs[i] = drivers[i]
	}

	if len(drs) == 0 {
		drs = append(drs, driver.NewMemory())
	}

	cli := baseClient{drivers: drs}

	return &Client{
		baseClient:   cli,
		StringClient: StringClient{baseClient: cli},
		HashClient:   HashClient{baseClient: cli},
		ListClient:   ListClient{baseClient: cli},
	}
}
