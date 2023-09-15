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

func newListValue() *listValue {
	defaultExpireAt := time.Now().Add(ValueMaxTTL)
	return &listValue{
		expireValue: expireValue{
			expireAt: &defaultExpireAt,
			expired:  false,
		},
		value: new(list.List),
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

		var result = make([]string, 0, len(data))
		for i := 0; i < len(data); i++ {
			marshal, err := marshalData(data[i])
			if err != nil {
				return 0, err
			}
			result = append(result, marshal)
		}

		for i := 0; i < len(result); i++ {
			val.value.PushFront(result[i])
		}

		s.values[key] = val

		return int64(val.value.Len()), nil
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

		listLen := int64(val.value.Len())
		if listLen == 0 {
			return nil
		}

		// 提取正确的索引位置
		if start < 0 {
			// len = 10, start = -1
			// start = 9
			start = listLen + start
		}

		if stop < 0 {
			stop = listLen + stop
		}

		elem := val.value.Front()
		index := int64(0)

		diffLen := (stop - start) + 1
		if diffLen <= 0 {
			val.value = new(list.List)
			return nil
		}

		// 不在范围内的数据
		for elem != nil {
			next := elem.Next()
			if index < start || index > stop {
				val.value.Remove(elem)
			}
			elem = next
			index++
		}

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

		listLen := int64(val.value.Len())
		if listLen == 0 {
			return []string{}, nil
		}

		// 提取正确的索引位置
		if start < 0 {
			// len = 10, start = -1
			// start = 9
			start = listLen + start
		}

		if stop < 0 {
			stop = listLen + stop
		}

		diffLen := int(stop-start) + 1
		if diffLen < 0 {
			diffLen = 0
		}

		if diffLen == 0 {
			return []string{}, nil
		}

		elem := val.value.Front()
		index := int64(0)
		result := make([]string, 0, diffLen)
		// 丢弃前面部分
		for elem != nil {
			if index >= start && index <= stop {
				result = append(result, elem.Value.(string))
			}
			elem = elem.Next()
			index++
		}
		return result, nil
	}
}

// Pop 推出列表尾的最后数据
func (s *listStore) Pop(ctx context.Context, key string) (string, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		val, ok := s.values[key].(*listValue)
		if !ok {
			return "", MemoryNil
		}

		elem := val.value.Back()
		if elem == nil {
			return "", MemoryNil
		}
		val.value.Remove(elem)
		return elem.Value.(string), nil
	}
}

// Shift 推出列表头的第一个数据
func (s *listStore) Shift(ctx context.Context, key string) (string, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		val, ok := s.values[key].(*listValue)
		if !ok {
			return "", MemoryNil
		}

		elem := val.value.Front()
		if elem == nil {
			return "", MemoryNil
		}
		val.value.Remove(elem)
		return elem.Value.(string), nil
	}
}
