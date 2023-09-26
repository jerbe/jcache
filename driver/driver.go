package driver

import (
	"context"
	"time"

	"github.com/jerbe/jcache/v2/errors"

	cerr "github.com/jerbe/go-errors"
	"github.com/redis/go-redis/v9"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/12 23:15
  @describe :
*/

// Cache 缓存器
type Cache interface {
	Common
	String
	Hash
	List
}

// String 字符串
type String interface {
	Common

	// Set 设置数据
	Set(ctx context.Context, key string, data interface{}, ttl time.Duration) StatusValuer

	// SetNX 如果key不存在才设置数据
	SetNX(ctx context.Context, key string, data interface{}, ttl time.Duration) BoolValuer

	// Get 获取数据
	Get(ctx context.Context, key string) StringValuer

	// MGet 获取多个key的数据
	MGet(ctx context.Context, keys ...string) SliceValuer
}

// Hash 哈希表
type Hash interface {
	Common

	// HExists 判断field是否存在
	HExists(ctx context.Context, key, field string) BoolValuer

	// HDel 哈希表删除指定字段(fields)
	HDel(ctx context.Context, key string, fields ...string) IntValuer

	// HSet 哈希表设置数据
	HSet(ctx context.Context, key string, data ...interface{}) IntValuer

	// HSetNX 设置哈希表field对应的值,当field不存在时才能成功
	HSetNX(ctx context.Context, key, field string, data interface{}) BoolValuer

	// HGet 哈希表获取一个数据
	HGet(ctx context.Context, key string, field string) StringValuer

	// HMGet 哈希表获取多个数据
	HMGet(ctx context.Context, key string, fields ...string) SliceValuer

	// HKeys 哈希表获取某个Key的所有字段(field)
	HKeys(ctx context.Context, key string) StringSliceValuer

	// HVals 哈希表获取所有值
	HVals(ctx context.Context, key string) StringSliceValuer

	// HGetAll 获取哈希表的键跟值
	HGetAll(ctx context.Context, key string) MapStringStringValuer

	// HLen 哈希表所有字段的数量
	HLen(ctx context.Context, key string) IntValuer
}

// List 列表
type List interface {
	Common

	// LTrim 对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除。
	//举个例子，执行命令 LTRIM list 0 2 ，表示只保留列表 list 的前三个元素，其余元素全部删除。
	//下标(index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。
	//你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
	LTrim(ctx context.Context, key string, start, stop int64) StatusValuer

	// LPush 将数据推入到列表中
	LPush(ctx context.Context, key string, data ...interface{}) IntValuer

	// LRang 提取列表范围内的数据
	LRang(ctx context.Context, key string, start, stop int64) StringSliceValuer

	// LPop 推出列表尾的最后数据
	LPop(ctx context.Context, key string) StringValuer

	// LShift 推出列表头的第一个数据
	LShift(ctx context.Context, key string) StringValuer

	// LLen 获取list列表的长度
	LLen(ctx context.Context, key string) IntValuer
}

// Common 通用接口
type Common interface {
	// Del 删除一个或多个key
	Del(ctx context.Context, keys ...string) IntValuer

	// Exists 判断某个Key是否存在
	Exists(ctx context.Context, keys ...string) IntValuer

	// Expire 设置某个key的存活时间
	Expire(ctx context.Context, key string, ttl time.Duration) BoolValuer

	// ExpireAt 设置某个key在指定时间内到期
	ExpireAt(ctx context.Context, key string, at time.Time) BoolValuer

	// Persist 将某个Key设置成持久性
	Persist(ctx context.Context, key string) BoolValuer
}

// ================================================================================================
// =================================== VALUER =====================================================
// ================================================================================================

// StatusValuer 状态数值接口
type StatusValuer interface {
	Val() string
	Err() error

	Result() (string, error)
}

// StringValuer 字符串数值接口
type StringValuer interface {
	Val() string
	Err() error
	Scan(dst interface{}) error

	Bytes() ([]byte, error)
	Bool() (bool, error)
	Int() (int, error)
	Int64() (int64, error)
	Uint64() (uint64, error)
	Float32() (float32, error)
	Float64() (float64, error)
	Time() (time.Time, error)

	Result() (string, error)
}

// StringSliceValuer 字符串切片数值接口
type StringSliceValuer interface {
	Val() []string
	Err() error
	ScanSlice(container interface{}) error

	Result() ([]string, error)
}

// MapStringStringValuer 以string为键string为值的map
type MapStringStringValuer interface {
	Val() map[string]string
	Err() error
	Scan(dst interface{}) error

	Result() (map[string]string, error)
}

// SliceValuer 切片数值接口
type SliceValuer interface {
	Val() []interface{}
	Err() error
	Scan(dst interface{}) error

	Result() ([]interface{}, error)
}

// IntValuer 整形数值接口
type IntValuer interface {
	Val() int64
	Err() error

	Result() (int64, error)
}

// BoolValuer 布尔数值接口
type BoolValuer interface {
	Val() bool
	Err() error

	Result() (bool, error)
}

// FloatValuer 浮点型数值接口
type FloatValuer interface {
	Val() float64
	Err() error

	Result() (float64, error)
}

func translateErr(err error) error {
	if cerr.IsIn(err, redis.Nil, MemoryNil) {
		return errors.Nil
	}
	return err
}
