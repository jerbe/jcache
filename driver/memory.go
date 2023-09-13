package driver

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/8 15:58
  @describe :
*/

var MemoryNil = errors.New("memory cache: nil")

type MemoryDriverStoreType string

const (
	MemoryDriverStoreTypeString    MemoryDriverStoreType = "String"
	MemoryDriverStoreTypeHash      MemoryDriverStoreType = "Hash"
	MemoryDriverStoreTypeList      MemoryDriverStoreType = "List"
	MemoryDriverStoreTypeSet       MemoryDriverStoreType = "Set"
	MemoryDriverStoreTypeSortedSet MemoryDriverStoreType = "SortedSet"
)

func NewMemory() Cache {
	return &Memory{
		rwMu: sync.RWMutex{},
		ss:   newStringStore(),
		hs:   newHashStore(),
		ls:   newListStore(),
	}
}

func NewMemoryString() String {
	return NewMemory()
}

type Memory struct {
	rwMu sync.RWMutex

	ss *stringStore

	hs *hashStore

	ls *listStore
}

// checkKeyType 检测Key的类型是否正确
func (mc *Memory) checkKeyType(key string, useFor MemoryDriverStoreType) (bool, error) {
	mc.ss.rwMutex.RLock()
	defer mc.ss.rwMutex.RUnlock()
	_, ok := mc.ss.values[key]
	if ok {
		if useFor != MemoryDriverStoreTypeString {
			return false, fmt.Errorf("该Key已经是String类型,不可用设置成%s类型", useFor)
		}
		return true, nil
	}
	return false, nil
}

// ================================================================================================
// ====================================== COMMON ==================================================
// ================================================================================================

// Exists 判断某个Key是否存在
func (mc *Memory) Exists(ctx context.Context, keys ...string) IntValuer {
	mc.rwMu.RLock()
	defer mc.rwMu.RUnlock()

	result := &redis.IntCmd{}
	cnt := int64(0)
	// 遍历字符串存储
	_, err := mc.ss.Exists(ctx, keys...)
	if err == nil {
		cnt++
	}

	result.SetVal(cnt)
	// 遍历Hash存储

	// 遍历List存储
	return result
}

// Del 删除一个或多个key
func (mc *Memory) Del(ctx context.Context, keys ...string) IntValuer {
	result := &redis.IntCmd{}
	cnt := int64(0)
	// 遍历字符串存储
	delcnt, err := mc.ss.Del(ctx, keys...)
	if err == nil {
		cnt += delcnt
	}

	result.SetVal(cnt)

	// 遍历Hash存储

	// 遍历List存储
	return result
}

// Expire 设置某个key的存活时间
func (mc *Memory) Expire(ctx context.Context, key string, ttl time.Duration) BoolValuer {
	result := &redis.BoolCmd{}
	// 遍历字符串存储
	// Todo 错误收集
	b, err := mc.ss.Expire(ctx, key, ttl)
	if err == nil && b {
		result.SetVal(true)
		return result
	}

	// 遍历Hash存储

	// 遍历List存储
	return result
}

// ExpireAt 设置某个key在指定时间内到期
func (mc *Memory) ExpireAt(ctx context.Context, key string, at *time.Time) BoolValuer {
	result := &redis.BoolCmd{}
	// 遍历字符串存储
	// Todo 错误收集
	b, err := mc.ss.ExpireAt(ctx, key, *at)
	if err == nil && b {
		result.SetVal(true)
		return result
	}

	// 遍历Hash存储

	// 遍历List存储
	return result
}

// ================================================================================================
// ====================================== STRING ==================================================
// ================================================================================================

// Set 设置数据
func (mc *Memory) Set(ctx context.Context, key string, data any, expiration time.Duration) StatusValuer {
	val := new(redis.StatusCmd)
	// 检测该Key是否被其他类型用了
	if _, err := mc.checkKeyType(key, MemoryDriverStoreTypeString); err != nil {
		val.SetErr(err)
		return val
	}

	err := mc.ss.Set(ctx, key, data, expiration)
	if err == nil {
		val.SetVal("OK")
	}
	val.SetErr(err)
	return val
}

// SetNX 设置数据,如果key不存在的话
func (mc *Memory) SetNX(ctx context.Context, key string, data any, expiration time.Duration) BoolValuer {
	val := new(redis.BoolCmd)
	// 检测该Key是否被其他类型用了
	if _, err := mc.checkKeyType(key, MemoryDriverStoreTypeString); err != nil {
		val.SetErr(err)
		return val
	}

	b, err := mc.ss.SetNX(ctx, key, data, expiration)
	if err == nil {
		val.SetVal(b)
	}
	val.SetErr(err)
	return val
}

