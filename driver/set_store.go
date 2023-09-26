package driver

import (
	"context"
	"fmt"
	"github.com/jerbe/jcache/v2/utils"
	"sync"
	"time"
)

/**
  @author : Jerbe - The porter from Earth
  @describe :
*/

type Z struct {
	Score  float64
	Member string
}

type sortSetRankList []*sortSetData

func (l sortSetRankList) String() string {
	var s string
	for i, data := range l {
		s += fmt.Sprintf("{%d:%+v}, ", i, data)
	}
	return fmt.Sprintf("[%s]", s)
}

type sortSetData struct {
	// Key 键
	Key string

	// Score 分数值
	Score float64

	// Rank 排名
	Rank int
}

// sortSetValue 可排序集合值
type sortSetValue struct {
	expireValue

	// rankList 排名顺序列表
	rankList sortSetRankList

	// mapping 字典映射
	mapping map[string]*sortSetData
}

func (v *sortSetValue) setExist(z *Z, data *sortSetData) (b bool) {
	b = false
	newScore := z.Score
	if data.Score == newScore {
		return
	}

	// 如果只有一个元素,就什么都不用变了
	if len(v.rankList) == 1 {
		data.Score = newScore
		return
	}

	// 判断判断位置是否需要更新
	var pre, next *sortSetData
	if data.Rank > 0 {
		pre = v.rankList[data.Rank-1]
	}
	if data.Rank < len(v.rankList)-1 {
		next = v.rankList[data.Rank+1]
	}

	// 如果分数值还在原来位置区间内，则表示位置没有变化
	// 如果newScore大于或者等于next，说明需要调整
	if (pre == nil || pre.Score <= newScore) && (next == nil || newScore < next.Score) {
		data.Score = newScore
		return
	}

	// 先判断是否小于第一个
	first := v.rankList[0]
	if newScore < first.Score {
		for i := data.Rank; i > 0; i-- {
			v.rankList[i-1].Rank++
			v.rankList[i] = v.rankList[i-1]
		}
		data.Rank = 0
		v.rankList[0] = data
		data.Score = newScore
		return
	}

	var current *sortSetData
	// 遍历一下排名列表,找到newScore合适的位置
	for i := 0; i < len(v.rankList); i++ {
		current, next = nil, nil
		current = v.rankList[i]

		if i < len(v.rankList)-1 {
			next = v.rankList[i+1]
		}

		// 如果新分数小于上一个分数
		// 需要符合多种条件
		// 1. 上一个元素分数比当前元素分数小或等于
		// 2. 下一个元素分数
		//    1. 不存在
		//	  2. 分数值比当前元素高
		if current.Score <= newScore && (next == nil || newScore < next.Score) {
			// 如果 data.rank < i
			// 则 data.rank 到 i 之间所有元素的rank需要减1
			if data.Rank < i {
				for j := data.Rank; j < i; j++ {
					v.rankList[j+1].Rank--
					v.rankList[j] = v.rankList[j+1]
				}
				data.Rank = i
				data.Score = newScore
				v.rankList[i] = data
				return
			}

			// 如果 i < data.rank
			// 则 i 到 data.rank之间的所有元素都需要加1
			if i < data.Rank {
				for j := data.Rank; j > i+1; j-- {
					v.rankList[j-1].Rank++
					v.rankList[j] = v.rankList[j-1]
				}
				data.Rank = i + 1
				data.Score = newScore
				v.rankList[i+1] = data
				return
			}

			data.Score = newScore
			return
		}
	}
	return
}

func (v *sortSetValue) setNoExist(z *Z, data *sortSetData) (b bool) {
	b = true
	// 将对象存入映射表中
	v.mapping[data.Key] = data

	newScore := z.Score
	if len(v.rankList) == 0 {
		data.Score = newScore
		data.Rank = 0
		v.rankList = append(v.rankList, data)
		return
	}

	// 省去遍历
	if len(v.rankList) > 0 {
		// 判断是否比第一个元素小,小的话排到第一个元素前面
		first := v.rankList[0]
		if newScore < first.Score {
			data.Score = newScore
			data.Rank = 0
			v.rankList = append(v.rankList, data)

			for i := len(v.rankList) - 1; i > 0; i-- {
				v.rankList[i-1].Rank++
				v.rankList[i] = v.rankList[i-1]
			}

			v.rankList[0] = data
			return
		}

		// 判断是否比最后一个大,大的话排到最后一个元素后面
		last := v.rankList[len(v.rankList)-1]

		if last.Score <= newScore {
			data.Score = newScore
			data.Rank = len(v.rankList)
			v.rankList = append(v.rankList, data)
			return
		}
	}

	// 其余的就需要更新位置了
	var current, next *sortSetData
	for i := 0; i < len(v.rankList); i++ {
		current = v.rankList[i]
		next = nil
		if i < len(v.rankList)-1 {
			next = v.rankList[i+1]
		}

		// 如果新分数小于上一个分数
		// 需要符合多种条件
		// 1. 上一个元素分数比当前元素分数小或等于
		// 2. 下一个元素分数
		//    1. 不存在
		//	  2. 分数值比当前元素高
		if current.Score <= newScore && (next == nil || newScore < next.Score) {
			// 新增一个元素,用于后续方便调整
			v.rankList = append(v.rankList, data)

			// 顶替掉next位置,然后原先位置排名顺延一位
			newRank := i + 1
			for j := len(v.rankList) - 1; j > newRank; j-- {
				v.rankList[j-1].Rank++
				v.rankList[j] = v.rankList[j-1]
			}

			data.Rank = newRank
			data.Score = newScore
			v.rankList[newRank] = data
			return
		}
	}
	return
}

