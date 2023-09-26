package driver

import (
	"context"
	"math/rand"
	"reflect"
	"strconv"
	"testing"
	"time"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/25 17:50
  @describe :
*/

func Test_sortSetStore_ZAdd(t *testing.T) {
	s := newSortSetStore()

	type args struct {
		ctx     context.Context
		key     string
		members []Z
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "添加成功,A,Score:1.0",
			args: args{
				ctx: context.TODO(),
				key: "abc",
				members: []Z{
					{Member: "A", Score: 1.0},
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "添加成功,A,Score:0",
			args: args{
				ctx: context.TODO(),
				key: "abc",
				members: []Z{
					{Member: "A", Score: 0.0},
				},
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "添加成功,B,Score:0.0",
			args: args{
				ctx: context.TODO(),
				key: "abc",
				members: []Z{
					{Member: "B", Score: 0.0},
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "添加成功,C,Score:-1.0",
			args: args{
				ctx: context.TODO(),
				key: "abc",
				members: []Z{
					{Member: "C", Score: -1.0},
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "添加成功,C,Score:2.0",
			args: args{
				ctx: context.TODO(),
				key: "abc",
				members: []Z{
					{Member: "C", Score: 2.0},
				},
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "添加成功,D,Score:99.0",
			args: args{
				ctx: context.TODO(),
				key: "abc",
				members: []Z{
					{Member: "D", Score: 99.0},
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "添加成功,E,Score:98.0",
			args: args{
				ctx: context.TODO(),
				key: "abc",
				members: []Z{
					{Member: "E", Score: 98.0},
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "添加成功,F,Score:-98.0",
			args: args{
				ctx: context.TODO(),
				key: "abc",
				members: []Z{
					{Member: "F", Score: -98.0},
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "添加成功,F,Score:99.0",
			args: args{
				ctx: context.TODO(),
				key: "abc",
				members: []Z{
					{Member: "F", Score: 99.0},
				},
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "添加成功,E,Score:0.0",
			args: args{
				ctx: context.TODO(),
				key: "abc",
				members: []Z{
					{Member: "E", Score: 0.0},
				},
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "添加成功,位置不变,D,Score:99.0",
			args: args{
				ctx: context.TODO(),
				key: "abc",
				members: []Z{
					{Member: "D", Score: 99},
				},
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "添加成功,位置不变,D,Score:3",
			args: args{
				ctx: context.TODO(),
				key: "abc",
				members: []Z{
					{Member: "D", Score: 3},
				},
			},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.ZAdd(tt.args.ctx, tt.args.key, tt.args.members...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZAdd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ZAdd() got = %v, want %v", got, tt.want)
			}
		})
	}
}

var sx = newSortSetStore()

func Benchmark_sortSetStore_ZAdd(b *testing.B) {
	s := sx
	b.SetParallelism(1000)
	rand.Seed(time.Now().UnixNano())
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			b.StopTimer()
			members := make([]Z, 0)
			var key = strconv.FormatInt(rand.Int63n(5), 16)
			for i := 0; i < 10; i++ {
				var member = strconv.FormatInt(rand.Int63(), 16)
				var score = rand.Float64()
				members = append(members, Z{
					Score:  score,
					Member: member,
				})
			}
			b.StartTimer()
			s.ZAdd(context.TODO(), key, members...)
		}
	})
}

func Test_sortSetStore_ZRange(t *testing.T) {
	s := newSortSetStore()
	s.ZAdd(context.TODO(), "abc",
		Z{Score: 2, Member: "A"},
		Z{Score: 1, Member: "B"},
		Z{Score: 3, Member: "C"},
		Z{Score: 4, Member: "D"},
		Z{Score: 5, Member: "E"},
	)

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
			name: "获取正确,0_-1",
			args: args{
				ctx:   context.TODO(),
				key:   "abc",
				start: 0,
				stop:  -1,
			},
			want: []string{
				"B", "A", "C", "D", "E",
			},
			wantErr: false,
		},
		{
			name: "获取正确,1_5",
			args: args{
				ctx:   context.TODO(),
				key:   "abc",
				start: 1,
				stop:  5,
			},
			want: []string{
				"A", "C", "D", "E",
			},
			wantErr: false,
		},
		{
			name: "获取正确,0_0",
			args: args{
				ctx:   context.TODO(),
				key:   "abc",
				start: 0,
				stop:  0,
			},
			want: []string{
				"B",
			},
			wantErr: false,
		},
		{
			name: "获取正确,-1_-1",
			args: args{
				ctx:   context.TODO(),
				key:   "abc",
				start: -1,
				stop:  -1,
			},
			want: []string{
				"E",
			},
			wantErr: false,
		},
		{
			name: "获取正确,3_2",
			args: args{
				ctx:   context.TODO(),
				key:   "abc",
				start: 3,
				stop:  2,
			},
			want:    []string{},
			wantErr: false,
		},
		{
			name: "获取正确,-2_-3",
			args: args{
				ctx:   context.TODO(),
				key:   "abc",
				start: -2,
				stop:  -3,
			},
			want:    []string{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.ZRange(tt.args.ctx, tt.args.key, tt.args.start, tt.args.stop)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZRange() got = %v, want %v", got, tt.want)
			}
		})
	}
}
