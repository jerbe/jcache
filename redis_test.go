package jcache

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/9 11:53
  @describe :
*/

type strSlice []string

func (ss strSlice) MarshalBinary() ([]byte, error) {
	return json.Marshal(ss)
}

func Test_Redis(t *testing.T) {
	cli, err := initRedis(&RedisConfig{
		Mode:       "single",
		MasterName: "",
		Addrs:      []string{"192.168.31.101:6379"},
		Database:   "",
		Username:   "",
		Password:   "root",
	})
	if err != nil {
		t.Fatal(err)
	}

	var q = []int{1, 2, 3, 4, 5}
	c := make([]int, 0)
	copy(c, q)
	c[0] = 2
	fmt.Println(q)
	fmt.Println(c)

	cli.Set(context.Background(), "key1", "value1", time.Minute)
	cli.Set(context.Background(), "key2", "value2", time.Minute)
	cli.HSet(context.Background(), "key3", "field4", "value4")
	cli.Set(context.Background(), "key4", "value4", time.Minute)
	fmt.Println(cli.Set(context.Background(), "key5", strSlice{"1", "2", "3", "4"}, time.Minute).Err())

	cmd := cli.MGet(context.Background(), "key1", "key2", "key3", "key4", "key5")
	fmt.Println(cmd.Val())
}

func Test_String(t *testing.T) {
	var data any
	data = []byte("123456")
	d, ok := data.(string)
	fmt.Println(d, ok)

}