func (v *sortSetValue) Set(z *Z) bool {
	defer func() {
		//log.Printf("最后排序:%+v", v.rankList)
		//log.Printf("最后数量:%d", len(v.rankList))
	}()
	data, ok := v.mapping[z.Member]
	if !ok {
		data = &sortSetData{Key: z.Member}
		return v.setNoExist(z, data)
	}
	return v.setExist(z, data)

}

// newSortSetValue 返回一个新的有序集合数值对象指针
func newSortSetValue() *sortSetValue {
	defaultExpireAt := time.Now().Add(ValueMaxTTL)
	return &sortSetValue{
		expireValue: expireValue{
			expireAt: &defaultExpireAt,
			expired:  false,
		},
		mapping:  make(map[string]*sortSetData),
		rankList: make([]*sortSetData, 0),
	}
}

type sortSetStore struct {
	baseStore
}

func newSortSetStore() *sortSetStore {
	ticker := time.NewTicker(time.Second * 10)
	store := &sortSetStore{
		baseStore: baseStore{
			values:       make(map[string]expireable),
			rwMutex:      sync.RWMutex{},
			expireTicker: ticker,
		},
	}

	go store.checkExpireTick()
	return store
}

func (s *sortSetStore) Type() driverStoreType {
	return driverStoreTypeSortedSet
}

// ZAdd 添加有序集合的元素
func (s *sortSetStore) ZAdd(ctx context.Context, key string, members ...Z) (int64, error) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if err := utils.ContextIsDone(ctx); err != nil {
		return 0, err
	}

	val, ok := s.values[key].(*sortSetValue)
	if !ok {
		val = newSortSetValue()
	}

	// TODO 待优化
	cnt := int64(0)
	for _, member := range members {
		if val.Set(&member) {
			cnt++
		}
	}
	s.values[key] = val

	return cnt, nil
}

// ZCard 获取有序集合的元素数量
func (s *sortSetStore) ZCard(ctx context.Context, key string) (int64, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	if err := utils.ContextIsDone(ctx); err != nil {
		return 0, err
	}

	val, ok := s.values[key].(*sortSetValue)
	if !ok {
		return 0, nil
	}

	return int64(len(val.rankList)), nil
}

func (s *sortSetStore) ZCount(ctx context.Context, key, min, max string) (int64, error) {
	return 0, nil
}

func (s *sortSetStore) ZIncrBy(ctx context.Context, key string, increment float64, member string) (float64, error) {
	return 0, nil
}

func (s *sortSetStore) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	if err := utils.ContextIsDone(ctx); err != nil {
		return []string{}, err
	}

	val, ok := s.values[key].(*sortSetValue)
	if !ok {
		return []string{}, nil
	}

	listLen := int64(len(val.rankList))

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
	//
	//tmp := listLen - stop
	//stop = listLen - start
	//start = tmp

	result := make([]string, stop-start+1, stop-start+1)
	for i := 0; i < len(result); i++ {
		result[i] = val.rankList[int(start)+i].Key
	}
	//sort.Slice(result, func(i, j int) bool {
	//	return true
	//})
	return result, nil

}

func (s *sortSetStore) ZRank(ctx context.Context, key, member string) (int64, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	if err := utils.ContextIsDone(ctx); err != nil {
		return 0, err
	}

	val, ok := s.values[key].(*sortSetValue)
	if !ok {
		return 0, MemoryNil
	}

	data, ok := val.mapping[member]
	if !ok {
		return 0, MemoryNil
	}
	return int64(data.Rank), nil
}

func (s *sortSetStore) ZRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return 0, nil
}
func (s *sortSetStore) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) (int64, error) {
	return 0, nil
}
func (s *sortSetStore) ZRemRangeByScore(ctx context.Context, key, min, max string) (int64, error) {
	return 0, nil
}
func (s *sortSetStore) ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return nil, nil
}
func (s *sortSetStore) ZRevRank(ctx context.Context, key, member string) (int64, error) {
	return 0, nil
}

func (s *sortSetStore) ZScore(ctx context.Context, key, member string) (float64, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	if err := utils.ContextIsDone(ctx); err != nil {
		return 0, err
	}

	val, ok := s.values[key].(*sortSetValue)
	if !ok {
		return 0, MemoryNil
	}

	data, ok := val.mapping[member]
	if !ok {
		return 0, MemoryNil
	}
	return data.Score, nil
}
