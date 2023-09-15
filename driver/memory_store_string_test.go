package driver

import (
	"context"
	"reflect"
	"testing"
	"time"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/15 09:07
  @describe :
*/

func Test_stringStore_Get(t *testing.T) {

	ss := newStringStore()
	ss.Set(context.Background(), "key1", "value1", 0)

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
				key: "key",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "获取成功",
			args: args{
				ctx: context.Background(),
				key: "key1",
			},
			want:    "value1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := ss.Get(tt.args.ctx, tt.args.key)
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

func Test_stringStore_MGet(t *testing.T) {
	ss := newStringStore()
	ss.Set(context.Background(), "key", "value", 0)

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
		{
			name: "获取成功",
			args: args{
				ctx:  context.Background(),
				keys: []string{"key"},
			},
			want:    []interface{}{"value"},
			wantErr: false,
		},
		{
			name: "获取失败",
			args: args{
				ctx:  context.Background(),
				keys: []string{"key1", "key2"},
			},
			want:    []interface{}{nil, nil},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := ss.MGet(tt.args.ctx, tt.args.keys...)
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

func Test_stringStore_Set(t *testing.T) {
	ss := newStringStore()
	type args struct {
		ctx        context.Context
		key        string
		data       interface{}
		expiration time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "设置成功",
			args: args{
				ctx:        context.Background(),
				key:        "key",
				data:       nil,
				expiration: 0,
			},
		},
		{
			name: "再次设置成功",
			args: args{
				ctx:        context.Background(),
				key:        "key",
				data:       nil,
				expiration: 0,
			},
		},
		{
			name: "再次设置成功",
			args: args{
				ctx:        context.Background(),
				key:        "key",
				data:       nil,
				expiration: -1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ss.Set(tt.args.ctx, tt.args.key, tt.args.data, tt.args.expiration); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_stringStore_SetNX(t *testing.T) {
	ss := newStringStore()
	type args struct {
		ctx        context.Context
		key        string
		data       interface{}
		expiration time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "设置成功",
			args: args{
				ctx:        context.Background(),
				key:        "key",
				data:       "value",
				expiration: 0,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "设置失败",
			args: args{
				ctx:        context.Background(),
				key:        "key",
				data:       "value",
				expiration: 0,
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ss.SetNX(tt.args.ctx, tt.args.key, tt.args.data, tt.args.expiration)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetNX() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SetNX() got = %v, want %v", got, tt.want)
			}
		})
	}
}
