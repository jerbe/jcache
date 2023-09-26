package jcache

import (
	"context"
	"testing"
	"time"

	"github.com/jerbe/jcache/v2/driver"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/16 22:24
  @describe :
*/

func newClient() *Client {
	return NewClient(driver.NewMemory())
}

func Test_Client_Del(t *testing.T) {
	cli := newClient()

	type fields struct {
		drivers []driver.Common
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
		{
			name: "删除成功",
			args: args{
				ctx:  context.Background(),
				keys: []string{"key", "key1", "key2"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cli.Del(tt.args.ctx, tt.args.keys...).Err(); (err != nil) != tt.wantErr {
				t.Errorf("Del() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_Client_Exists(t *testing.T) {
	cli := newClient()
	cli.Set(context.Background(), "key1", "value1", 0)

	type args struct {
		ctx  context.Context
		keys []string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "判断失败",
			args: args{
				ctx:  context.Background(),
				keys: []string{"key"},
			},
			want: 0,
		},
		{
			name: "判断成功,1",
			args: args{
				ctx:  context.Background(),
				keys: []string{"key1"},
			},
			want: 1,
		},
		{
			name: "判断成功,2",
			args: args{
				ctx:  context.Background(),
				keys: []string{"key1", "key"},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cli.Exists(tt.args.ctx, tt.args.keys...).Result()
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

func Test_Client_Expire(t *testing.T) {
	cli := newClient()
	cli.Set(context.Background(), "key1", "value1", 0)

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
		{
			name: "设置失败",
			args: args{
				ctx:        context.Background(),
				key:        "key",
				expiration: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cli.Expire(tt.args.ctx, tt.args.key, tt.args.expiration).Err(); (err != nil) != tt.wantErr {
				t.Errorf("Expire() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_Client_ExpireAt(t *testing.T) {
	type fields struct {
		drivers []driver.Common
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
			cli := &BaseClient{
				drivers: tt.fields.drivers,
			}
			if err := cli.ExpireAt(tt.args.ctx, tt.args.key, tt.args.at); (err != nil) != tt.wantErr {
				t.Errorf("ExpireAt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
