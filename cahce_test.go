package jcache

import (
	"context"
	"reflect"
	"testing"
	"time"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/8 22:52
  @describe :
*/

func TestCheckAndGet(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckAndGet(tt.args.ctx, tt.args.key)
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

func TestCheckAndScan(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckAndScan(tt.args.ctx, tt.args.dst, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("CheckAndScan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDel(t *testing.T) {
	type args struct {
		ctx  context.Context
		keys []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Del(tt.args.ctx, tt.args.keys...); (err != nil) != tt.wantErr {
				t.Errorf("Del() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExists(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Exists(tt.args.ctx, tt.args.key)
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

func TestExpire(t *testing.T) {
	type args struct {
		ctx        context.Context
		key        string
		expiration time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Expire(tt.args.ctx, tt.args.key, tt.args.expiration); (err != nil) != tt.wantErr {
				t.Errorf("Expire() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGet(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get(tt.args.ctx, tt.args.key)
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

func TestGetAndScan(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetAndScan(tt.args.ctx, tt.args.dst, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("GetAndScan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHDel(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		fields []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := HDel(tt.args.ctx, tt.args.key, tt.args.fields...); (err != nil) != tt.wantErr {
				t.Errorf("HDel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHGet(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		field string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HGet(tt.args.ctx, tt.args.key, tt.args.field)
			if (err != nil) != tt.wantErr {
				t.Errorf("HGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HGet() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHGetAndScan(t *testing.T) {
	type args struct {
		ctx   context.Context
		dst   any
		key   string
		field string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := HGetAndScan(tt.args.ctx, tt.args.dst, tt.args.key, tt.args.field); (err != nil) != tt.wantErr {
				t.Errorf("HGetAndScan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHKeys(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HKeys(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("HKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HKeys() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHKeysAndScan(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
		dst any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := HKeysAndScan(tt.args.ctx, tt.args.key, tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("HKeysAndScan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHLen(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HLen(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("HLen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HLen() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHMGet(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		fields []string
	}
	tests := []struct {
		name    string
		args    args
		want    []any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HMGet(tt.args.ctx, tt.args.key, tt.args.fields...)
			if (err != nil) != tt.wantErr {
				t.Errorf("HMGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HMGet() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHMGetAndScan(t *testing.T) {
	type args struct {
		ctx    context.Context
		dst    any
		key    string
		fields []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := HMGetAndScan(tt.args.ctx, tt.args.dst, tt.args.key, tt.args.fields...); (err != nil) != tt.wantErr {
				t.Errorf("HMGetAndScan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHSet(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		values []any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := HSet(tt.args.ctx, tt.args.key, tt.args.values...); (err != nil) != tt.wantErr {
				t.Errorf("HSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHVals(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HVals(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("HVals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HVals() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHValsScan(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := HValsScan(tt.args.ctx, tt.args.dst, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("HValsScan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMGet(t *testing.T) {
	type args struct {
		ctx  context.Context
		keys []string
	}
	tests := []struct {
		name    string
		args    args
		want    []interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MGet(tt.args.ctx, tt.args.keys...)
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

func TestMGetAndScan(t *testing.T) {
	type args struct {
		ctx  context.Context
		dst  any
		keys []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MGetAndScan(tt.args.ctx, tt.args.dst, tt.args.keys...); (err != nil) != tt.wantErr {
				t.Errorf("MGetAndScan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPop(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Pop(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pop() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Pop() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPopAndScan(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PopAndScan(tt.args.ctx, tt.args.dst, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("PopAndScan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPush(t *testing.T) {
	type args struct {
		ctx  context.Context
		key  string
		data []any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Push(tt.args.ctx, tt.args.key, tt.args.data...); (err != nil) != tt.wantErr {
				t.Errorf("Push() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRang(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		limit int64
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Rang(tt.args.ctx, tt.args.key, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Rang() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Rang() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRangAndScan(t *testing.T) {
	type args struct {
		ctx   context.Context
		dst   any
		key   string
		limit int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RangAndScan(tt.args.ctx, tt.args.dst, tt.args.key, tt.args.limit); (err != nil) != tt.wantErr {
				t.Errorf("RangAndScan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSet(t *testing.T) {
	type args struct {
		ctx        context.Context
		key        string
		data       any
		expiration time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Set(tt.args.ctx, tt.args.key, tt.args.data, tt.args.expiration); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetNX(t *testing.T) {
	type args struct {
		ctx        context.Context
		key        string
		data       any
		expiration time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetNX(tt.args.ctx, tt.args.key, tt.args.data, tt.args.expiration); (err != nil) != tt.wantErr {
				t.Errorf("SetNX() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
