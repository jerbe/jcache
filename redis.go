package jcache

import (
	"errors"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

// initRedis 初始化redis
func initRedis(cfg *RedisConfig) (redis.UniversalClient, error) {
	var cli redis.UniversalClient
	if len(cfg.Addrs) == 0 {
		return nil, errors.New("未设置Addrs")
	}

	var dialTO = time.Second * 5

	switch strings.ToLower(cfg.Mode) {
	case "sentinel": // 哨兵模式
		if cfg.MasterName == "" {
			return nil, errors.New("redis.mode 哨兵(sentinel)模式下,未指定master_name")
		}

		// 返回 *redis.FailoverClient
		cli = redis.NewUniversalClient(&redis.UniversalOptions{
			MasterName:  cfg.MasterName,
			Addrs:       cfg.Addrs,
			Username:    cfg.Username,
			Password:    cfg.Password,
			DialTimeout: dialTO,
		})
	case "cluster": //集群模式
		// 返回 *redis.ClusterClient
		cli = redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs:       cfg.Addrs,
			Username:    cfg.Username,
			Password:    cfg.Password,
			DialTimeout: dialTO,
		})
	default: // 单例模式
		if len(cfg.Addrs) > 1 {
			return nil, errors.New("redis.mode 单例(single)模式下,addrs只允许一个元素")
		}
		// 返回 *redis.Client
		cli = redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs:       cfg.Addrs[0:1],
			Username:    cfg.Username,
			Password:    cfg.Password,
			DialTimeout: dialTO,
		})
	}
	return cli, nil
}