// Get 获取数据
func (mc *Memory) Get(ctx context.Context, key string) StringValuer {
	val := new(redis.StringCmd)
	// 检测该Key是否被其他类型用了
	if _, err := mc.checkKeyType(key, MemoryDriverStoreTypeString); err != nil {
		val.SetErr(err)
		return val
	}

	v, err := mc.ss.Get(ctx, key)
	if err == nil {
		val.SetVal(v)
	}
	val.SetErr(err)
	return val
}

// MGet 获取数据
func (mc *Memory) MGet(ctx context.Context, keys ...string) SliceValuer {
	val := new(redis.SliceCmd)
	anies, err := mc.ss.MGet(ctx, keys...)
	if err == nil {
		val.SetVal(anies)
	}
	val.SetErr(err)
	return val
}

// ================================================================================================
// ======================================== HASH ==================================================
// ================================================================================================

// HDel 哈希表删除指定字段(fields)
func (mc *Memory) HDel(ctx context.Context, key string, fields ...string) IntValuer {
	val := new(redis.IntCmd)
	cnt, err := mc.hs.HDel(ctx, key, fields...)
	if err == nil {
		val.SetVal(cnt)
	}
	val.SetErr(err)
	return val
}

// HSet 哈希表设置数据
func (mc *Memory) HSet(ctx context.Context, key string, data ...any) IntValuer {
	val := new(redis.IntCmd)
	cnt, err := mc.hs.HSet(ctx, key, data...)
	if err == nil {
		val.SetVal(cnt)
	}
	val.SetErr(err)
	return val
}

// HGet 哈希表获取一个数据
func (mc *Memory) HGet(ctx context.Context, key string, field string) StringValuer {
	val := new(redis.StringCmd)
	v, err := mc.hs.HGet(ctx, key, field)
	if err == nil {
		val.SetVal(v)
	}
	val.SetErr(err)
	return val
}

// HMGet 哈希表获取多个数据
func (mc *Memory) HMGet(ctx context.Context, key string, fields ...string) SliceValuer {
	val := new(redis.SliceCmd)
	v, err := mc.hs.HMGet(ctx, key, fields...)
	if err == nil {
		val.SetVal(v)
	}
	val.SetErr(err)
	return val
}

// HKeys 哈希表获取某个Key的所有字段(field)
func (mc *Memory) HKeys(ctx context.Context, key string) StringSliceValuer {
	val := new(redis.StringSliceCmd)
	v, err := mc.hs.HKeys(ctx, key)
	if err == nil {
		val.SetVal(v)
	}
	val.SetErr(err)
	return val
}

// HVals 哈希表获取所有值
func (mc *Memory) HVals(ctx context.Context, key string) StringSliceValuer {
	val := new(redis.StringSliceCmd)
	v, err := mc.hs.HVals(ctx, key)
	if err == nil {
		val.SetVal(v)
	}
	val.SetErr(err)
	return val
}

// HLen 哈希表所有字段的数量
func (mc *Memory) HLen(ctx context.Context, key string) IntValuer {
	val := new(redis.IntCmd)
	v, err := mc.hs.HLen(ctx, key)
	if err == nil {
		val.SetVal(v)
	}
	val.SetErr(err)
	return val
}

// ================================================================================================
// ======================================== LIST ==================================================
// ================================================================================================

// Trim 对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除。
// 举个例子，执行命令 LTRIM list 0 2 ，表示只保留列表 list 的前三个元素，其余元素全部删除。
// 下标(index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。
// 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
func (mc *Memory) Trim(ctx context.Context, key string, start, stop int64) StatusValuer {
	val := new(redis.StatusCmd)
	err := mc.ls.Trim(ctx, key, start, stop)
	if err == nil {
		val.SetVal("OK")
	}
	val.SetErr(err)
	return val
}

// Push 将数据推入到列表中
func (mc *Memory) Push(ctx context.Context, key string, data ...any) IntValuer {
	val := new(redis.IntCmd)
	v, err := mc.ls.Push(ctx, key, data...)
	if err == nil {
		val.SetVal(v)
	}
	val.SetErr(err)
	return val
}

// Rang 提取列表范围内的数据
func (mc *Memory) Rang(ctx context.Context, key string, start, stop int64) StringSliceValuer {
	val := new(redis.StringSliceCmd)
	v, err := mc.ls.Rang(ctx, key, start, stop)
	if err == nil {
		val.SetVal(v)
	}
	val.SetErr(err)
	return val

}

// Pop 推出列表尾的最后数据
func (mc *Memory) Pop(ctx context.Context, key string) StringValuer {
	val := new(redis.StringCmd)
	v, err := mc.ls.Pop(ctx, key)
	if err == nil {
		val.SetVal(v)
	}
	val.SetErr(err)
	return val
}

// Shift 推出列表头的第一个数据
func (mc *Memory) Shift(ctx context.Context, key string) StringValuer {
	val := new(redis.StringCmd)
	v, err := mc.ls.Shift(ctx, key)
	if err == nil {
		val.SetVal(v)
	}
	val.SetErr(err)
	return val
}
