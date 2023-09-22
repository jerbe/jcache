package driver

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	// Mode 模式
	// 支持:single,sentinel,cluster
	Mode       string   `yaml:"mode"`
	MasterName string   `yaml:"master_name"`
	Addrs      []string `yaml:"addrs"`
	Database   string   `yaml:"database"`
	Username   string   `yaml:"username"`
	Password   string   `yaml:"password"`
}

// ========================================================
// ====================== Redis 实例 =======================
// ========================================================

// Redis 驱动器
type Redis struct {
	cli redis.UniversalClient
}

type RedisOptions struct {
	Config *RedisConfig
	Client redis.UniversalClient
}

func NewRedisOptionsWithConfig(cfg *RedisConfig) *RedisOptions {
	return &RedisOptions{Config: cfg}
}

func NewRedisOptionsWithClient(cli redis.UniversalClient) *RedisOptions {
	return &RedisOptions{Client: cli}
}

var _ Cache = new(Redis)

func NewRedis(opt *RedisOptions) Cache {
	r := &Redis{}

	if opt.Client != nil {
		r.cli = opt.Client
		return r
	}

	if cfg := opt.Config; cfg != nil {
		var cli redis.UniversalClient
		var dialTimeout = time.Second * 5
		switch strings.ToLower(cfg.Mode) {
		case "sentinel": // 哨兵模式
			// 返回 *redis.FailoverClient
			cli = redis.NewUniversalClient(&redis.UniversalOptions{
				MasterName:  cfg.MasterName,
				Addrs:       cfg.Addrs,
				Username:    cfg.Username,
				Password:    cfg.Password,
				DialTimeout: dialTimeout,
			})
		case "cluster": //集群模式
			// 返回 *redis.ClusterClient
			cli = redis.NewUniversalClient(&redis.UniversalOptions{
				Addrs:       cfg.Addrs,
				Username:    cfg.Username,
				Password:    cfg.Password,
				DialTimeout: dialTimeout,
			})
		default: // 单例模式
			// 返回 *redis.Client
			cli = redis.NewUniversalClient(&redis.UniversalOptions{
				Addrs:       cfg.Addrs[0:1],
				Username:    cfg.Username,
				Password:    cfg.Password,
				DialTimeout: dialTimeout,
			})
		}
		r.cli = cli
		return r
	}

	panic(errors.New("redis driver: invalid options 'Client' and 'Config' is nil"))
}

func NewRedisString(opt *RedisOptions) String {
	return NewRedis(opt)
}

func NewRedisHashDriver(opt *RedisOptions) Hash {
	return NewRedis(opt)
}

func NewRedisListDriver(opt *RedisOptions) Hash {
	return NewRedis(opt)
}

// ============================
// ========= Common ===========
// ============================

// Del 删除一个或多个key
func (r *Redis) Del(ctx context.Context, keys ...string) IntValuer {
	return r.cli.Del(ctx, keys...)
}

// Exists 判断某个Key是否存在
func (r *Redis) Exists(ctx context.Context, keys ...string) IntValuer {
	return r.cli.Exists(ctx, keys...)
}

// Expire 设置某个key的存活时间
func (r *Redis) Expire(ctx context.Context, key string, ttl time.Duration) BoolValuer {
	return r.cli.Expire(ctx, key, ttl)
}

// ExpireAt 设置某个key在指定时间内到期
func (r *Redis) ExpireAt(ctx context.Context, key string, at time.Time) BoolValuer {
	return r.cli.ExpireAt(ctx, key, at)
}

// Persist 移除某个key的TTL,设置成持久性
func (r *Redis) Persist(ctx context.Context, key string) BoolValuer {
	return r.cli.Persist(ctx, key)
}

// ============================
// ========= String ===========
// ============================

// Set 设置数据
func (r *Redis) Set(ctx context.Context, key string, data interface{}, ttl time.Duration) StatusValuer {
	return r.cli.Set(ctx, key, data, ttl)
}

