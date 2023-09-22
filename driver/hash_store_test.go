package driver

import (
	"context"
	"reflect"
	"sort"
	"testing"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/15 10:07
  @describe :
*/

func Test_hashStore_HDel(t *testing.T) {
	s := newHashStore()
	s.HSet(context.Background(), "key", "f1", "v1", "f2", "v2")
	s.HSet(context.Background(), "key", "1", "1", "2", "2")
	type args struct {
		ctx    context.Context
		key    string
		fields []string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "删除没有任何数据删除",
			args: args{
				ctx:    context.Background(),
				key:    "key",
				fields: []string{"field1", "field2"},
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "删除了2个数据",
			args: args{
				ctx:    context.Background(),
				key:    "key",
				fields: []string{"f1", "f2"},
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "删除了2个数据",
			args: args{
				ctx:    context.Background(),
				key:    "key",
				fields: []string{"1", "2"},
			},
			want:    2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := s.HDel(tt.args.ctx, tt.args.key, tt.args.fields...)
			if (err != nil) != tt.wantErr {
				t.Errorf("HDel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HDel() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hashStore_HGet(t *testing.T) {
	s := newHashStore()

	s.HSet(context.Background(), "key", "f1", "v1", "f2", "v2", "today", "today", "yesterday", "yesterday")
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
		{
			name: "获取成功",
			args: args{
				ctx:   context.Background(),
				key:   "key",
				field: "f1",
			},
			want:    "v1",
			wantErr: false,
		},
		{
			name: "获取失败,找不到field",
			args: args{
				ctx:   context.Background(),
				key:   "key",
				field: "f3",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "获取失败,找不到key",
			args: args{
				ctx:   context.Background(),
				key:   "key1",
				field: "f2",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "获取获取成功",
			args: args{
				ctx:   context.Background(),
				key:   "key",
				field: "today",
			},
			want:    "today",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.HGet(tt.args.ctx, tt.args.key, tt.args.field)
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

func Test_hashStore_HKeys(t *testing.T) {
	s := newHashStore()

	s.HSet(context.Background(), "key", "f1", "v1", "f2", "v2", "today", "today", "yesterday", "yesterday")

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
		{
			name: "获取成功",
			args: args{
				ctx: context.Background(),
				key: "key",
			},
			want: []string{
				"f1", "f2", "today", "yesterday", // 顺序无解,map遍历是无序的
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := s.HKeys(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("HKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 顺序无解,map遍历是无序的
			sort.Slice(got, func(i, j int) bool {
				return got[i] < got[j]
			})
			sort.Slice(tt.want, func(i, j int) bool {
				return tt.want[i] < tt.want[j]
			})

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HKeys() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hashStore_HLen(t *testing.T) {
	s := newHashStore()

	s.HSet(context.Background(), "key", "f1", "v1", "f2", "v2", "today", "today", "yesterday", "yesterday")

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
			name: "获取成功",
			args: args{
				ctx: context.Background(),
				key: "key",
			},
			want:    4,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.HLen(tt.args.ctx, tt.args.key)
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

func Test_hashStore_HMGet(t *testing.T) {
	s := newHashStore()

	s.HSet(context.Background(), "key", "f1", "v1", "f2", "v2", "today", "today", "yesterday", "yesterday")

	type args struct {
		ctx    context.Context
		key    string
		fields []string
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
				ctx:    context.Background(),
				key:    "key",
				fields: []string{"f1", "f2"},
			},
			want:    []interface{}{"v1", "v2"},
			wantErr: false,
		},
		{
			name: "获取成功",
			args: args{
				ctx:    context.Background(),
				key:    "key",
				fields: []string{"f2", "f1"},
			},
			want:    []interface{}{"v2", "v1"},
			wantErr: false,
		},
		{
			name: "获取成功",
			args: args{
				ctx:    context.Background(),
				key:    "key",
				fields: []string{"f2", "f1", "f1"},
			},
			want:    []interface{}{"v2", "v1", "v1"},
			wantErr: false,
		},
		{
			name: "获取成功",
			args: args{
				ctx:    context.Background(),
				key:    "key",
				fields: []string{"f2", "f1", "f1", "v3"},
			},
			want:    []interface{}{"v2", "v1", "v1", nil},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.HMGet(tt.args.ctx, tt.args.key, tt.args.fields...)
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

func Test_hashStore_HSet(t *testing.T) {
	s := newHashStore()
	type HashData struct {
		Today string `redis:"today"`

		Yesterday string `redis:"yesterday"`

		Tomorrow string `memory:"tomorrow,omitempty"`
	}
	type args struct {
		ctx  context.Context
		key  string
		data []string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "设置成功,字符串",
			args: args{
				ctx:  context.Background(),
				key:  "key",
				data: []string{"f1", "v1", "f2", "v2"},
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "设置成功,字符串数组",
			args: args{
				ctx:  context.Background(),
				key:  "key",
				data: []string{"f11", "v11", "f12", "v12"},
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "设置失败,参数数量不对",
			args: args{
				ctx:  context.Background(),
				key:  "key",
				data: []string{"f11", "v11", "f12"},
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.HSet(tt.args.ctx, tt.args.key, tt.args.data...)
			if (err != nil) != tt.wantErr {
				t.Errorf("HSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HSet() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hashStore_HVals(t *testing.T) {
	s := newHashStore()

	s.HSet(context.Background(), "key", "f1", "v1", "f2", "v2", "f3", "v3", "f4", "v4", "f5", "v5", "f6", "v6", "today", "今天", "yesterday", "昨天")

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
		{
			name: "获取成功",
			args: args{
				ctx: context.Background(),
				key: "key",
			},
			want: []string{
				"v1", "v2", "v3", "v4", "v5", "v6", "今天", "昨天",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.HVals(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("HVals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 顺序无解,map遍历是无序的
			sort.Slice(got, func(i, j int) bool {
				return got[i] < got[j]
			})
			sort.Slice(tt.want, func(i, j int) bool {
				return tt.want[i] < tt.want[j]
			})

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HVals() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hashStore_HGetAll(t *testing.T) {
	s := newHashStore()

	type HashData struct {
		Today string `redis:"today"`

		Yesterday string `redis:"yesterday"`

		Tomorrow string `memory:"tomorrow,omitempty"`
	}

	s.HSet(context.Background(), "key", "f1", "v1", "f2", "v2", "f3", "v3", "f4", "v4", "f5", "v5", "f6", "v6", "today", "今天", "yesterday", "昨天")

	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "获取成功",
			args: args{
				ctx: context.Background(),
				key: "key",
			},
			want: map[string]string{
				"f1":        "v1",
				"f2":        "v2",
				"f3":        "v3",
				"f4":        "v4",
				"f5":        "v5",
				"f6":        "v6",
				"today":     "今天",
				"yesterday": "昨天",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.HGetAll(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("HGetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//
			//// 顺序无解,map遍历是无序的
			//sort.Slice(got, func(i, j int) bool {
			//	return got[i] < got[j]
			//})
			//sort.Slice(tt.want, func(i, j int) bool {
			//	return tt.want[i] < tt.want[j]
			//})

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HGetAll() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hashStore_HSetNX(t *testing.T) {
	s := newHashStore()

	type args struct {
		ctx   context.Context
		key   string
		field string
		data  string
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
				ctx:   context.Background(),
				key:   "key",
				field: "f1",
				data:  "v1",
			},
			want: true,
		},
		{
			name: "设置成功",
			args: args{
				ctx:   context.Background(),
				key:   "key",
				field: "f1",
				data:  "v1",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.HSetNX(tt.args.ctx, tt.args.key, tt.args.field, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("HGetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HGetAll() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hashStore_HExists(t *testing.T) {
	s := newHashStore()
	s.HSet(context.Background(), "key", "f1", "v1", "f2", "v2")

	type args struct {
		ctx   context.Context
		key   string
		field string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "获取成功,f1",
			args: args{
				ctx:   context.Background(),
				key:   "key",
				field: "f1",
			},
			want: true,
		},
		{
			name: "设置成功,f2",
			args: args{
				ctx:   context.Background(),
				key:   "key",
				field: "f2",
			},
			want: true,
		},
		{
			name: "设置失败,f3",
			args: args{
				ctx:   context.Background(),
				key:   "key",
				field: "f3",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.HExists(tt.args.ctx, tt.args.key, tt.args.field)
			if (err != nil) != tt.wantErr {
				t.Errorf("HGetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//
			//// 顺序无解,map遍历是无序的
			//sort.Slice(got, func(i, j int) bool {
			//	return got[i] < got[j]
			//})
			//sort.Slice(tt.want, func(i, j int) bool {
			//	return tt.want[i] < tt.want[j]
			//})

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HGetAll() got = %v, want %v", got, tt.want)
			}
		})
	}
}
