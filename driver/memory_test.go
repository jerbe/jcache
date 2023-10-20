package driver

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/19 23:24
  @describe :
*/

func TestNewDistributeMemory(t *testing.T) {
	count := int64(1)
	mems := make([]Cache, 0)
	for i := 0; i < int(count); i++ {
		/*
			cfg := DistributeMemoryConfig{Port: 2000 + i, Prefix: "mydear", Username: "root", Password: "root", EtcdCfg: clientv3.Config{Endpoints: []string{"192.168.31.101:2379"}}}
			mem, err := NewDistributeMemory(cfg)

			if err != nil {
				log.Fatal(err)
			}
		*/
		mem := NewMemory()

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

	fns := make([]func(cache Cache), 0)

	fns = append(fns, func(cache Cache) {
		start := time.Now()
		val := cache.Del(context.Background(), randKey(), randKey(), randKey())
		log.Println(time.Now().Sub(start), "cache.Del(context.Background(), randKey(), randKey(), randKey())", val)
	})

	fns = append(fns, func(cache Cache) {
		start := time.Now()
		val := cache.Expire(context.Background(), randKey(), time.Second*5)
		log.Println(time.Now().Sub(start), "cache.Expire(context.Background(), randKey(), time.Second*5)", val)
	})

	fns = append(fns, func(cache Cache) {
		start := time.Now()
		val := cache.ExpireAt(context.Background(), randKey(), time.Now().Add(time.Minute))
		log.Println(time.Now().Sub(start), "cache.ExpireAt(context.Background(), randKey(), time.Now().Add(time.Minute))", val)
	})

	fns = append(fns, func(cache Cache) {
		start := time.Now()
		val := cache.Persist(context.Background(), randKey())
		log.Println(time.Now().Sub(start), "cache.Persist(context.Background(), randKey())", val)
	})

	fns = append(fns, func(cache Cache) {
		start := time.Now()
		val := cache.Set(context.Background(), randKey(), randValue(), time.Second*5)
		log.Println(time.Now().Sub(start), "cache.Set(context.Background(), randKey(), randValue(), time.Second*5)", val)
	})

	fns = append(fns, func(cache Cache) {
		start := time.Now()
		val := cache.SetNX(context.Background(), randKey(), randValue(), time.Second*5)
		log.Println(time.Now().Sub(start), "cache.SetNX(context.Background(), randKey(), randValue(), time.Second*5)", val)
	})

	fns = append(fns, func(cache Cache) {
		start := time.Now()
		val := cache.HDel(context.Background(), randKey(), randField(), randField(), randField())
		log.Println(time.Now().Sub(start), "cache.HDel(context.Background(), randKey(), randField(), randField(), randField())", val)
	})
	fns = append(fns, func(cache Cache) {
		start := time.Now()
		val := cache.HSet(context.Background(), randKey(), randField(), randValue(), randField(), randValue(), randField(), randValue())
		log.Println(time.Now().Sub(start), "cache.HSet(context.Background(), randKey(), randField(), randValue(), randField(), randValue(), randField(), randValue())", val)
	})
	fns = append(fns, func(cache Cache) {
		start := time.Now()
		val := cache.HSetNX(context.Background(), randKey(), randField(), randValue())
		log.Println(time.Now().Sub(start), "cache.HSetNX(context.Background(), randKey(), randField(), randValue())", val)
	})
	fns = append(fns, func(cache Cache) {
		start := time.Now()
		val := cache.LPush(context.Background(), randKey(), randValue(), randValue(), randValue())

		log.Println(time.Now().Sub(start), "cache.LPush(context.Background(), randKey(), randValue(), randValue(), randValue())", val)
	})
	fns = append(fns, func(cache Cache) {
		start := time.Now()
		val := cache.LPop(context.Background(), randKey())
		log.Println(time.Now().Sub(start), "cache.LPop(context.Background(), randKey())", val)
	})
	fns = append(fns, func(cache Cache) {
		start := time.Now()
		val := cache.LShift(context.Background(), randKey())
		log.Println(time.Now().Sub(start), "cache.LShift(context.Background(), randKey())", val)
	})
	fns = append(fns, func(cache Cache) {
		start := time.Now()
		val := cache.LTrim(context.Background(), randKey(), rand.Int63n(10), rand.Int63n(10))
		log.Println(time.Now().Sub(start), "cache.LTrim(context.Background(), randKey(), rand.Int63n(10), rand.Int63n(10))", val)
	})
	fns = append(fns, func(cache Cache) {
		start := time.Now()
		val := cache.LBPop(context.Background(), time.Second*10, randKey())
		log.Println(time.Now().Sub(start), "cache.LBPop(context.Background(), time.Second*10, randKey())", val)
	})

	sig := make(chan os.Signal)
	//for i := 0; i < 10; i++ {
	go func() {
		ticker := time.NewTicker(time.Second)

		for {
			select {
			case <-ticker.C:
				log.Println("ticking...")
				//cache := mems[rand.Intn(len(mems))]
				//go func() {
				//	log.Println("cache.LBPop(context.Background(), time.Second*5, randKey())", cache.LBPop(context.Background(), time.Second*10, randKey()))
				//}()
				//
				//log.Println("cache.LTrim(context.Background(), randKey(), rand.Int63n(10), rand.Int63n(10))", cache.LTrim(context.Background(), randKey(), rand.Int63n(10), rand.Int63n(10)))

				for y := 0; y < len(fns); y++ {
					fn := fns[rand.Intn(len(fns))]
					cache := mems[rand.Intn(len(mems))]
					go fn(cache)
				}
			}
		}

		//for {
		//
		//	log.Println("time.Sleep(time.Second):", x)
		//	time.Sleep(time.Second)
		//}
	}()
	//}

	signal.Notify(sig, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT)
	s := <-sig
	log.Println("程序退出.原因:", s)

}

func TestLBPop(t *testing.T) {
	count := int64(2)
	mems := make([]Cache, 0)
	for i := 0; i < int(count); i++ {

		cfg := DistributeMemoryConfig{Port: 2000 + i, Prefix: "mydear", Username: "root", Password: "root", EtcdCfg: clientv3.Config{Endpoints: []string{"192.168.31.101:2379"}}}
		mem, err := NewDistributeMemory(cfg)
		if err != nil {
			log.Fatal(err)
		}

		//mem := NewMemory()

		mems = append(mems, mem)
	}

	randCache := func() Cache {
		return mems[rand.Intn(len(mems))]
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

	type Func func(cache Cache, id int, key, field, value string)
	funcSlice := make([]Func, 0)
	funcSlice = append(funcSlice, func(cache Cache, id int, key, field, value string) {
		start := time.Now()
		val := cache.Del(context.Background(), key)
		args := []interface{}{
			time.Now().Sub(start),
			id, key, field, value, val,
		}
		log.Printf("Del:[dur:%s,id:%d,key:%s,field:%s,value:%s,resule:%s]", args...)
	})

	funcSlice = append(funcSlice, func(cache Cache, id int, key, field, value string) {
		start := time.Now()
		val := cache.Expire(context.Background(), key, time.Second*5)
		args := []interface{}{
			time.Now().Sub(start),
			id, key, field, value, val,
		}
		log.Printf("Expire:[dur:%s,id:%d,key:%s,field:%s,value:%s,resule:%s]", args...)
	})
	funcSlice = append(funcSlice, func(cache Cache, id int, key, field, value string) {
		start := time.Now()
		val := cache.ExpireAt(context.Background(), key, time.Now().Add(time.Minute))
		args := []interface{}{
			time.Now().Sub(start),
			id, key, field, value, val,
		}
		log.Printf("ExpireAt:[dur:%s,id:%d,key:%s,field:%s,value:%s,resule:%s]", args...)
	})

	funcSlice = append(funcSlice, func(cache Cache, id int, key, field, value string) {
		start := time.Now()
		val := cache.Persist(context.Background(), key)
		args := []interface{}{
			time.Now().Sub(start),
			id, key, field, value, val,
		}
		log.Printf("Persist:[dur:%s,id:%d,key:%s,field:%s,value:%s,resule:%s]", args...)
	})

	funcSlice = append(funcSlice, func(cache Cache, id int, key, field, value string) {
		start := time.Now()
		val := cache.Set(context.Background(), key, value, time.Second*5)
		args := []interface{}{
			time.Now().Sub(start),
			id, key, field, value, val,
		}
		log.Printf("Set:[dur:%s,id:%d,key:%s,field:%s,value:%s,resule:%s]", args...)
	})

	funcSlice = append(funcSlice, func(cache Cache, id int, key, field, value string) {
		start := time.Now()
		val := cache.SetNX(context.Background(), key, value, time.Second*5)
		args := []interface{}{
			time.Now().Sub(start),
			id, key, field, value, val,
		}
		log.Printf("SetNX:[dur:%s,id:%d,key:%s,field:%s,value:%s,resule:%s]", args...)
	})

	funcSlice = append(funcSlice, func(cache Cache, id int, key, field, value string) {
		start := time.Now()
		val := cache.HDel(context.Background(), key, field)
		args := []interface{}{
			time.Now().Sub(start),
			id, key, field, value, val,
		}
		log.Printf("HDel:[dur:%s,id:%d,key:%s,field:%s,value:%s,resule:%s]", args...)
	})

	funcSlice = append(funcSlice, func(cache Cache, id int, key, field, value string) {
		start := time.Now()
		val := cache.HSet(context.Background(), key, field, value)
		args := []interface{}{
			time.Now().Sub(start),
			id, key, field, value, val,
		}
		log.Printf("HSet:[dur:%s,id:%d,key:%s,field:%s,value:%s,resule:%s]", args...)
	})

	funcSlice = append(funcSlice, func(cache Cache, id int, key, field, value string) {
		start := time.Now()
		val := cache.HSetNX(context.Background(), key, field, value)
		args := []interface{}{
			time.Now().Sub(start),
			id, key, field, value, val,
		}
		log.Printf("HSetNX:[dur:%s,id:%d,key:%s,field:%s,value:%s,resule:%s]", args...)
	})
	funcSlice = append(funcSlice, func(cache Cache, id int, key, field, value string) {
		start := time.Now()
		val := cache.LPush(context.Background(), key, value)
		args := []interface{}{
			time.Now().Sub(start),
			id, key, field, value, val,
		}
		log.Printf("LPush:[dur:%s,id:%d,key:%s,field:%s,value:%s,resule:%s]", args...)
	})

	funcSlice = append(funcSlice, func(cache Cache, id int, key, field, value string) {
		start := time.Now()
		val := cache.LPop(context.Background(), key)
		args := []interface{}{
			time.Now().Sub(start),
			id, key, field, value, val,
		}
		log.Printf("LPop:[dur:%s,id:%d,key:%s,field:%s,value:%s,resule:%s]", args...)
	})

	funcSlice = append(funcSlice, func(cache Cache, id int, key, field, value string) {
		start := time.Now()
		val := cache.LShift(context.Background(), key)
		args := []interface{}{
			time.Now().Sub(start),
			id, key, field, value, val,
		}
		log.Printf("LShift:[dur:%s,id:%d,key:%s,field:%s,value:%s,resule:%s]", args...)
	})

	funcSlice = append(funcSlice, func(cache Cache, id int, key, field, value string) {
		start := time.Now()
		val := cache.LTrim(context.Background(), key, rand.Int63n(10), rand.Int63n(10))
		args := []interface{}{
			time.Now().Sub(start),
			id, key, field, value, val,
		}
		log.Printf("LTrim:[dur:%s,id:%d,key:%s,field:%s,value:%s,resule:%s]", args...)
	})

	funcSlice = append(funcSlice, func(cache Cache, id int, key, field, value string) {
		start := time.Now()
		val := cache.LBPop(context.Background(), time.Second*10, key)
		args := []interface{}{
			time.Now().Sub(start),
			id, key, field, value, val,
		}
		log.Printf("LBPop:[dur:%s,id:%d,key:%s,field:%s,value:%s,resule:%s]", args...)
	})

	timer := time.NewTicker(time.Second)
	for range timer.C {
		log.Println("ticking...")
		for i := 0; i < 20; i++ {
			fn := funcSlice[rand.Intn(len(funcSlice))]
			go fn(randCache(), rand.Int(), randKey(), randField(), randValue())
		}
	}

}
