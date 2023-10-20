package jcache

import (
	"context"
	"github.com/jerbe/jcache/v2/driver"
	"github.com/redis/go-redis/v9"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"reflect"
	"testing"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/29 11:34
  @describe :
*/

func testNewSortedSetClient() *SortedSetClient {
	redisCli := redis.NewClient(&redis.Options{
		Addr:     "192.168.31.101:6379",
		Password: "root",
	})
	redisCli.Del(context.Background(), "ABC")

	opts := &driver.RedisOptions{
		//Config: &driver.RedisConfig{Addrs: []string{"192.168.31.101:6379"}, Password: "root"},
		Client: redisCli,
	}

	redisDriver := driver.NewRedis(opts)
	distributeMemoryCfg := driver.DistributeMemoryConfig{
		Prefix:   "mytest",
		Port:     8088,
		Username: "root",
		Password: "root",
		EtcdCfg:  clientv3.Config{Endpoints: []string{"127.0.0.1:2379"}},
		Context:  nil,
	}
	m, err := driver.NewDistributeMemory(distributeMemoryCfg)
	if err != nil {
		log.Fatal(err)
	}
	return NewSortedSetClient(m, redisDriver)
}

func Test_NewSortedSetClient(t *testing.T) {
	t.Run("测试", func(t *testing.T) {
		cli := testNewSortedSetClient()
		cli.ZAdd(context.Background(), "ABC", driver.Z{Member: "abc", Score: 1})
		v := cli.ZRange(context.TODO(), "ABC", 0, -1)
		t.Logf("数据:%+v,%+v", v.Val(), v.Err())
	})

}

func TestSortedSetClient_ZAdd(t *testing.T) {
	cli := testNewSortedSetClient()
	type args struct {
		ctx     context.Context
		key     string
		members []driver.Z
	}
	tests := []struct {
		name string
		args args
		want driver.IntValuer
	}{
		{
			name: "添加成功,0,1,2",
			args: args{
				ctx: context.Background(),
				key: "ABC",
				members: []driver.Z{
					{
						Score:  0,
						Member: "0",
					},
					{
						Score:  1,
						Member: "1",
					},
					{
						Score:  2,
						Member: "2",
					},
				},
			},
			want: func() driver.IntValuer {
				cmd := &redis.IntCmd{}
				cmd.SetVal(3)
				return cmd
			}(),
		},
		{
			name: "添加成功,0,1,2",
			args: args{
				ctx: context.Background(),
				key: "ABC",
				members: []driver.Z{
					{
						Score:  0,
						Member: "0",
					},
					{
						Score:  1,
						Member: "1",
					},
					{
						Score:  2,
						Member: "2",
					},
				},
			},
			want: func() driver.IntValuer {
				cmd := &redis.IntCmd{}
				cmd.SetVal(0)
				return cmd
			}(),
		},
		{
			name: "添加成功,0,1,2,3",
			args: args{
				ctx: context.Background(),
				key: "ABC",
				members: []driver.Z{
					{
						Score:  0,
						Member: "0",
					},
					{
						Score:  1,
						Member: "1",
					},
					{
						Score:  2,
						Member: "2",
					},
					{
						Score:  3,
						Member: "3",
					},
				},
			},
			want: func() driver.IntValuer {
				cmd := &redis.IntCmd{}
				cmd.SetVal(1)
				return cmd
			}(),
		},
		{
			name: "添加成功,99,98,97,0",
			args: args{
				ctx: context.Background(),
				key: "ABC",
				members: []driver.Z{
					{
						Score:  99,
						Member: "99",
					},
					{
						Score:  98,
						Member: "98",
					},
					{
						Score:  97,
						Member: "97",
					},
					{
						Score:  0,
						Member: "0",
					},
				},
			},
			want: func() driver.IntValuer {
				cmd := &redis.IntCmd{}
				cmd.SetVal(3)
				return cmd
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cli.ZAdd(tt.args.ctx, tt.args.key, tt.args.members...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZAdd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortedSetClient_ZCard(t *testing.T) {
	cli := testNewSortedSetClient()

	members := make([]driver.Z, 0, 10)
	for i := 0; i < 10; i++ {
		members = append(members, driver.Z{Member: i, Score: float64(i)})
	}

	cli.ZAdd(context.Background(), "ABC", members...)

	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name   string
		args   args
		before func()
		want   driver.IntValuer
	}{
		{
			name: "获取正常",
			args: args{
				ctx: context.Background(),
				key: "ABC",
			},
			want: func() driver.IntValuer {
				cmd := &redis.IntCmd{}
				cmd.SetVal(10)
				return cmd
			}(),
		},
		{
			name: "获取正常",
			args: args{
				ctx: context.Background(),
				key: "ABC",
			},
			want: func() driver.IntValuer {
				cmd := &redis.IntCmd{}
				cmd.SetVal(9)
				return cmd
			}(),
			before: func() {
				cli.ZRem(context.Background(), "ABC", "0")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before()
			}
			if got := cli.ZCard(tt.args.ctx, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZCard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortedSetClient_ZCount(t *testing.T) {
	cli := testNewSortedSetClient()
	KEY := "ABC"

	members := make([]driver.Z, 0, 10)
	for i := 0; i < 10; i++ {
		members = append(members, driver.Z{
			Member: i,
			Score:  float64(i),
		})
	}
	cli.ZAdd(context.Background(), KEY, members...)

	type args struct {
		ctx context.Context
		key string
		min string
		max string
	}
	tests := []struct {
		name   string
		before func()
		args   args
		want   driver.IntValuer
	}{
		{
			name: "获取成功,1,4",
			args: args{
				ctx: context.Background(),
				key: KEY,
				min: "1",
				max: "4",
			},
			want: func() driver.IntValuer {
				cmd := &redis.IntCmd{}
				cmd.SetVal(4)
				return cmd
			}(),
		},
		{
			name: "获取成功,(1,4",
			args: args{
				ctx: context.Background(),
				key: KEY,
				min: "(1",
				max: "4",
			},
			want: func() driver.IntValuer {
				cmd := &redis.IntCmd{}
				cmd.SetVal(3)
				return cmd
			}(),
		},
		{
			name: "获取成功,(1,(4",
			args: args{
				ctx: context.Background(),
				key: KEY,
				min: "(1",
				max: "(4",
			},
			want: func() driver.IntValuer {
				cmd := &redis.IntCmd{}
				cmd.SetVal(2)
				return cmd
			}(),
		},
		{
			name: "获取成功,(4,(2",
			args: args{
				ctx: context.Background(),
				key: KEY,
				min: "(4",
				max: "(2",
			},
			want: func() driver.IntValuer {
				cmd := &redis.IntCmd{}
				cmd.SetVal(0)
				return cmd
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cli.ZCount(tt.args.ctx, tt.args.key, tt.args.min, tt.args.max); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortedSetClient_ZIncrBy(t *testing.T) {
	cli := testNewSortedSetClient()
	KEY := "ABC"

	type args struct {
		ctx       context.Context
		key       string
		increment float64
		member    string
	}
	tests := []struct {
		name   string
		before func()
		args   args
		want   driver.FloatValuer
	}{
		{
			name: "递增正常,abc,0",
			args: args{
				ctx:       context.Background(),
				key:       KEY,
				increment: 0,
				member:    "abc",
			},
			want: func() driver.FloatValuer {
				cmd := new(redis.FloatCmd)
				cmd.SetVal(0)
				return cmd
			}(),
		},
		{
			name: "递增正常,abc,1",
			args: args{
				ctx:       context.Background(),
				key:       KEY,
				increment: 1,
				member:    "abc",
			},
			want: func() driver.FloatValuer {
				cmd := new(redis.FloatCmd)
				cmd.SetVal(1)
				return cmd
			}(),
		},
		{
			name: "递增正常,abc,-1",
			args: args{
				ctx:       context.Background(),
				key:       KEY,
				increment: -1,
				member:    "abc",
			},
			want: func() driver.FloatValuer {
				cmd := new(redis.FloatCmd)
				cmd.SetVal(0)
				return cmd
			}(),
		},
		{
			name: "递增正常,abc,-99",
			args: args{
				ctx:       context.Background(),
				key:       KEY,
				increment: -99,
				member:    "abc",
			},
			want: func() driver.FloatValuer {
				cmd := new(redis.FloatCmd)
				cmd.SetVal(-99)
				return cmd
			}(),
		},
		{
			name: "递增正常,abc,99",
			args: args{
				ctx:       context.Background(),
				key:       KEY,
				increment: 99,
				member:    "abc",
			},
			want: func() driver.FloatValuer {
				cmd := new(redis.FloatCmd)
				cmd.SetVal(0)
				return cmd
			}(),
		},
		{
			name: "递增正常,abc,-199",
			args: args{
				ctx:       context.Background(),
				key:       KEY,
				increment: -199,
				member:    "abc",
			},
			want: func() driver.FloatValuer {
				cmd := new(redis.FloatCmd)
				cmd.SetVal(-199)
				return cmd
			}(),
		},
		{
			name: "递增正常,abc,99",
			args: args{
				ctx:       context.Background(),
				key:       KEY,
				increment: 99,
				member:    "abc",
			},
			want: func() driver.FloatValuer {
				cmd := new(redis.FloatCmd)
				cmd.SetVal(-100)
				return cmd
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cli.ZIncrBy(tt.args.ctx, tt.args.key, tt.args.increment, tt.args.member); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZIncrBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortedSetClient_ZRange(t *testing.T) {
	cli := testNewSortedSetClient()
	KEY := "ABC"

	members := make([]driver.Z, 0, 10)
	for i := 0; i < 10; i++ {
		members = append(members, driver.Z{
			Member: i,
			Score:  float64(i),
		})
	}
	cli.ZAdd(context.Background(), KEY, members...)

	type args struct {
		ctx   context.Context
		key   string
		start int64
		stop  int64
	}
	tests := []struct {
		name string
		args args
		want driver.StringSliceValuer
	}{
		{
			name: "获取正常,0,3",
			args: args{
				ctx:   context.Background(),
				key:   KEY,
				start: 0,
				stop:  3,
			},
			want: func() driver.StringSliceValuer {
				cmd := new(redis.StringSliceCmd)
				cmd.SetVal([]string{
					"0", "1", "2", "3",
				})
				return cmd
			}(),
		},
		{
			name: "获取正常,-2,-1",
			args: args{
				ctx:   context.Background(),
				key:   KEY,
				start: -2,
				stop:  -1,
			},
			want: func() driver.StringSliceValuer {
				cmd := new(redis.StringSliceCmd)
				cmd.SetVal([]string{
					"8", "9",
				})
				return cmd
			}(),
		},
		{
			name: "获取正常,0,-1",
			args: args{
				ctx:   context.Background(),
				key:   KEY,
				start: 0,
				stop:  -1,
			},
			want: func() driver.StringSliceValuer {
				cmd := new(redis.StringSliceCmd)
				cmd.SetVal([]string{
					"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
				})
				return cmd
			}(),
		},
		{
			name: "获取正常,-1,0",
			args: args{
				ctx:   context.Background(),
				key:   KEY,
				start: -1,
				stop:  0,
			},
			want: func() driver.StringSliceValuer {
				cmd := new(redis.StringSliceCmd)
				cmd.SetVal([]string{})
				return cmd
			}(),
		},
		{
			name: "获取正常,10,11",
			args: args{
				ctx:   context.Background(),
				key:   KEY,
				start: 10,
				stop:  11,
			},
			want: func() driver.StringSliceValuer {
				cmd := new(redis.StringSliceCmd)
				cmd.SetVal([]string{})
				return cmd
			}(),
		},
		{
			name: "获取正常,9,11",
			args: args{
				ctx:   context.Background(),
				key:   KEY,
				start: 9,
				stop:  11,
			},
			want: func() driver.StringSliceValuer {
				cmd := new(redis.StringSliceCmd)
				cmd.SetVal([]string{
					"9",
				})
				return cmd
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cli.ZRange(tt.args.ctx, tt.args.key, tt.args.start, tt.args.stop); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortedSetClient_ZRangeByScore(t *testing.T) {
	cli := testNewSortedSetClient()
	KEY := "ABC"

	members := make([]driver.Z, 0, 10)
	for i := 0; i < 10; i++ {
		members = append(members, driver.Z{
			Member: i,
			Score:  float64(i),
		})
	}
	cli.ZAdd(context.Background(), KEY, members...)

	type args struct {
		ctx context.Context
		key string
		opt *driver.ZRangeBy
	}
	tests := []struct {
		name string
		args args
		want driver.StringSliceValuer
	}{
		{
			name: "获取成功,0,+inf,5,0",
			args: args{
				ctx: context.Background(),
				key: KEY,
				opt: &driver.ZRangeBy{
					Min: "0", Max: "+inf",
					Offset: 5, Count: 0,
				},
			},
			want: func() driver.StringSliceValuer {
				cmd := new(redis.StringSliceCmd)
				cmd.SetVal([]string{})
				return cmd
			}(),
		},
		{
			name: "获取成功,0,+inf,5,1",
			args: args{
				ctx: context.Background(),
				key: KEY,
				opt: &driver.ZRangeBy{
					Min: "0", Max: "+inf",
					Offset: 5, Count: 1,
				},
			},
			want: func() driver.StringSliceValuer {
				cmd := new(redis.StringSliceCmd)
				cmd.SetVal([]string{
					"5",
				})
				return cmd
			}(),
		},
		{
			name: "获取成功,0,+inf,0,0",
			args: args{
				ctx: context.Background(),
				key: KEY,
				opt: &driver.ZRangeBy{
					Min: "0", Max: "+inf",
					Offset: 0, Count: 0,
				},
			},
			want: func() driver.StringSliceValuer {
				cmd := new(redis.StringSliceCmd)
				cmd.SetVal([]string{
					"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
				})
				return cmd
			}(),
		},
		{
			name: "获取成功,0,+inf,0,5",
			args: args{
				ctx: context.Background(),
				key: KEY,
				opt: &driver.ZRangeBy{
					Min: "0", Max: "+inf",
					Offset: 0, Count: 5,
				},
			},
			want: func() driver.StringSliceValuer {
				cmd := new(redis.StringSliceCmd)
				cmd.SetVal([]string{
					"0", "1", "2", "3", "4",
				})
				return cmd
			}(),
		},
		{
			name: "获取成功,0,+inf,2,5",
			args: args{
				ctx: context.Background(),
				key: KEY,
				opt: &driver.ZRangeBy{
					Min: "0", Max: "+inf",
					Offset: 2, Count: 5,
				},
			},
			want: func() driver.StringSliceValuer {
				cmd := new(redis.StringSliceCmd)
				cmd.SetVal([]string{
					"2", "3", "4", "5", "6",
				})
				return cmd
			}(),
		},
		{
			name: "获取成功,0,+inf,-1,10",
			args: args{
				ctx: context.Background(),
				key: KEY,
				opt: &driver.ZRangeBy{
					Min: "0", Max: "+inf",
					Offset: -1, Count: 10,
				},
			},
			want: func() driver.StringSliceValuer {
				cmd := new(redis.StringSliceCmd)
				cmd.SetVal([]string{})
				return cmd
			}(),
		},
		{
			name: "获取成功,0,+inf,0,-1",
			args: args{
				ctx: context.Background(),
				key: KEY,
				opt: &driver.ZRangeBy{
					Min: "0", Max: "+inf",
					Offset: 0, Count: -1,
				},
			},
			want: func() driver.StringSliceValuer {
				cmd := new(redis.StringSliceCmd)
				cmd.SetVal([]string{
					"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
				})
				return cmd
			}(),
		},
		{
			name: "获取成功,0,+inf,0,1",
			args: args{
				ctx: context.Background(),
				key: KEY,
				opt: &driver.ZRangeBy{
					Min: "0", Max: "+inf",
					Offset: 0, Count: 1,
				},
			},
			want: func() driver.StringSliceValuer {
				cmd := new(redis.StringSliceCmd)
				cmd.SetVal([]string{
					"0",
				})
				return cmd
			}(),
		},
		{
			name: "获取成功,-inf,+inf,0,1",
			args: args{
				ctx: context.Background(),
				key: KEY,
				opt: &driver.ZRangeBy{
					Min: "-inf", Max: "+inf",
					Offset: 0, Count: 1,
				},
			},
			want: func() driver.StringSliceValuer {
				cmd := new(redis.StringSliceCmd)
				cmd.SetVal([]string{
					"0",
				})
				return cmd
			}(),
		},
		{
			name: "获取成功,-inf,(5,0,4",
			args: args{
				ctx: context.Background(),
				key: KEY,
				opt: &driver.ZRangeBy{
					Min: "-inf", Max: "(5",
					Offset: 0, Count: 4,
				},
			},
			want: func() driver.StringSliceValuer {
				cmd := new(redis.StringSliceCmd)
				cmd.SetVal([]string{
					"0", "1", "2", "3",
				})
				return cmd
			}(),
		},
		{
			name: "获取成功,-inf,(5,-1,4",
			args: args{
				ctx: context.Background(),
				key: KEY,
				opt: &driver.ZRangeBy{
					Min: "-inf", Max: "(5",
					Offset: -1, Count: 4,
				},
			},
			want: func() driver.StringSliceValuer {
				cmd := new(redis.StringSliceCmd)
				cmd.SetVal([]string{})
				return cmd
			}(),
		},
		{
			name: "获取成功,-inf,(5,2,4",
			args: args{
				ctx: context.Background(),
				key: KEY,
				opt: &driver.ZRangeBy{
					Min: "-inf", Max: "(5",
					Offset: 2, Count: 4,
				},
			},
			want: func() driver.StringSliceValuer {
				cmd := new(redis.StringSliceCmd)
				cmd.SetVal([]string{
					"2", "3", "4",
				})
				return cmd
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cli.ZRangeByScore(tt.args.ctx, tt.args.key, tt.args.opt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZRangeByScore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortedSetClient_ZRank(t *testing.T) {
	cli := testNewSortedSetClient()
	KEY := "ABC"

	members := make([]driver.Z, 0, 10)
	for i := 0; i < 10; i++ {
		members = append(members, driver.Z{
			Member: i,
			Score:  float64(i),
		})
	}
	cli.ZAdd(context.Background(), KEY, members...)
	type args struct {
		ctx    context.Context
		key    string
		member string
	}
	tests := []struct {
		name string
		args args
		want driver.IntValuer
	}{
		{
			name: "获取正常,1",
			args: args{
				ctx:    context.Background(),
				key:    KEY,
				member: "1",
			},
			want: func() driver.IntValuer {
				cmd := &redis.IntCmd{}
				cmd.SetVal(1)
				return cmd
			}(),
		},
		{
			name: "获取正常,0",
			args: args{
				ctx:    context.Background(),
				key:    KEY,
				member: "0",
			},
			want: func() driver.IntValuer {
				cmd := &redis.IntCmd{}
				cmd.SetVal(0)
				return cmd
			}(),
		},
		{
			name: "获取正常,9",
			args: args{
				ctx:    context.Background(),
				key:    KEY,
				member: "9",
			},
			want: func() driver.IntValuer {
				cmd := &redis.IntCmd{}
				cmd.SetVal(9)
				return cmd
			}(),
		},
		{
			name: "获取失败,-1",
			args: args{
				ctx:    context.Background(),
				key:    KEY,
				member: "-1",
			},
			want: func() driver.IntValuer {
				cmd := &redis.IntCmd{}
				cmd.SetErr(Nil)
				return cmd
			}(),
		},
		{
			name: "获取失败,-99",
			args: args{
				ctx:    context.Background(),
				key:    KEY,
				member: "-99",
			},
			want: func() driver.IntValuer {
				cmd := &redis.IntCmd{}
				cmd.SetErr(Nil)
				return cmd
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cli.ZRank(tt.args.ctx, tt.args.key, tt.args.member); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZRank() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortedSetClient_ZRem(t *testing.T) {
	cli := testNewSortedSetClient()
	KEY := "ABC"

	members := make([]driver.Z, 0, 10)
	for i := 0; i < 10; i++ {
		members = append(members, driver.Z{
			Member: i,
			Score:  float64(i),
		})
	}
	cli.ZAdd(context.Background(), KEY, members...)
	type args struct {
		ctx     context.Context
		key     string
		members []interface{}
	}
	tests := []struct {
		name string
		args args
		want driver.IntValuer
	}{
		{
			name: "删除失败,-1",
			args: args{
				ctx: context.Background(),
				key: KEY,
				members: []interface{}{
					"-1",
				},
			},
			want: func() driver.IntValuer {
				cmd := &redis.IntCmd{}
				return cmd
			}(),
		},
		{
			name: "删除成功,1",
			args: args{
				ctx: context.Background(),
				key: KEY,
				members: []interface{}{
					"1",
				},
			},
			want: func() driver.IntValuer {
				cmd := &redis.IntCmd{}
				cmd.SetVal(1)
				return cmd
			}(),
		},
		{
			name: "再次删除,1",
			args: args{
				ctx: context.Background(),
				key: KEY,
				members: []interface{}{
					"1",
				},
			},
			want: func() driver.IntValuer {
				cmd := &redis.IntCmd{}
				cmd.SetVal(0)
				return cmd
			}(),
		},
		{
			name: "删除成功,3,4,5",
			args: args{
				ctx: context.Background(),
				key: KEY,
				members: []interface{}{
					"3", "4", "5",
				},
			},
			want: func() driver.IntValuer {
				cmd := &redis.IntCmd{}
				cmd.SetVal(3)
				return cmd
			}(),
		},
		{
			name: "删除成功,3,4,5,6",
			args: args{
				ctx: context.Background(),
				key: KEY,
				members: []interface{}{
					"3", "4", "5", "6",
				},
			},
			want: func() driver.IntValuer {
				cmd := &redis.IntCmd{}
				cmd.SetVal(1)
				return cmd
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cli.ZRem(tt.args.ctx, tt.args.key, tt.args.members...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZRem() = %v, want %v", got, tt.want)
			}
		})
	}
}
