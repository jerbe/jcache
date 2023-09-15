package jcache

import (
	"context"
	"github.com/jerbe/jcache/driver"
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

	memoryDriver := driver.NewMemoryString()
	return NewStringClient(memoryDriver)
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
			_, err = cli.Get(context.Background(), key)
			if err != nil {
				b.Logf("Get:发生错误:%+v", err)
			}

			err = cli.Expire(context.Background(), key, time.Minute)
			if err != nil {
				b.Logf("Expire:发生错误:%+v", err)
			}

			err = cli.ExpireAt(context.Background(), key, time.Now().Add(time.Minute))
			if err != nil {
				b.Logf("ExpireAt:发生错误:%+v", err)
			}
			err = cli.Del(context.Background(), key)
			if err != nil {
				b.Logf("Del:发生错误:%+v", err)
			}

			//b.Log(v)
		}
	})
}

func TestStringClient_CheckAndGet(t *testing.T) {
	cli := newStringClient()
	cli.Set(context.Background(), "my_love", "you_love", time.Hour)

	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "获取失败",
			args: args{
				ctx: context.Background(),
				key: "test",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "获取成功",
			args: args{
				ctx: context.Background(),
				key: "my_love",
			},
			want:    "you_love",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cli.CheckAndGet(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckAndGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckAndGet() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringClient_CheckAndScan(t *testing.T) {
	cli := newStringClient()
	cli.Set(context.Background(), "my_love", "you_love", time.Hour)
	cli.Set(context.Background(), "you_love", "", time.Hour)
	type args struct {
		ctx context.Context
		dst interface{}
		key string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "获取失败",
			args: args{
				ctx: context.Background(),
				dst: new(string),
				key: "my_love1",
			},
			wantErr: true,
		},
		{
			name: "获取成功",
			args: args{
				ctx: context.Background(),
				dst: new(string),
				key: "my_love",
			},
			wantErr: false,
		},
		{
			name: "获取失败",
			args: args{
				ctx: context.Background(),
				dst: new(string),
				key: "you_love",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cli.CheckAndScan(tt.args.ctx, tt.args.dst, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("CheckAndScan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
