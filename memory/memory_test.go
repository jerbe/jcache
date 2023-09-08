package memory

import (
	"context"
	"reflect"
	"testing"
	"time"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/8 22:53
  @describe :
*/

func TestCache_Exists(t *testing.T) {
	type fields struct {
		strStore stringStore
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := &Cache{
				strStore: tt.fields.strStore,
			}
			got, err := mc.Exists(tt.args.ctx, tt.args.key)
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

func TestCache_Get(t *testing.T) {
	type fields struct {
		strStore stringStore
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := &Cache{
				strStore: tt.fields.strStore,
			}
			got, err := mc.Get(tt.args.ctx, tt.args.key)
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

func TestCache_Set(t *testing.T) {
	type fields struct {
		strStore stringStore
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
			mc := &Cache{
				strStore: tt.fields.strStore,
			}
			if err := mc.Set(tt.args.ctx, tt.args.key, tt.args.data, tt.args.expiration); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCache_SetNX(t *testing.T) {
	type fields struct {
		strStore stringStore
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
			mc := &Cache{
				strStore: tt.fields.strStore,
			}
			if err := mc.SetNX(tt.args.ctx, tt.args.key, tt.args.data, tt.args.expiration); (err != nil) != tt.wantErr {
				t.Errorf("SetNX() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCache_checkKeyType(t *testing.T) {
	type fields struct {
		strStore stringStore
	}
	type args struct {
		key   string
		useto StoreType
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := &Cache{
				strStore: tt.fields.strStore,
			}
			got, err := mc.checkKeyType(tt.args.key, tt.args.useto)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkKeyType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkKeyType() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_expireValue_IsExpire(t *testing.T) {
	type fields struct {
		expireAt *time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev := &expireValue{
				expireAt: tt.fields.expireAt,
			}
			if got := ev.IsExpire(); got != tt.want {
				t.Errorf("IsExpire() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_expireValue_SetExpire(t *testing.T) {
	type fields struct {
		expireAt *time.Time
	}
	type args struct {
		d time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev := &expireValue{
				expireAt: tt.fields.expireAt,
			}
			ev.SetExpire(tt.args.d)
		})
	}
}

func Test_expireValue_SetExpireAt(t *testing.T) {
	type fields struct {
		expireAt *time.Time
	}
	type args struct {
		t *time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev := &expireValue{
				expireAt: tt.fields.expireAt,
			}
			ev.SetExpireAt(tt.args.t)
		})
	}
}

func Test_marshalVal(t *testing.T) {
	type args struct {
		data any
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := marshalVal(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("marshalVal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("marshalVal() got = %v, want %v", got, tt.want)
			}
		})
	}
}
