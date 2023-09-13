package jcache

import (
	"context"
	"github.com/jerbe/jcache/driver"
	"math/rand"
	"reflect"
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
	return NewStringClient(driver.NewMemoryString())
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
		dst any
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

func TestStringClient_Del(t *testing.T) {
	type fields struct {
		drivers []driver.String
	}
	type args struct {
		ctx  context.Context
		keys []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &StringClient{
				drivers: tt.fields.drivers,
			}
			if err := cli.Del(tt.args.ctx, tt.args.keys...); (err != nil) != tt.wantErr {
				t.Errorf("Del() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStringClient_Exists(t *testing.T) {
	type fields struct {
		drivers []driver.String
	}
	type args struct {
		ctx  context.Context
		keys []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &StringClient{
				drivers: tt.fields.drivers,
			}
			got, err := cli.Exists(tt.args.ctx, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Exists() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringClient_Expire(t *testing.T) {
	type fields struct {
		drivers []driver.String
	}
	type args struct {
		ctx        context.Context
		key        string
		expiration time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &StringClient{
				drivers: tt.fields.drivers,
			}
			if err := cli.Expire(tt.args.ctx, tt.args.key, tt.args.expiration); (err != nil) != tt.wantErr {
				t.Errorf("Expire() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStringClient_ExpireAt(t *testing.T) {
	type fields struct {
		drivers []driver.String
	}
	type args struct {
		ctx context.Context
		key string
		at  time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &StringClient{
				drivers: tt.fields.drivers,
			}
			if err := cli.ExpireAt(tt.args.ctx, tt.args.key, tt.args.at); (err != nil) != tt.wantErr {
				t.Errorf("ExpireAt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStringClient_Get(t *testing.T) {
	type fields struct {
		drivers []driver.String
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &StringClient{
				drivers: tt.fields.drivers,
			}
			got, err := cli.Get(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringClient_GetAndScan(t *testing.T) {
	type fields struct {
		drivers []driver.String
	}
	type args struct {
		ctx context.Context
		dst any
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &StringClient{
				drivers: tt.fields.drivers,
			}
			if err := cli.GetAndScan(tt.args.ctx, tt.args.dst, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("GetAndScan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStringClient_MGet(t *testing.T) {
	type fields struct {
		drivers []driver.String
	}
	type args struct {
		ctx  context.Context
		keys []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &StringClient{
				drivers: tt.fields.drivers,
			}
			got, err := cli.MGet(tt.args.ctx, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("MGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MGet() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringClient_MGetAndScan(t *testing.T) {
	type fields struct {
		drivers []driver.String
	}
	type args struct {
		ctx  context.Context
		dst  any
		keys []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &StringClient{
				drivers: tt.fields.drivers,
			}
			if err := cli.MGetAndScan(tt.args.ctx, tt.args.dst, tt.args.keys...); (err != nil) != tt.wantErr {
				t.Errorf("MGetAndScan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStringClient_Set(t *testing.T) {
	type fields struct {
		drivers []driver.String
	}
	type args struct {
		ctx        context.Context
		key        string
		data       any
		expiration time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &StringClient{
				drivers: tt.fields.drivers,
			}
			if err := cli.Set(tt.args.ctx, tt.args.key, tt.args.data, tt.args.expiration); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStringClient_SetNX(t *testing.T) {
	type fields struct {
		drivers []driver.String
	}
	type args struct {
		ctx        context.Context
		key        string
		data       any
		expiration time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &StringClient{
				drivers: tt.fields.drivers,
			}
			if err := cli.SetNX(tt.args.ctx, tt.args.key, tt.args.data, tt.args.expiration); (err != nil) != tt.wantErr {
				t.Errorf("SetNX() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
