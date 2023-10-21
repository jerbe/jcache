package main

import (
	"bufio"
	"context"
	"flag"
	"github.com/jerbe/jcache/v2/driver"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/22 14:09
  @describe :
*/

var (
	prefix = flag.String("prefix", "cache", "前缀")
	port   = flag.Int("port", 9090, "端口号")
)

func main() {
	flag.Parse()
	cfg := driver.MemoryConfig{Port: *port, Prefix: *prefix, EtcdConfig: clientv3.Config{Endpoints: []string{"127.0.0.1:2379"}}}
	mem, err := driver.NewMemoryWithConfig(cfg)
	if err != nil {
		log.Fatalf("初始化内存驱动失败. 原因:[%v]", err)
	}

	inputScanner := bufio.NewScanner(os.Stdin)
	for inputScanner.Scan() {
		values := strings.Split(inputScanner.Text(), " ")
		if len(values) == 0 {
			continue
		}
		action := values[0]
		switch action {
		case "exists":
			val := mem.Exists(context.Background(), values[1:]...)
			log.Printf(">> val:[%v], err:[%v]", val.Val(), val.Err())
		case "set":
			var ttl int64
			if len(values) == 4 {
				ttl, err = strconv.ParseInt(values[3], 10, 64)
				if err != nil {
					log.Printf("set failure. reason:[%v]", err)
					continue
				}
			}
			dur := time.Second * time.Duration(ttl)
			val := mem.Set(context.Background(), values[1], values[2], dur)
			log.Printf(">> val:[%v], err:[%v]", val.Val(), val.Err())
		case "get":
			val := mem.Get(context.Background(), values[1])
			log.Printf(">> val:[%v], err:[%v]", val.Val(), val.Err())
		case "setnx":
			ttl, err := strconv.ParseInt(values[3], 10, 64)
			if err != nil {
				log.Printf("setnx failure. reason:[%v]", err)
				continue
			}
			dur := time.Second * time.Duration(ttl)
			val := mem.SetNX(context.Background(), values[1], values[2], dur)
			log.Printf(">> val:[%v], err:[%v]", val.Val(), val.Err())
		case "del":
			val := mem.Del(context.Background(), values[1:]...)
			log.Printf(">> val:[%v], err:[%v]", val.Val(), val.Err())
		case "expire":
			ttl, err := strconv.ParseInt(values[2], 10, 64)
			if err != nil {
				log.Printf("expire failure. reason:[%v]", err)
				continue
			}
			dur := time.Second * time.Duration(ttl)
			val := mem.Expire(context.Background(), values[1], dur)
			log.Printf(">> val:[%v], err:[%v]", val.Val(), val.Err())
		case "expireat":
			at, err := time.Parse(time.DateTime, values[2])
			if err != nil {
				log.Printf("expireat failure. reason:[%v]", err)
				continue
			}
			val := mem.ExpireAt(context.Background(), values[1], at)
			log.Printf(">> val:[%v], err:[%v]", val.Val(), val.Err())
		case "hset":
			args := make([]interface{}, 0, len(values)-2)
			for i := 0; i < cap(args); i++ {
				args = append(args, values[i+2])
			}
			val := mem.HSet(context.Background(), values[1], args...)
			log.Printf(">> val:[%v], err:[%v]", val.Val(), val.Err())
		case "hsetnx":
			val := mem.HSetNX(context.Background(), values[1], values[2], values[3])
			log.Printf(">> val:[%v], err:[%v]", val.Val(), val.Err())
		case "hget":
			val := mem.HGet(context.Background(), values[1], values[2])
			log.Printf(">> val:[%v], err:[%v]", val.Val(), val.Err())
		case "hdel":
			val := mem.HDel(context.Background(), values[1], values[2:]...)
			log.Printf(">> val:[%v], err:[%v]", val.Val(), val.Err())
		case "hlen":
			val := mem.HLen(context.Background(), values[1])
			log.Printf(">> val:[%v], err:[%v]", val.Val(), val.Err())
		case "hkeys":
			val := mem.HKeys(context.Background(), values[1])
			log.Printf(">> val:[%v], err:[%v]", val.Val(), val.Err())
		case "hvals":
			val := mem.HVals(context.Background(), values[1])
			log.Printf(">> val:[%v], err:[%v]", val.Val(), val.Err())
		case "hgetall":
			val := mem.HGetAll(context.Background(), values[1])
			log.Printf(">> val:[%v], err:[%v]", val.Val(), val.Err())
		case "hexists":
			val := mem.HExists(context.Background(), values[1], values[2])
			log.Printf(">> val:[%v], err:[%v]", val.Val(), val.Err())
		case "lpush":
			args := make([]interface{}, 0, len(values)-2)
			for i := 0; i < cap(args); i++ {
				args = append(args, values[i+2])
			}
			val := mem.LPush(context.Background(), values[1], args...)
			log.Printf(">> val:[%v], err:[%v]", val.Val(), val.Err())
		case "ltrim":
			start, err := strconv.ParseInt(values[2], 10, 64)
			if err != nil {
				log.Printf("ltrim failure. reason:[%v]", err)
				continue
			}
			stop, err := strconv.ParseInt(values[3], 10, 64)
			if err != nil {
				log.Printf("ltrim failure. reason:[%v]", err)
				continue
			}
			val := mem.LTrim(context.Background(), values[1], start, stop)
			log.Printf(">> val:[%v], err:[%v]", val.Val(), val.Err())
		case "lrang":
			start, err := strconv.ParseInt(values[2], 10, 64)
			if err != nil {
				log.Printf("ltrim failure. reason:[%v]", err)
				continue
			}
			stop, err := strconv.ParseInt(values[3], 10, 64)
			if err != nil {
				log.Printf("ltrim failure. reason:[%v]", err)
				continue
			}
			val := mem.LRang(context.Background(), values[1], start, stop)
			log.Printf(">> val:[%v], err:[%v]", val.Val(), val.Err())
		case "lpop":
			val := mem.LPop(context.Background(), values[1])
			log.Printf(">> val:[%v], err:[%v]", val.Val(), val.Err())
		case "lshift":
			val := mem.LShift(context.Background(), values[1])
			log.Printf(">> val:[%v], err:[%v]", val.Val(), val.Err())
		}

	}

}
