package jcache

import (
	"context"
	"github.com/jerbe/jcache/v2/driver"
	clientv3 "go.etcd.io/etcd/client/v3"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/13 17:05
  @describe :
*/

func newStringClient() *StringClient {

	/*
		rdCfg := &driver.RedisConfig{
			Mode:       "single",
			MasterName: "",
			Addrs:      []string{"192.168.31.101:6379"},
			Database:   "",
			Username:   "",
			Password:   "root",
		}
		redisDriver := driver.NewRedisOptionsWithConfig(rdCfg)
	*/
	var cnt = 2
	l := make([]driver.Cache, 0)
	for i := 0; i < cnt; i++ {
		cfg := driver.MemoryConfig{Port: 9890 + i, Prefix: "/kings", EtcdConfig: clientv3.Config{Endpoints: []string{"127.0.0.1:2379"}}}
		memoryDriver, err := driver.NewMemoryWithConfig(cfg)
		if err != nil {
			panic(err)
		}
		l = append(l, memoryDriver)
	}

	return NewStringClient(l[0])
}

func BenchmarkStringClient(b *testing.B) {
	cli := newStringClient()

	rand.Seed(time.Now().UnixNano())

	b.SetParallelism(1000)
	b.RunParallel(func(pb *testing.PB) {
		key := strconv.FormatInt(rand.Int63(), 10)
		val := strconv.FormatInt(rand.Int63(), 10)

		for pb.Next() {
			err := cli.Set(context.Background(), key, val, time.Hour)
			if err != nil {
				b.Logf("Set:发生错误:%+v", err)
			}
			//_, err = cli.Get(context.Background(), key)
			//if err != nil {
			//	b.Logf("Get:发生错误:%+v", err)
			//}
			//
			//err = cli.Expire(context.Background(), key, time.Minute)
			//if err != nil {
			//	b.Logf("Expire:发生错误:%+v", err)
			//}
			//
			//err = cli.ExpireAt(context.Background(), key, time.Now().Add(time.Minute))
			//if err != nil {
			//	b.Logf("ExpireAt:发生错误:%+v", err)
			//}
			//err = cli.Del(context.Background(), key)
			//if err != nil {
			//	b.Logf("Del:发生错误:%+v", err)
			//}

			//b.Log(v)
		}
	})
}
