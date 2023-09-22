package driver

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"math/rand"
	"testing"
	"time"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/19 23:24
  @describe :
*/

func TestNewDistributeMemory(t *testing.T) {
	count := int64(2)
	mems := make([]Cache, 0)
	for i := 0; i < int(count); i++ {
		cfg := DistributeMemoryConfig{Port: 2000 + i, Prefix: "mydear", EtcdCfg: clientv3.Config{Endpoints: []string{"127.0.0.1:2379"}}}
		mem, err := NewDistributeMemory(cfg)
		if err != nil {
			log.Println(err)
		}
		mems = append(mems, mem)
	}

	randKey := func() string {
		return fmt.Sprintf("key-%d", rand.Int63n(20))
	}
	randField := func() string {
		return fmt.Sprintf("field-%d", rand.Int63n(20))
	}

	randValue := func() string {
		return fmt.Sprintf("value-%d", rand.Int63n(20))
	}

	for {
		mem := mems[rand.Int63n(count)]
		log.Println(mem.Del(context.Background(), randKey(), randKey(), randKey()))
		mem = mems[rand.Int63n(count)]
		log.Println(mem.Expire(context.Background(), randKey(), time.Second*5))

		mem = mems[rand.Int63n(count)]
		log.Println(mem.ExpireAt(context.Background(), randKey(), time.Now().Add(time.Minute)))
		mem = mems[rand.Int63n(count)]
		log.Println(mem.Persist(context.Background(), randKey()))
		mem = mems[rand.Int63n(count)]
		log.Println(mem.Set(context.Background(), randKey(), randValue(), time.Second*5))
		mem = mems[rand.Int63n(count)]
		log.Println(mem.SetNX(context.Background(), randKey(), randValue(), time.Second*5))
		mem = mems[rand.Int63n(count)]
		log.Println(mem.HDel(context.Background(), randKey(), randField(), randField(), randField()))
		mem = mems[rand.Int63n(count)]
		log.Println(mem.HSet(context.Background(), randKey(), randField(), randValue(), randField(), randValue(), randField(), randValue()))
		mem = mems[rand.Int63n(count)]
		log.Println(mem.HSetNX(context.Background(), randKey(), randField(), randValue()))
		mem = mems[rand.Int63n(count)]
		log.Println(mem.LPush(context.Background(), randKey(), randValue(), randValue(), randValue()))
		mem = mems[rand.Int63n(count)]
		log.Println(mem.LPop(context.Background(), randKey()))
		mem = mems[rand.Int63n(count)]
		log.Println(mem.LShift(context.Background(), randKey()))
		mem = mems[rand.Int63n(count)]
		log.Println(mem.LTrim(context.Background(), randKey(), rand.Int63n(10), rand.Int63n(10)))

		time.Sleep(time.Second)
	}
}
