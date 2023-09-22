package driver

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/jerbe/go-errors"
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

func (s *listStore) Type() driverStoreType {
	return driverStoreTypeList
}

// LPush 将数据推入到列表中
// 推入后列表顺序,先推入在左,后推入在右 [a,b,c,d,e...]
func (s *listStore) LPush(ctx context.Context, key string, data ...string) (int64, error) {
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

		val.value = append(val.value, data...)
		s.values[key] = val
		return int64(len(val.value)), nil
	}
}

// LTrim 对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除。
// 举个例子，执行命令 LTRIM list 0 2 ，表示只保留列表 list 的前三个元素，其余元素全部删除。
// 下标(index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。
// 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
// 列表顺序 先推入在左,后推入在右 [a,b,c,d,e]
// 如果 start = 0, stop = 1, 裁剪后应该保留 [d,e]
// 如果 start = -2, stop = -1, 裁剪后应该保留 [a,b]
func (s *listStore) LTrim(ctx context.Context, key string, start, stop int64) error {
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
		listLen := int64(len(val.value))
		if listLen == 0 {
			return nil
		}

		if start < 0 {
			start = listLen + start
		}
		if stop < 0 {
			stop = listLen + stop
		}

		if stop < start || start >= listLen || stop < 0 {
			val.value = make([]string, 0)
			return nil
		}

		// 提取正确的索引位置
		if start < 0 {
			start = 0
		}

		if stop >= listLen {
			stop = listLen - 1
		}

		tmp := listLen - stop
		stop = listLen - start
		start = tmp

		result := make([]string, stop-start+1, stop-start+1)
		if len(result) == 0 {
			val.value = nil
			delete(s.values, key)
			return nil
		}
		copy(result, val.value[start-1:stop])
		val.value = result
		return nil
	}
}

// LRang 提取列表范围内的数据
// 列表顺序 先推入在左,后推入在右 [a,b,c,d,e]
// 如果 start = 0, stop = 1, 取得范围应该 [e,d]
// 如果 start = -2, stop = -1, 取得范围应该 [b,a]
func (s *listStore) LRang(ctx context.Context, key string, start, stop int64) ([]string, error) {
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

		listLen := int64(len(val.value))

		if start < 0 {
			start = listLen + start
		}
		if stop < 0 {
			stop = listLen + stop
		}

		if listLen == 0 || stop < start || start >= listLen || stop < 0 {
			return []string{}, nil
		}

		// 提取正确的索引位置
		if start < 0 {
			start = 0
		}

		if stop >= listLen {
			stop = listLen - 1
		}

		tmp := listLen - stop
		stop = listLen - start
		start = tmp

		result := make([]string, stop-start+1, stop-start+1)

		copy(result, val.value[start-1:stop])
		sort.Slice(result, func(i, j int) bool {
			return true
		})

		return result, nil
	}
}

// LPop 推出列表尾的最后数据
// 列表顺序 先推入在左,后推入在右 [a,b,c,d,e]
// LPop得到的数据应该是'a'
func (s *listStore) LPop(ctx context.Context, key string) (string, error) {
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

		item := val.value[0]
		if len(val.value) == 1 {
			val.value = nil
			delete(s.values, key)
		} else {
			val.value = val.value[1:listLen:listLen]
		}
		return item, nil
	}
}

// LShift 推出列表头的第一个数据
// 列表顺序 先推入在左,后推入在右 [a,b,c,d,e]
// LShift得到的数据应该是'e'
func (s *listStore) LShift(ctx context.Context, key string) (string, error) {
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

		item := val.value[listLen-1]
		if len(val.value) == 1 {
			val.value = nil
			delete(s.values, key)
		} else {
			val.value = val.value[: listLen-1 : listLen-1]
		}

		return item, nil
	}
}

// LLen 列表长度
func (s *listStore) LLen(ctx context.Context, key string) (int64, error) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		val, ok := s.values[key].(*listValue)
		if !ok {
			return 0, nil
		}

		listLen := len(val.value)
		return int64(listLen), nil
	}
}