// SetNX 如果key不存在才设置数据
func (r *Redis) SetNX(ctx context.Context, key string, data interface{}, ttl time.Duration) BoolValuer {
	return r.cli.SetNX(ctx, key, data, ttl)
}

// Get 获取数据
func (r *Redis) Get(ctx context.Context, key string) StringValuer {
	return r.cli.Get(ctx, key)
}

// MGet 获取多个key的数据
func (r *Redis) MGet(ctx context.Context, keys ...string) SliceValuer {
	return r.cli.MGet(ctx, keys...)
}

// ============================
// ========== Hash ============
// ============================

// HExists 判断哈希表的field是否存在
func (r *Redis) HExists(ctx context.Context, key, field string) BoolValuer {
	return r.cli.HExists(ctx, key, field)
}

// HDel 哈希表删除指定字段(fields)
func (r *Redis) HDel(ctx context.Context, key string, fields ...string) IntValuer {
	return r.cli.HDel(ctx, key, fields...)
}

// HSet 哈希表设置数据
func (r *Redis) HSet(ctx context.Context, key string, data ...interface{}) IntValuer {
	return r.cli.HSet(ctx, key, data...)
}

// HSetNX 设置哈希表field对应的值,当field不存在时才能成功
func (r *Redis) HSetNX(ctx context.Context, key, field string, data interface{}) BoolValuer {
	return r.cli.HSetNX(ctx, key, field, data)
}

// HGet 哈希表获取一个数据
func (r *Redis) HGet(ctx context.Context, key string, field string) StringValuer {
	return r.cli.HGet(ctx, key, field)
}

// HMGet 哈希表获取多个数据
func (r *Redis) HMGet(ctx context.Context, key string, fields ...string) SliceValuer {
	return r.cli.HMGet(ctx, key, fields...)
}

// HKeys 哈希表获取某个Key的所有字段(field)
func (r *Redis) HKeys(ctx context.Context, key string) StringSliceValuer {
	return r.cli.HKeys(ctx, key)
}

// HVals 哈希表获取所有值
func (r *Redis) HVals(ctx context.Context, key string) StringSliceValuer {
	return r.cli.HVals(ctx, key)
}

// HGetAll 哈希表获取所有值,包括field跟value
func (r *Redis) HGetAll(ctx context.Context, key string) MapStringStringValuer {
	return r.cli.HGetAll(ctx, key)
}

// HLen 哈希表所有字段的数量
func (r *Redis) HLen(ctx context.Context, key string) IntValuer {
	return r.cli.HLen(ctx, key)
}

// ============================
// ========== List ============
// ============================

// LTrim 对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除。
// 举个例子，执行命令 LTRIM list 0 2 ，表示只保留列表 list 的前三个元素，其余元素全部删除。
// 下标(index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。
// 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
func (r *Redis) LTrim(ctx context.Context, key string, start, stop int64) StatusValuer {
	return r.cli.LTrim(ctx, key, start, stop)
}

// LPush 将数据推入到列表中
func (r *Redis) LPush(ctx context.Context, key string, data ...interface{}) IntValuer {
	return r.cli.LPush(ctx, key, data)
}

// LRang 提取列表范围内的数据
func (r *Redis) LRang(ctx context.Context, key string, start, stop int64) StringSliceValuer {
	return r.cli.LRange(ctx, key, start, stop)
}

// LPop 推出列表尾的最后数据
func (r *Redis) LPop(ctx context.Context, key string) StringValuer {
	return r.cli.RPop(ctx, key)
}

// LShift 推出列表头的第一个数据
func (r *Redis) LShift(ctx context.Context, key string) StringValuer {
	return r.cli.LPop(ctx, key)
}

// LLen 获取列表的长度
func (r *Redis) LLen(ctx context.Context, key string) IntValuer {
	return r.cli.LLen(ctx, key)
}
