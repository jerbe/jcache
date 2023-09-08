package memory

import (
	"context"
	"log"
	"math/rand"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/8 22:53
  @describe :
*/

func Test_stringStore_Get(t *testing.T) {
	type fields struct {
		rwLock sync.RWMutex
		values map[string]*stringValue
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := &stringStore{
				rwLock: tt.fields.rwLock,
				values: tt.fields.values,
			}
			got, err := ss.Get(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stringStore_MGet(t *testing.T) {
	type fields struct {
		rwLock sync.RWMutex
		values map[string]*stringValue
	}
	type args struct {
		ctx  context.Context
		keys []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    [][]byte
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := &stringStore{
				rwLock: tt.fields.rwLock,
				values: tt.fields.values,
			}
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
	type fields struct {
		rwLock sync.RWMutex
		values map[string]*stringValue
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
		{
			name: "正常设置",
			fields: fields{
				rwLock: sync.RWMutex{},
				values: make(map[string]*stringValue),
			},
			args: args{
				ctx:        context.Background(),
				key:        "test",
				data:       nil,
				expiration: 0,
			},
			wantErr: false,
		},
		{
			name: "再次设置",
			fields: fields{
				rwLock: sync.RWMutex{},
				values: make(map[string]*stringValue),
			},
			args: args{
				ctx:        context.Background(),
				key:        "test",
				data:       "nil",
				expiration: 0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := &stringStore{
				rwLock: tt.fields.rwLock,
				values: tt.fields.values,
			}
			if err := ss.Set(tt.args.ctx, tt.args.key, tt.args.data, tt.args.expiration); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func BenchmarkStringStore_SetAndGet(b *testing.B) {

	b.StopTimer()
	store := newStringStore()
	b.SetParallelism(10000)
	var data = ""
	for i := 0; i < 1000; i++ {
		data += "Redis 命令参考"
	}

	rand.Seed(time.Now().UnixNano())
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var key = "Key:" + strconv.FormatInt(rand.Int63(), 10)

			err := store.Set(context.Background(), key, data, -1)
			if err != nil {
				log.Println(1, err)
			}

			_, err = store.Get(context.Background(), key)
			if err != nil {
				log.Println(2, err)
				return
			}

			//log.Println("3", string(bytes), data)

		}
	})
}

func Test_stringStore_SetNX(t *testing.T) {
	type fields struct {
		rwLock sync.RWMutex
		values map[string]*stringValue
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
			ss := &stringStore{
				rwLock: tt.fields.rwLock,
				values: tt.fields.values,
			}
			if err := ss.SetNX(tt.args.ctx, tt.args.key, tt.args.data, tt.args.expiration); (err != nil) != tt.wantErr {
				t.Errorf("SetNX() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
