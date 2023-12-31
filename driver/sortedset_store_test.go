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
		members []SZ
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "添加成功,A,Score:1.0,[A]",
			args: args{
				ctx: context.TODO(),
				key: "abc",
				members: []SZ{
					{Member: "A", Score: 1.0},
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "添加成功,A,Score:0,[A]",
			args: args{
				ctx: context.TODO(),
				key: "abc",
				members: []SZ{
					{Member: "A", Score: 0.0},
				},
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "添加成功,B,Score:0.0[A B]",
			args: args{
				ctx: context.TODO(),
				key: "abc",
				members: []SZ{
					{Member: "B", Score: 0.0},
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "添加成功,C,Score:-1.0,[C A B]",
			args: args{
				ctx: context.TODO(),
				key: "abc",
				members: []SZ{
					{Member: "C", Score: -1.0},
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "添加成功,C,Score:2.0,[A B C]",
			args: args{
				ctx: context.TODO(),
				key: "abc",
				members: []SZ{
					{Member: "C", Score: 2.0},
				},
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "添加成功,D,Score:99.0,[A B C D]",
			args: args{
				ctx: context.TODO(),
				key: "abc",
				members: []SZ{
					{Member: "D", Score: 99.0},
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "添加成功,E,Score:98.0,[A B C E D]",
			args: args{
				ctx: context.TODO(),
				key: "abc",
				members: []SZ{
					{Member: "E", Score: 98.0},
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "添加成功,F,Score:-98.0,[F A B C E D]",
			args: args{
				ctx: context.TODO(),
				key: "abc",
				members: []SZ{
					{Member: "F", Score: -98.0},
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "添加成功,F,Score:99.0,[A B C E D F]",
			args: args{
				ctx: context.TODO(),
				key: "abc",
				members: []SZ{
					{Member: "F", Score: 99.0},
				},
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "添加成功,E,Score:0.0,[A B E C D F]",
			args: args{
				ctx: context.TODO(),
				key: "abc",
				members: []SZ{
					{Member: "E", Score: 0.0},
				},
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "添加成功,位置不变,D,Score:99.0,[A B E C D F]",
			args: args{
				ctx: context.TODO(),
				key: "abc",
				members: []SZ{
					{Member: "D", Score: 99},
				},
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "添加成功,位置不变,D,Score:3,[A B E C D F]",
			args: args{
				ctx: context.TODO(),
				key: "abc",
				members: []SZ{
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

func Test_sortSetStore_X_ZAdd(t *testing.T) {
	sx := newSortSetStore()
	m := []SZ{}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 1000; i++ {
		m = append(m, SZ{Member: fmt.Sprintf("%x", rand.Int63()), Score: rand.Float64()})
	}
	start := time.Now()
	cnt, err := sx.ZAdd(context.Background(), "abc", m...)
	//fmt.Printf("数据:%+v", sx.values["abc"].(*sortedSetValue).rankList)
	fmt.Println("耗时:", time.Now().Sub(start), "成员数量:", cnt, "错误:", err)
}

func Benchmark_sortSetStore_ZAdd(b *testing.B) {
	s := newSortSetStore()
	b.SetParallelism(10000)
	rand.Seed(time.Now().UnixNano())
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			b.StopTimer()
			members := make([]SZ, 0)
			var key = strconv.FormatInt(rand.Int63n(5), 16)
			for i := 0; i < 20; i++ {
				var member = strconv.FormatInt(rand.Int63(), 16)
				var score = rand.Float64()
				members = append(members, SZ{
					Score:  score,
					Member: member,
				})
			}
			b.StartTimer()
			s.ZAdd(context.TODO(), key, members...)
			//b.StopTimer()
			//go log.Println(len(s.values[key].(*sortedSetValue).rankList))
			//b.StartTimer()
		}
	})
}

func Test_sortSetStore_ZRange(t *testing.T) {
	s := newSortSetStore()
	s.ZAdd(context.TODO(), "abc",
		SZ{Score: 2, Member: "A"},
		SZ{Score: 1, Member: "B"},
		SZ{Score: 3, Member: "C"},
		SZ{Score: 4, Member: "D"},
		SZ{Score: 5, Member: "E"},
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

func Test_sortSetStore_ZRem(t *testing.T) {
	s := newSortSetStore()
	KEY := "ABC"

	member := []SZ{}
	for i := 0; i < 10; i++ {
		member = append(member, SZ{Member: strconv.Itoa(i), Score: float64(i)})
	}
	s.ZAdd(context.TODO(), KEY, member...)
	type args struct {
		ctx     context.Context
		key     string
		members []string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "正确删除,2,3",
			args: args{
				ctx:     context.TODO(),
				key:     KEY,
				members: []string{"2", "3"},
			},
			want: 2,
		},
		{
			name: "正确删除,9,1",
			args: args{
				ctx:     context.TODO(),
				key:     KEY,
				members: []string{"9", "1"},
			},
			want: 2,
		},
		{
			name: "错误删除,19,11",
			args: args{
				ctx:     context.TODO(),
				key:     KEY,
				members: []string{"19", "11"},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.ZRem(tt.args.ctx, tt.args.key, tt.args.members...)
			t.Logf("剩余数据:%+v", s.values[KEY].(*sortedSetValue).rankList)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZRem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ZRem() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sortSetStore_ZRemRangeByRank(t *testing.T) {
	s := newSortSetStore()
	KEY := "ABC"

	member := []SZ{}
	for i := 0; i < 10; i++ {
		member = append(member, SZ{Member: strconv.Itoa(i), Score: float64(i)})
	}
	s.ZAdd(context.TODO(), KEY, member...)
	type args struct {
		ctx   context.Context
		key   string
		start int64
		stop  int64
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "删除成功:0,1",
			args: args{
				ctx:   context.TODO(),
				key:   KEY,
				start: 0,
				stop:  1,
			},
			want: 2,
		},
		{
			name: "删除成功:5,7",
			args: args{
				ctx:   context.TODO(),
				key:   KEY,
				start: 5,
				stop:  7,
			},
			want: 3,
		},
		{
			name: "删除成功:-2,-1",
			args: args{
				ctx:   context.TODO(),
				key:   KEY,
				start: -2,
				stop:  -1,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.ZRemRangeByRank(tt.args.ctx, tt.args.key, tt.args.start, tt.args.stop)
			t.Logf("剩余数据:%+v", s.values[KEY].(*sortedSetValue).rankList)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZRemRangeByRank() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ZRemRangeByRank() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sortedSetStore_ZRemRangeByScore(t *testing.T) {
	s := newSortSetStore()
	KEY := "ABC"

	member := []SZ{}
	for i := 0; i < 10; i++ {
		member = append(member, SZ{Member: strconv.Itoa(i), Score: float64(i)})
	}
	s.ZAdd(context.TODO(), KEY, member...)
	type args struct {
		ctx context.Context
		key string
		min string
		max string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "删除成功,0,1",
			args: args{
				ctx: context.TODO(),
				key: KEY,
				min: "0",
				max: "1",
			},
			want: 2,
		},
		{
			name: "删除成功,(2,3",
			args: args{
				ctx: context.TODO(),
				key: KEY,
				min: "(2",
				max: "3",
			},
			want: 1,
		},
		{
			name: "删除成功,2,(4",
			args: args{
				ctx: context.TODO(),
				key: KEY,
				min: "2",
				max: "(4",
			},
			want: 1,
		},
		{
			name: "删除成功,-inf,4",
			args: args{
				ctx: context.TODO(),
				key: KEY,
				min: "-inf",
				max: "4",
			},
			want: 1,
		},
		{
			name: "删除成功,(6,8",
			args: args{
				ctx: context.TODO(),
				key: KEY,
				min: "(6",
				max: "8",
			},
			want: 2,
		},
		{
			name: "删除成功,(6,+inf",
			args: args{
				ctx: context.TODO(),
				key: KEY,
				min: "(6",
				max: "+inf",
			},
			want: 1,
		},
		{
			name: "删除成功,-inf,+inf",
			args: args{
				ctx: context.TODO(),
				key: KEY,
				min: "-inf",
				max: "+inf",
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.ZRemRangeByScore(tt.args.ctx, tt.args.key, tt.args.min, tt.args.max)
			t.Logf("剩余数据:%+v", s.values[KEY].(*sortedSetValue).rankList)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZRemRangeByScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ZRemRangeByScore() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sortSetStore_ZRevRange(t *testing.T) {
	s := newSortSetStore()
	KEY := "ABC"

	member := []SZ{}
	for i := 0; i < 10; i++ {
		member = append(member, SZ{Member: strconv.Itoa(i), Score: float64(i)})
	}
	s.ZAdd(context.TODO(), KEY, member...)
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
			name: "获取正常,0,-1",
			args: args{
				ctx:   context.TODO(),
				key:   KEY,
				start: 0,
				stop:  -1,
			},
			want: []string{"9", "8", "7", "6", "5", "4", "3", "2", "1", "0"},
		},
		{
			name: "获取正常,5,7",
			args: args{
				ctx:   context.TODO(),
				key:   KEY,
				start: 5,
				stop:  7,
			},
			want: []string{"4", "3", "2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.ZRevRange(tt.args.ctx, tt.args.key, tt.args.start, tt.args.stop)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZRevRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZRevRange() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sortedSetStore_ZRangeByScore(t *testing.T) {
	s := newSortSetStore()
	KEY := "ABC"

	member := []SZ{}
	for i := 0; i < 10; i++ {
		member = append(member, SZ{Member: strconv.Itoa(i), Score: float64(i)})
	}
	s.ZAdd(context.TODO(), KEY, member...)
	type args struct {
		ctx context.Context
		key string
		opt *ZRangeBy
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "获取正常,0,5",
			args: args{
				ctx: context.TODO(),
				key: KEY,
				opt: &ZRangeBy{Min: "0", Max: "5"},
			},
			want: []string{"0", "1", "2", "3", "4", "5"},
		},
		{
			name: "获取正常,(0,5",
			args: args{
				ctx: context.TODO(),
				key: KEY,
				opt: &ZRangeBy{Min: "(0", Max: "5"},
			},
			want: []string{"1", "2", "3", "4", "5"},
		},
		{
			name: "获取正常,(0,(5",
			args: args{
				ctx: context.TODO(),
				key: KEY,
				opt: &ZRangeBy{Min: "(0", Max: "(5"},
			},
			want: []string{"1", "2", "3", "4"},
		},
		{
			name: "获取正常,(0,(5, 2,2",
			args: args{
				ctx: context.TODO(),
				key: KEY,
				opt: &ZRangeBy{Min: "(0", Max: "(5", Offset: 2, Count: 2},
			},
			want: []string{"3", "4"},
		},
		{
			name: "获取正常,(0,(2, 2,2",
			args: args{
				ctx: context.TODO(),
				key: KEY,
				opt: &ZRangeBy{Min: "(0", Max: "(2", Offset: 2, Count: 2},
			},
			want: nil,
		},
		{
			name: "获取正常,0,2, 2,2",
			args: args{
				ctx: context.TODO(),
				key: KEY,
				opt: &ZRangeBy{Min: "0", Max: "2", Offset: 2, Count: 2},
			},
			want: []string{"2"},
		},
		{
			name: "获取正常,0,2, 0,2",
			args: args{
				ctx: context.TODO(),
				key: KEY,
				opt: &ZRangeBy{Min: "0", Max: "2", Offset: 0, Count: 2},
			},
			want: []string{"0", "1"},
		},
		{
			name: "获取正常,5,9, 3,2",
			args: args{
				ctx: context.TODO(),
				key: KEY,
				opt: &ZRangeBy{Min: "5", Max: "9", Offset: 3, Count: 2},
			},
			want: []string{"8", "9"},
		},
		{
			name: "获取正常,-inf,9, 3,2",
			args: args{
				ctx: context.TODO(),
				key: KEY,
				opt: &ZRangeBy{Min: "-inf", Max: "9", Offset: 3, Count: 2},
			},
			want: []string{"3", "4"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.ZRangeByScore(tt.args.ctx, tt.args.key, tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZRangeByScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZRangeByScore() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sortedSetStore_ZCard(t *testing.T) {
	s := newSortSetStore()
	KEY := "ABC"

	member := []SZ{}
	for i := 0; i < 10; i++ {
		member = append(member, SZ{Member: strconv.Itoa(i), Score: float64(i)})
	}
	s.ZAdd(context.TODO(), KEY, member...)

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
			got, err := s.ZCard(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ZCard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sortedSetStore_ZRank(t *testing.T) {
	s := newSortSetStore()
	KEY := "ABC"

	member := []SZ{}
	for i := 0; i < 10; i++ {
		member = append(member, SZ{Member: strconv.Itoa(i), Score: float64(i)})
	}
	s.ZAdd(context.TODO(), KEY, member...)
	type args struct {
		ctx    context.Context
		key    string
		member string
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
			got, err := s.ZRank(tt.args.ctx, tt.args.key, tt.args.member)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZRank() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ZRank() got = %v, want %v", got, tt.want)
			}
		})
	}
}
