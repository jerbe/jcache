package jcache

import (
	"context"
	"time"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/8/27 18:08
  @describe :
*/

// Exists 判断某个Key是否存在
func Exists(ctx context.Context, key string) (bool, error) {
	if redisCli != nil {
		cmd := redisCli.Exists(ctx, key)
		val, err := cmd.Result()

		if err == nil {
			return val > 0, nil
		}
	}

	return false, ErrNoCacheClient
}

// Set 设置数据
func Set(ctx context.Context, key string, data any, expiration time.Duration) error {
	if redisCli != nil {
		if rs := redisCli.Set(globCtx, key, data, expiration); rs.Err() != nil {
			// @TODO 写入队列进行重试
		}
	}

	return nil
}

// SetNX 设置数据,如果key不存在的话
func SetNX(ctx context.Context, key string, data any, expiration time.Duration) error {
	if redisCli != nil {
		if rs := redisCli.SetNX(globCtx, key, data, expiration); rs.Err() != nil {
			// @TODO 写入队列进行重试
		}
	}
	return nil
}

// HSet 写入hash数据
func HSet(ctx context.Context, key string, values ...any) error {
	if redisCli != nil {
		rs := redisCli.HSet(globCtx, key, values...)
		if rs.Err() != nil {
			// @TODO 写入队列进行重试
		}
	}

	return nil
}

// HVals 获取Hash表的所有值
func HVals(ctx context.Context, key string) ([]string, error) {
	if redisCli != nil {
		cmd := redisCli.HVals(globCtx, key)
		if cmd.Err() == nil {
			return cmd.Result()
		}
	}

	return nil, ErrNoCacheClient
}

// HKeys 获取Hash表的所有键
func HKeys(ctx context.Context, key string) ([]string, error) {
	if redisCli != nil {
		cmd := redisCli.HKeys(globCtx, key)

		if cmd.Err() == nil {
			return cmd.Result()
		}
	}

	return nil, ErrNoCacheClient
}

// HKeysAndScan 获取Hash表的所有键并扫描到dst中
func HKeysAndScan(ctx context.Context, key string, dst any) error {
	if redisCli != nil {
		cmd := redisCli.HKeys(globCtx, key)

		if cmd.Err() == nil {
			err := cmd.ScanSlice(dst)
			if err == nil {
				return nil
			}
		}
	}

	return ErrNoCacheClient
}

// HLen 获取Hash表的所有键个数
func HLen(ctx context.Context, key string) (int64, error) {
	if redisCli != nil {
		cmd := redisCli.HLen(globCtx, key)
		if cmd.Err() == nil {
			return cmd.Result()
		}
	}

	return 0, ErrNoCacheClient
}

// HGet 获取Hash表指定字段的值
func HGet(ctx context.Context, key, field string) (string, error) {
	if redisCli != nil {
		cmd := redisCli.HGet(globCtx, key, field)

		if cmd.Err() == nil {
			return cmd.Val(), nil
		}
	}

	return "", ErrNoCacheClient
}

// HGetAndScan 获取Hash表指定字段的值
func HGetAndScan(ctx context.Context, dst any, key, field string) error {
	if redisCli != nil {
		cmd := redisCli.HGet(globCtx, key, field)

		if cmd.Err() == nil {
			err := cmd.Scan(dst)
			if err == nil {
				return nil
			}
		}
	}

	return ErrNoCacheClient
}

// HMGet 获取Hash表指定字段的值
func HMGet(ctx context.Context, key string, fields ...string) ([]any, error) {
	if redisCli != nil {
		cmd := redisCli.HMGet(globCtx, key, fields...)
		if cmd.Err() == nil {
			return cmd.Val(), nil
		}
	}

	return nil, ErrNoCacheClient
}

// HMGetAndScan 获取Hash表指定字段的值并扫描进入到dst中
func HMGetAndScan(ctx context.Context, dst any, key string, fields ...string) error {
	if redisCli != nil {
		cmd := redisCli.HMGet(globCtx, key, fields...)

		if cmd.Err() == nil {
			err := cmd.Scan(dst)
			if err == nil {
				return nil
			}
		}
	}

	return ErrNoCacheClient
}

// HVals 获取Hash表的所有值
func HValsScan(ctx context.Context, dst any, key string) error {
	if redisCli != nil {
		cmd := redisCli.HVals(globCtx, key)

		if cmd.Err() == nil {
			err := cmd.ScanSlice(dst)
			if err == nil {
				return nil
			}
		}
	}

	return ErrNoCacheClient
}

// HDel 删除hash数据
func HDel(ctx context.Context, key string, fields ...string) error {
	if redisCli != nil {
		rs := redisCli.HDel(globCtx, key, fields...)
		if rs.Err() != nil {
			// @TODO 写入队列进行重试
		}
	}

	return nil
}

// Get 获取数据
func Get(ctx context.Context, key string) (string, error) {
	if redisCli != nil {
		cmd := redisCli.Get(globCtx, key)
		if cmd.Err() == nil {
			return cmd.Result()
		}
	}

	return "", ErrNoCacheClient
}

// MGet 获取多个Keys的值
func MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	if redisCli != nil {
		cmd := redisCli.MGet(globCtx, keys...)
		if cmd.Err() == nil {
			return cmd.Result()
		}
	}

	return nil, ErrNoCacheClient
}

