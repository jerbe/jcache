package driver

import (
	"context"
	"errors"
	"sort"
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
	value []string
}

func newListValue() *listValue {
	defaultExpireAt := time.Now().Add(ValueMaxTTL)
	return &listValue{
		expireValue: expireValue{
			expireAt: &defaultExpireAt,
			expired:  false,
		},
		value: make([]string, 0, 1<<8), // 预先进行容量设定,防止append时的扩容耗能
	}
}

type listStore struct {
	baseStore
}

func newListStore() *listStore {
	ticker := time.NewTicker(time.Second * 10)
	store := &listStore{

		baseStore: baseStore{
			values:       make(map[string]expireable),
			rwMutex:      sync.RWMutex{},
			expireTicker: ticker,
		},
	}

	go store.checkExpireTick()
	return store
}

// Push 将数据推入到列表中
func (s *listStore) Push(ctx context.Context, key string, data ...interface{}) (int64, error) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		val, ok := s.values[key].(*listValue)
		if !ok {
			val = newListValue()
		}
		dataLen := len(data)

		if dataLen == 0 {
			return int64(len(val.value)), errors.New("the number of parameters is incorrect")
		}

		var result = make([]string, 0, dataLen)
		for i := 0; i < dataLen; i++ {
			marshal, err := marshalData(data[i])
			if err != nil {
				return 0, err
			}
			result = append(result, marshal)
		}

		val.value = append(val.value, result...)

		s.values[key] = val
		return int64(len(val.value)), nil
	}
}

// Trim 对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除。
// 举个例子，执行命令 LTRIM list 0 2 ，表示只保留列表 list 的前三个元素，其余元素全部删除。
// 下标(index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。
// 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
func (s *listStore) Trim(ctx context.Context, key string, start, stop int64) error {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		val, ok := s.values[key].(*listValue)
		if !ok {
			return nil
		}

		if stop-start < 0 {
			val.value = make([]string, 0)
			return nil
		}

		listLen := int64(len(val.value))
		if listLen == 0 {
			return nil
		}

		// 提取正确的索引位置
		if start < 0 {
			start = listLen + start
		}

		if start < 0 {
			start = 0
		}

		if start >= listLen {
			val.value = make([]string, 0)
			return nil
		}

		if stop < 0 {
			stop = listLen + stop
		}
		if stop < 0 {
			val.value = make([]string, 0)
			return nil
		}

		if stop >= listLen {
			stop = listLen - 1
		}

		start = listLen - start - 1
		stop = listLen - stop
		if start > stop {
			stop, start = start, stop
		}

		result := make([]string, stop-start+1, stop-start+1)
		copy(result, val.value[start-1:stop])
		val.value = result
		return nil
	}
}

// Rang 提取列表范围内的数据
func (s *listStore) Rang(ctx context.Context, key string, start, stop int64) ([]string, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		val, ok := s.values[key].(*listValue)
		if !ok {
			return []string{}, nil
		}

		if stop-start < 0 {
			return []string{}, nil
		}

		listLen := int64(len(val.value))
		if listLen == 0 {
			return []string{}, nil
		}

		// 提取正确的索引位置
		if start < 0 {
			start = listLen + start
		}

		if start < 0 {
			start = 0
		}

		if start >= listLen {
			return []string{}, nil
		}

		if stop < 0 {
			stop = listLen + stop
		}
		if stop < 0 {
			return []string{}, nil
		}

		if stop >= listLen {
			stop = listLen - 1
		}

		start = listLen - start
		stop = listLen - stop
		if start > stop {
			stop, start = start, stop
		}

		result := make([]string, stop-start+1, stop-start+1)
		copy(result, val.value[start-1:stop])
		sort.Slice(result, func(i, j int) bool {
			return true
		})

		return result, nil
	}
}

// Pop 推出列表尾的最后数据
func (s *listStore) Pop(ctx context.Context, key string) (string, error) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		val, ok := s.values[key].(*listValue)
		if !ok {
			return "", MemoryNil
		}

		listLen := len(val.value)

		if listLen == 0 {
			return "", MemoryNil
		}

		back := val.value[0]

		val.value = val.value[1:listLen]
		return back, nil
	}
}

// Shift 推出列表头的第一个数据
func (s *listStore) Shift(ctx context.Context, key string) (string, error) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		val, ok := s.values[key].(*listValue)
		if !ok {
			return "", MemoryNil
		}

		listLen := len(val.value)

		if listLen == 0 {
			return "", MemoryNil
		}

		first := val.value[listLen-1]
		val.value = val.value[0 : listLen-1]
		return first, nil
	}
}
