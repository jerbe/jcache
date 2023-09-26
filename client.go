package jcache

import (
	"context"
	"time"

	"github.com/jerbe/jcache/v2/driver"
	"github.com/jerbe/jcache/v2/errors"

	jerrors "github.com/jerbe/go-errors"
	"github.com/redis/go-redis/v9"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/8/27 18:08
  @describe :
*/

var (
	Nil              = errors.Nil
	ErrNoCacheClient = errors.ErrNoCacheClient
)

// Z 表示已排序的集合成员。
type Z struct {
	Score  float64
	Member interface{}
}

// returnable 检测值是否可以返回
func returnable(val errors.ErrorValuer) bool {
	return val.Err() == nil || !jerrors.IsIn(val.Err(), redis.Nil, driver.MemoryNil)
}

func (cli *BaseClient) preCheck(ctx context.Context) (context.Context, context.CancelFunc) {
	if len(cli.drivers) == 0 {
		panic(ErrNoCacheClient)
	}

	var cancelFunc context.CancelFunc
	if ctx == nil {
		ctx, cancelFunc = context.WithTimeout(context.Background(), time.Second*5)
	} else {
		ctx, cancelFunc = context.WithCancel(context.Background())
	}

	return ctx, cancelFunc
}

// =======================================================
// ================= BaseClient ==========================
// =======================================================

type BaseClient struct {
	drivers []driver.Common
}

// Exists 判断某个Key是否存在
func (cli *BaseClient) Exists(ctx context.Context, keys ...string) driver.IntValuer {
	ctx, _ = cli.preCheck(ctx)

	var value driver.IntValuer
	for _, c := range cli.drivers {
		if value = c.Exists(ctx, keys...); returnable(value) {
			return value
		}
	}
	return value
}

// Del 删除键
func (cli *BaseClient) Del(ctx context.Context, keys ...string) driver.IntValuer {
	ctx, _ = cli.preCheck(ctx)

	var value driver.IntValuer
	for i, c := range cli.drivers {
		if v := c.Del(ctx, keys...); i == 0 {
			value = v
		}
	}
	return value
}

// Expire 设置某个Key的TTL时长
func (cli *BaseClient) Expire(ctx context.Context, key string, expiration time.Duration) driver.BoolValuer {
	ctx, _ = cli.preCheck(ctx)
	var value driver.BoolValuer
	for i, c := range cli.drivers {
		if v := c.Expire(ctx, key, expiration); i == 0 {
			value = v
		}
	}

	return value
}

// ExpireAt 设置某个key在指定时间内到期
func (cli *BaseClient) ExpireAt(ctx context.Context, key string, at time.Time) driver.BoolValuer {
	ctx, _ = cli.preCheck(ctx)

	var value driver.BoolValuer
	for i, c := range cli.drivers {
		if v := c.ExpireAt(ctx, key, at); i == 0 {
			value = v
		}
	}

	return value
}

// =======================================================
// ================= Client ==============================
// =======================================================

type Client struct {
	BaseClient
	StringClient
	HashClient
	ListClient
}

// NewClient 实例化出一个客户端
func NewClient(drivers ...driver.Cache) *Client {
	drs := make([]driver.Common, len(drivers))

	for i := 0; i < len(drivers); i++ {
		drs[i] = drivers[i]
	}

	if len(drs) == 0 {
		drs = append(drs, driver.NewMemory())
	}

	cli := BaseClient{drivers: drs}

	return &Client{
		BaseClient:   cli,
		StringClient: StringClient{BaseClient: cli},
		HashClient:   HashClient{BaseClient: cli},
		ListClient:   ListClient{BaseClient: cli},
	}
}