// GetAndScan 获取并扫描
func GetAndScan(ctx context.Context, dst any, key string) error {
	if redisCli != nil {
		cmd := redisCli.Get(globCtx, key)

		if cmd.Err() == nil {
			err := cmd.Scan(dst)
			if err == nil {
				return nil
			}
		}
	}

	return ErrNoCacheClient
}

// MGetAndScan 获取多个Keys的值并扫描进dst中
func MGetAndScan(ctx context.Context, dst any, keys ...string) error {
	if redisCli != nil {
		cmd := redisCli.MGet(globCtx, keys...)

		if cmd.Err() == nil {
			err := cmd.Scan(dst)
			if err == nil {
				return nil
			}
		}
	}

	return ErrNoCacheClient
}

// CheckAndGet 检测并获取数据
func CheckAndGet(ctx context.Context, key string) (string, error) {
	if redisCli != nil {
		cmd := redisCli.Get(globCtx, key)
		var err = cmd.Err()

		if err == nil {
			// 如果找得到数据,并且没有错误并且数据为空,则返回找不到数据
			if cmd.Val() == "" {
				return "", ErrEmpty
			}
			return cmd.Val(), nil
		}
	}
	return "", ErrNoCacheClient
}

// CheckAndScan 获取数据
func CheckAndScan(ctx context.Context, dst any, key string) error {
	if redisCli != nil {
		cmd := redisCli.Get(globCtx, key)
		var err = cmd.Err()

		if err == nil {
			// 如果找得到数据,并且没有错误并且数据为空,则返回找不到数据
			if cmd.Val() == "" {
				return ErrEmpty
			}

			// 将不为空的数据扫描进struct中去
			if err = cmd.Scan(dst); err == nil {
				return nil
			}
		}
	}
	return ErrNoCacheClient
}

// Del 删除键
func Del(ctx context.Context, keys ...string) error {
	// 删除Redis
	if redisCli != nil {
		redisCli.Del(globCtx, keys...)
	}

	// 删除其他的
	return nil
}

// Push 推送数据
func Push(ctx context.Context, key string, data ...any) error {
	if redisCli != nil {
		// 从右边插入
		cmd := redisCli.RPush(globCtx, key, data...)
		if cmd.Err() != nil {
			// @TODO 重做?或提送到队列服务再重做?
		}
	}

	return nil
}

// Rang 获取列表内的范围数据
func Rang(ctx context.Context, key string, limit int64) ([]string, error) {
	if redisCli != nil {
		cmd := redisCli.LRange(globCtx, key, 0, limit-1)
		if cmd.Err() == nil {
			return cmd.Val(), nil
		}
	}

	return nil, ErrNoCacheClient
}

// RangAndScan 通过扫描方式获取列表内的范围内数据
func RangAndScan(ctx context.Context, dst any, key string, limit int64) error {
	if redisCli != nil {
		cmd := redisCli.LRange(globCtx, key, 0, limit-1)

		if cmd.Err() == nil {
			err := cmd.ScanSlice(dst)
			if err == nil {
				return nil
			}
		}
	}

	return ErrNoCacheClient
}

// Pop 取出列表内的第一个数据
func Pop(ctx context.Context, key string) (string, error) {
	if redisCli != nil {
		cmd := redisCli.LPop(globCtx, key)
		if cmd.Err() == nil {
			return cmd.Val(), nil
		}
	}

	return "", ErrNoCacheClient
}

// PopAndScan 通过扫描方式取出列表内的第一个数据
func PopAndScan(ctx context.Context, dst any, key string) error {
	if redisCli != nil {
		cmd := redisCli.LPop(globCtx, key)
		if cmd.Err() == nil {

			if err := cmd.Scan(dst); err == nil {
				return nil
			}
		}
	}
	return ErrNoCacheClient
}

// Expire 设置某个Key的TTL时长
func Expire(ctx context.Context, key string, expiration time.Duration) error {
	if redisCli != nil {
		redisCli.Expire(globCtx, key, expiration)
	}
	// @TODO 其他缓存方法
	return nil
}
