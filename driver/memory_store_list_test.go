package driver

import (
	"context"
	"reflect"
	"testing"
)

/*
*

	@author : Jerbe - The porter from Earth
	@time : 2023/9/15 12:29
	@describe :
*/
func Test_listStore_Push(t *testing.T) {
	s := newListStore()
	type args struct {
		ctx  context.Context
		key  string
		data []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "推入成功",
			args: args{
				ctx:  context.Background(),
				key:  "key",
				data: []interface{}{"v1", "v2", "v3", "v4"},
			},
			want:    4,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := s.Push(tt.args.ctx, tt.args.key, tt.args.data...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Push() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Push() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_listStore_Pop(t *testing.T) {
	s := newListStore()
	s.Push(context.Background(), "key", "v1", "v2")
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
			name: "获取尾部数据成功,v1",
			args: args{
				ctx: context.Background(),
				key: "key",
			},
			want:    "v1",
			wantErr: false,
		},
		{
			name: "获取尾部数据成功,v2",
			args: args{
				ctx: context.Background(),
				key: "key",
			},
			want:    "v2",
			wantErr: false,
		},
		{
			name: "获取尾部数据失败",
			args: args{
				ctx: context.Background(),
				key: "key",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Pop(tt.args.ctx, tt.args.key)
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

func Test_listStore_Rang(t *testing.T) {
	s := newListStore()
	s.Push(context.Background(), "key", "0", "1", "2", "3", "4", "5") // [5,4,3,2,1,0]
	type args struct {
		ctx   context.Context
		key   string
		start int64
		stop  int64
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "获取成功,0_0",
			args: args{
				ctx:   context.Background(),
				key:   "key",
				start: 0,
				stop:  0,
			},
			want:    []string{"5"},
			wantErr: false,
		},
		{
			name: "获取成功,1_3",
			args: args{
				ctx:   context.Background(),
				key:   "key",
				start: 1,
				stop:  3,
			},
			want:    []string{"4", "3", "2"},
			wantErr: false,
		},
		{
			name: "获取成功,-2_-1",
			args: args{
				ctx:   context.Background(),
				key:   "key",
				start: -2,
				stop:  -1,
			},
			want:    []string{"1", "0"},
			wantErr: false,
		},
		{
			name: "获取失败,3,0",
			args: args{
				ctx:   context.Background(),
				key:   "key",
				start: 3,
				stop:  0,
			},
			want:    []string{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := s.Rang(tt.args.ctx, tt.args.key, tt.args.start, tt.args.stop)
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

func Test_listStore_Shift(t *testing.T) {
	s := newListStore()
	s.Push(context.Background(), "key", "0", "1", "2") // [2,1,0]

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
			name: "推出头部成功,2",
			args: args{
				ctx: context.Background(),
				key: "key",
			},
			want:    "2",
			wantErr: false,
		},
		{
			name: "推出头部成功,1",
			args: args{
				ctx: context.Background(),
				key: "key",
			},
			want:    "1",
			wantErr: false,
		},
		{
			name: "推出头部成功,0",
			args: args{
				ctx: context.Background(),
				key: "key",
			},
			want:    "0",
			wantErr: false,
		},
		{
			name: "推出头部失败",
			args: args{
				ctx: context.Background(),
				key: "key",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := s.Shift(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Shift() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Shift() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_listStore_Trim(t *testing.T) {
	s := newListStore()

	type args struct {
		ctx   context.Context
		key   string
		start int64
		stop  int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "裁剪并只保留范围内数据",
			args: args{
				ctx: context.Background(),
				key: func() string {
					s.Push(context.Background(), "key", "0", "1", "2", "3", "4", "5") // [5,4,3,2,1,0]
					return "key"
				}(),
				start: 0,
				stop:  0,
			},
			wantErr: false,
		},
		{
			name: "裁剪并只保留范围内数据",
			args: args{
				ctx: context.Background(),
				key: func() string {
					s.Push(context.Background(), "key1", "0", "1", "2", "3", "4", "5") // [5,4,3,2,1,0]
					return "key1"
				}(),
				start: 1,
				stop:  -1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.Trim(tt.args.ctx, tt.args.key, tt.args.start, tt.args.stop); (err != nil) != tt.wantErr {
				t.Errorf("Trim() error = %v, wantErr %v", err, tt.wantErr)
			}
			r, err := s.Rang(context.Background(), "key", 0, 0)
			t.Log("rang 1", r, err)
			r, err = s.Rang(context.Background(), "key1", 0, -1)
			t.Log("rang 2", r, err)
		})
	}
}
