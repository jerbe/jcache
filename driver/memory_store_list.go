package driver

import (
	"container/list"
	"context"
	"sync"
	"time"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/13 15:47
  @describe :
*/

// listValue 字符串值
type listValue struct {
	expireValue

	// value 值
	value *list.List
}

type listStore struct {
	baseStore
}

func newListStore() *listStore {
	ticker := time.NewTicker(time.Second * 10)
	store := &listStore{

		baseStore: baseStore{
			values:       make(map[string]expirable),
			rwMutex:      sync.RWMutex{},
			expireTicker: ticker,
		},
	}

	go store.checkExpireTick()
	return store
}

// Trim 对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除。
// 举个例子，执行命令 LTRIM list 0 2 ，表示只保留列表 list 的前三个元素，其余元素全部删除。
// 下标(index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。
// 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
func (s *listStore) Trim(ctx context.Context, key string, start, stop int64) error {
	return nil
}

// Push 将数据推入到列表中
func (s *listStore) Push(ctx context.Context, key string, data ...any) (int64, error) {
	return 0, nil
}

// Rang 提取列表范围内的数据
func (s *listStore) Rang(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return nil, nil
}

// Pop 推出列表尾的最后数据
func (s *listStore) Pop(ctx context.Context, key string) (string, error) {
	return "", nil
}

// Shift 推出列表头的第一个数据
func (s *listStore) Shift(ctx context.Context, key string) (string, error) {
	return "", nil
}
