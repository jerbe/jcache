package driver

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"testing"
	"time"
)

/*
*

	@author : Jerbe - The porter from Earth
	@time : 2023/9/15 12:29
	@describe :
*/
func Test_listStore_LPush(t *testing.T) {
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

			got, err := s.LPush(tt.args.ctx, tt.args.key, tt.args.data...)
			if (err != nil) != tt.wantErr {
				t.Errorf("LPush() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LPush() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_listStore_LPop(t *testing.T) {
	s := newListStore()
	s.LPush(context.Background(), "key", "v1", "v2")
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
			got, err := s.LPop(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("LPop() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LPop() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_listStore_LRang(t *testing.T) {
	s := newListStore()
	s.LPush(context.Background(), "key", "x", "1", "你好", "3", "g", "5") // [5,g,3,你好,1,x]
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
			want:    []string{"g", "3", "你好"},
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
			want:    []string{"1", "x"},
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
			got, err := s.LRang(tt.args.ctx, tt.args.key, tt.args.start, tt.args.stop)
			if (err != nil) != tt.wantErr {
				t.Errorf("LRang() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LRang() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_listStore_LShift(t *testing.T) {
	s := newListStore()
	s.LPush(context.Background(), "key", "0", "1", "2") // [2,1,0]

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

			got, err := s.LShift(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("LShift() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LShift() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_listStore_LTrim(t *testing.T) {
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
					s.LPush(context.Background(), "key", "0", "1", "2", "3", "4", "5") // [5,4,3,2,1,0]
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
					s.LPush(context.Background(), "key1", "0", "1", "2", "3", "4", "5") // [5,4,3,2,1,0]
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
			if err := s.LTrim(tt.args.ctx, tt.args.key, tt.args.start, tt.args.stop); (err != nil) != tt.wantErr {
				t.Errorf("LTrim() error = %v, wantErr %v", err, tt.wantErr)
			}
			r, err := s.LRang(context.Background(), "key", 0, 0)
			t.Log("rang 1", r, err)
			r, err = s.LRang(context.Background(), "key1", 0, -1)
			t.Log("rang 2", r, err)
		})
	}
}

func Test_listStore_LLen(t *testing.T) {
	s := newListStore()
	s.LPush(context.Background(), "key", "0", 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20)

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
		{
			name: "获取正确",
			args: args{
				ctx: context.Background(),
				key: "key",
			},
			want:    21,
			wantErr: false,
		},
		{
			name: "获取失败,没有该key",
			args: args{
				ctx: context.Background(),
				key: "key1",
			},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.LLen(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("LLen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LLen() got = %v, want %v", got, tt.want)
			}

		})
	}
}

func Benchmark_listStore_LPush(b *testing.B) {
	s := newListStore()
	b.SetParallelism(10000)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var key = strconv.FormatInt(rand.Int63n(100), 10)
			var value = strconv.FormatInt(rand.Int63(), 10)
			s.LPush(context.Background(), key, value)
		}
	})

}

func Benchmark_listStore_LPop(b *testing.B) {
	s := newListStore()
	b.SetParallelism(10000)

	b.StopTimer()
	b.Log("init testing...")
	for i := 0; i < 10000000; i++ {
		var key = strconv.FormatInt(rand.Int63n(100), 10)
		var value = strconv.FormatInt(rand.Int63(), 10)
		s.LPush(context.Background(), key, value)
	}
	b.StartTimer()
	b.Log("start testing...")
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var key = strconv.FormatInt(rand.Int63n(100), 10)
			s.LPop(context.Background(), key)

		}
	})
}

func Benchmark_listStore_LShift(b *testing.B) {
	s := newListStore()
	b.SetParallelism(10000)

	b.StopTimer()
	for i := 0; i < 10000000; i++ {
		var key = strconv.FormatInt(rand.Int63n(100), 10)
		var value = strconv.FormatInt(rand.Int63(), 10)
		s.LPush(context.Background(), key, value)
	}

	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var key = strconv.FormatInt(rand.Int63n(100), 10)
			s.LShift(context.Background(), key)

		}
	})
}

func Benchmark_listStore_LRang(b *testing.B) {
	s := newListStore()
	b.SetParallelism(10000)

	b.StopTimer()
	for i := 0; i < 10000000; i++ {
		var key = strconv.FormatInt(rand.Int63n(100), 10)
		var value = strconv.FormatInt(rand.Int63(), 10)
		s.LPush(context.Background(), key, value)
	}
	b.StartTimer()

	b.Log("开始:>> ")
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var key = strconv.FormatInt(rand.Int63n(100), 10)
			s.LRang(context.Background(), key, rand.Int63n(10000000), rand.Int63n(10000000))
		}
	})
	b.Log("<<:结束")
}

func Test_listStore_LX(t *testing.T) {
	s := newListStore()
	now := time.Now()
	for i := 0; i < 1000000; i++ {
		var key = strconv.FormatInt(rand.Int63n(100), 10)
		var value = strconv.FormatInt(rand.Int63(), 10)
		s.LPush(context.Background(), key, value)
	}
	fmt.Println("初始化耗时", time.Now().Sub(now))
	now = time.Now()
	for i := 0; i < 10000; i++ {
		var key = strconv.FormatInt(rand.Int63n(100), 10)
		var value = strconv.FormatInt(rand.Int63(), 10)
		s.LPush(context.Background(), key, value)
	}
	fmt.Println("再插入一批看看", time.Now().Sub(now))
	key := strconv.FormatInt(rand.Int63n(100), 10)
	s.LRang(context.Background(), key, rand.Int63n(10000000), rand.Int63n(10000000))

	fmt.Println("获取耗时:", time.Now().Sub(now))

}
