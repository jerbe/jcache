package jcache

import (
	"context"
	"github.com/jerbe/jcache/memory"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"time"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/8/28 13:00
  @describe :
*/

var (
	BG_CTX = context.Background()

	redisCli redis.UniversalClient

	memoryCli *memory.Cache
)

func Redis() redis.UniversalClient {
	return redisCli
}

func Memory() *memory.Cache {
	return memoryCli
}

const (
	// KeepTTL 保持原先的过期时间(TTL)
	KeepTTL = -1

	// DefaultEmptySetNXDuration 默认空对象设置过期时效
	DefaultEmptySetNXDuration = time.Second * 10

	// DefaultExpirationDuration 默认缓存过期时效
	DefaultExpirationDuration = time.Hour
)

// RandomExpirationDuration 以 DefaultExpirationDuration 为基础,返回一个 DefaultExpirationDuration ± 30m 内的时间长度
func RandomExpirationDuration() time.Duration {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	r := (1 - rnd.Int63n(2)) * rnd.Int63n(30)
	return DefaultExpirationDuration + (time.Minute * time.Duration(r))
}

// Init 初始化缓存
func Init(cfg *Config) error {
	var err error
	if cfg.Redis != nil {
		redisCli, err = initRedis(cfg.Redis)
		if err != nil {
			return err
		}
	}

	if cfg.Memory != nil {
		memoryCli = new(memory.Cache)
	}
	return nil
}
