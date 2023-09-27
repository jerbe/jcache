package driver

import (
	"context"
	"fmt"
	"github.com/jerbe/jcache/v2/utils"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

/**
  @author : Jerbe - The porter from Earth
  @describe :
*/

const processCnt = 1000

type Z struct {
	Score  float64
	Member string
}

type sortedSetRankList []*sortedSetData

func (l sortedSetRankList) String() string {
	var s string
	for _, data := range l {
		//s += fmt.Sprintf("{%d:%+v}, ", i, data)
		s += fmt.Sprintf("(%d)%s:%0.2f ", data.Rank, data.Member, data.Score)
		//s += fmt.Sprintf("%s ", data.Member)
	}
	return fmt.Sprintf("[%s]", s)
}

type sortedSetData struct {
	// Member 成员信息
	Member string

	// Score 分数值
	Score float64

	// Rank 排名
	Rank int
}

// sortedSetValue 可排序集合值
type sortedSetValue struct {
	expireValue

	// rankList 排名顺序列表
	rankList sortedSetRankList

	// mapping 字典映射
	mapping map[string]*sortedSetData
}

// Refresh 重新排序,并刷新排行
func (v *sortedSetValue) Refresh() {
	sort.SliceStable(v.rankList, func(i, j int) bool {
		if v.rankList[i].Score < v.rankList[j].Score {
			return true
		}

		if v.rankList[i].Score == v.rankList[j].Score && v.rankList[i].Member < v.rankList[j].Member {
			return true
		}

		return false
	})

	for i, data := range v.rankList {
		data.Rank = i
	}
}

// Set 设置数据
func (v *sortedSetValue) Set(m []Z) int64 {
	//defer func() {
	//log.Printf("%+v", v.rankList)
	//}()
	newCnt := int64(0)
	rankLen := len(v.rankList)
	for _, z := range m {
		if data, ok := v.mapping[z.Member]; ok {
			data.Score = z.Score
		} else {
			data := &sortedSetData{Member: z.Member, Score: z.Score, Rank: rankLen}
			v.mapping[z.Member] = data
			v.rankList = append(v.rankList, data)
			rankLen++
			newCnt++
		}
	}

	// 刷新排行
	v.Refresh()
	return newCnt
}

// newSortSetValue 返回一个新的有序集合数值对象指针
func newSortSetValue() *sortedSetValue {
	defaultExpireAt := time.Now().Add(ValueMaxTTL)
	return &sortedSetValue{
		expireValue: expireValue{
			expireAt: &defaultExpireAt,
			expired:  false,
		},
		mapping:  make(map[string]*sortedSetData),
		rankList: make([]*sortedSetData, 0),
	}
}

type sortedSetStore struct {
	baseStore
}

func newSortSetStore() *sortedSetStore {
	ticker := time.NewTicker(time.Second * 10)
	store := &sortedSetStore{
		baseStore: baseStore{
			values:       make(map[string]expireable),
			rwMutex:      sync.RWMutex{},
			expireTicker: ticker,
		},
	}

	go store.checkExpireTick()
	return store
}

func (s *sortedSetStore) Type() driverStoreType {
	return driverStoreTypeSortedSet
}

// ZAdd 添加有序集合的元素
func (s *sortedSetStore) ZAdd(ctx context.Context, key string, members ...Z) (int64, error) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if err := utils.ContextIsDone(ctx); err != nil {
		return 0, err
	}

	val, ok := s.values[key].(*sortedSetValue)
	if !ok {
		val = newSortSetValue()
	}

	// 使用多个变量一起插入
	cnt := val.Set(members)

	s.values[key] = val

	return cnt, nil
}

// ZCard 获取有序集合的元素数量
func (s *sortedSetStore) ZCard(ctx context.Context, key string) (int64, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	if err := utils.ContextIsDone(ctx); err != nil {
		return 0, err
	}

	val, ok := s.values[key].(*sortedSetValue)
	if !ok {
		return 0, nil
	}

	return int64(len(val.rankList)), nil
}

// ZCount 返回有序集 key 中， score 值在 min 和 max 之间(默认包括 score 值等于 min 或 max )的成员的数量。
func (s *sortedSetStore) ZCount(ctx context.Context, key, min, max string) (int64, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	if err := utils.ContextIsDone(ctx); err != nil {
		return 0, err
	}

	val, ok := s.values[key].(*sortedSetValue)
	if !ok {
		return 0, nil
	}

	if len(val.rankList) == 0 {
		return 0, nil
	}

	minop, maxop := utils.LTE, utils.LTE
	minf, maxf := float64(0), float64(0)
	var err error
	if strings.Index(min, "(") > -1 {
		minop = utils.LT
		min = strings.ReplaceAll(min, "(", "")
	}
	minf, err = strconv.ParseFloat(strings.ReplaceAll(min, "(", ""), 64)
	if err != nil {
		return 0, err
	}

	if strings.Index(max, "(") > -1 {
		maxop = utils.LT
		max = strings.ReplaceAll(max, "(", "")
	}

	maxf, err = strconv.ParseFloat(strings.ReplaceAll(max, "(", ""), 64)
	if err != nil {
		return 0, err
	}

	cnt := int64(0)
	for _, z := range val.rankList {
		if minop(minf, z.Score) && maxop(z.Score, maxf) {
			cnt++
		}
	}

	return cnt, nil
}

// ZIncrBy 为有序集 key 的成员 member 的 score 值加上增量 increment 。
// 可以通过传递一个负数值 increment ，让 score 减去相应的值，比如 ZINCRBY key -5 member ，就是让 member 的 score 值减去 5
// @return member 成员的新 score 值
func (s *sortedSetStore) ZIncrBy(ctx context.Context, key string, increment float64, member string) (float64, error) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if err := utils.ContextIsDone(ctx); err != nil {
		return 0, err
	}

	val, ok := s.values[key].(*sortedSetValue)
	if !ok {
		return 0, nil
	}

	data, ok := val.mapping[member]
	if !ok {
		return 0, MemoryNil
	}
	data.Score += increment
	val.Refresh()

	return data.Score, nil
}

// ZRange 返回有序集 key 中，指定区间内的成员。
// 其中成员的位置按 score 值递增(从小到大)来排序。
// 具有相同 score 值的成员按字典序(lexicographical order )来排列。
// 如果你需要成员按 score 值递减(从大到小)来排列，请使用 ZREVRANGE 命令。
// 下标参数 start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。
// 你也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。
func (s *sortedSetStore) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	if err := utils.ContextIsDone(ctx); err != nil {
		return []string{}, err
	}

	val, ok := s.values[key].(*sortedSetValue)
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
		result[i] = val.rankList[int(start)+i].Member
	}
	//sort.Slice(result, func(i, j int) bool {
	//	return true
	//})
	return result, nil
}

// ZRank 返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递增(从小到大)顺序排列。
// 排名以 0 为底，也就是说， score 值最小的成员排名为 0 。
func (s *sortedSetStore) ZRank(ctx context.Context, key, member string) (int64, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	if err := utils.ContextIsDone(ctx); err != nil {
		return 0, err
	}

	val, ok := s.values[key].(*sortedSetValue)
	if !ok {
		return 0, MemoryNil
	}

	data, ok := val.mapping[member]
	if !ok {
		return 0, MemoryNil
	}
	return int64(data.Rank), nil
}

// ZRem 移除有序集 key 中的一个或多个成员，不存在的成员将被忽略。
// @return 被成功移除的成员的数量，不包括被忽略的成员。
func (s *sortedSetStore) ZRem(ctx context.Context, key string, members ...string) (int64, error) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if err := utils.ContextIsDone(ctx); err != nil {
		return 0, err
	}

	if len(members) == 0 {
		return 0, nil
	}

	val, ok := s.values[key].(*sortedSetValue)
	if !ok {
		return 0, nil
	}

	affectCnt := 0
	// 过滤一遍,将需要删除的对象排名设置成 -1
	// 并将其在映射表中删除
	for _, member := range members {
		if data, ok := val.mapping[member]; ok {
			data.Rank = -1
			affectCnt++
			delete(val.mapping, member)
		}
	}

	// 根据元素排名进行排序
	// 其余的位置是不会改变的
	sort.SliceStable(val.rankList, func(i, j int) bool {
		return val.rankList[i].Rank < val.rankList[j].Rank
	})

	// 丢掉原来的列表,用新列表顶替
	rankLen := len(val.rankList) - affectCnt
	result := make([]*sortedSetData, rankLen, rankLen)
	copy(result, val.rankList[affectCnt:])
	val.rankList = result

	// 重新刷新
	val.Refresh()

	return int64(affectCnt), nil
}

// ZRemRangeByRank 移除有序集 key 中，指定排名(rank)区间内的所有成员。
// 区间分别以下标参数 start 和 stop 指出，包含 start 和 stop 在内。
// 下标参数 start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。
// 你也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。
func (s *sortedSetStore) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) (int64, error) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if err := utils.ContextIsDone(ctx); err != nil {
		return 0, err
	}

	val, ok := s.values[key].(*sortedSetValue)
	if !ok {
		return 0, nil
	}

	affectCnt := 0

	listLen := int64(len(val.rankList))

	if start < 0 {
		start = listLen + start
	}
	if stop < 0 {
		stop = listLen + stop
	}

	if listLen == 0 || stop < start || start >= listLen || stop < 0 {
		return 0, nil
	}

	// 提取正确的索引位置
	if start < 0 {
		start = 0
	}

	if stop >= listLen {
		stop = listLen - 1
	}

	for i := start; i <= stop; i++ {
		data := val.rankList[i]
		delete(val.mapping, data.Member)
	}

	affectCnt = int(stop - start + 1)
	result := make([]*sortedSetData, 0, affectCnt)

	result = append(val.rankList[0:start], val.rankList[stop+1:]...)
	val.rankList = result
	val.Refresh()

	return int64(affectCnt), nil
}

// ZRemRangeByScore 返回有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。
// 有序集成员按 score 值递增(从小到大)次序排列。
func (s *sortedSetStore) ZRemRangeByScore(ctx context.Context, key, min, max string) (int64, error) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if err := utils.ContextIsDone(ctx); err != nil {
		return 0, err
	}

	val, ok := s.values[key].(*sortedSetValue)
	if !ok {
		return 0, nil
	}

	listLen := len(val.rankList)
	if listLen == 0 {
		return 0, nil
	}
	minop, maxop, maxbreakop := utils.LTE, utils.LTE, utils.LT
	minf, maxf := float64(0), float64(0)
	var err error
	if strings.Index(min, "(") > -1 {
		minop = utils.LT
		min = strings.ReplaceAll(min, "(", "")
	}
	minf, err = strconv.ParseFloat(strings.ReplaceAll(min, "(", ""), 64)
	if err != nil {
		return 0, err
	}
	if strings.Index(max, "(") > -1 {
		maxop = utils.LT
		max = strings.ReplaceAll(max, "(", "")
		maxbreakop = utils.LTE
	}
	maxf, err = strconv.ParseFloat(strings.ReplaceAll(max, "(", ""), 64)
	if err != nil {
		return 0, err
	}

	affectCnt := 0

	start := -1
	stop := -1
	for i := 0; i < listLen; i++ {
		data := val.rankList[i]
		stop++
		if minop(minf, data.Score) && maxop(data.Score, maxf) {
			if start == -1 {
				start = i
			}
			delete(val.mapping, data.Member)
		}
		if maxbreakop(maxf, data.Score) {
			break
		}
	}

	affectCnt = stop - start + 1
	result := make([]*sortedSetData, 0, affectCnt)

	result = append(val.rankList[0:start], val.rankList[stop+1:]...)
	val.rankList = result
	val.Refresh()

	return int64(affectCnt), nil
}

// ZRevRange 返回有序集 key 中，指定区间内的成员。
// 其中成员的位置按 score 值递减(从大到小)来排列。
// 具有相同 score 值的成员按字典序的逆序(reverse lexicographical order)排列。
func (s *sortedSetStore) ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	if err := utils.ContextIsDone(ctx); err != nil {
		return []string{}, err
	}

	val, ok := s.values[key].(*sortedSetValue)
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

	// 反相获取
	tmp := listLen - stop
	stop = listLen - start
	start = tmp

	result := make([]string, stop-start+1, stop-start+1)
	for i := 0; i < len(result); i++ {
		result[i] = val.rankList[int(start-1)+i].Member
	}

	sort.Slice(result, func(i, j int) bool {
		return true
	})
	return result, nil
}

// ZRevRank 返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递减(从大到小)排序。
// 排名以 0 为底，也就是说， score 值最大的成员排名为 0 。
func (s *sortedSetStore) ZRevRank(ctx context.Context, key, member string) (int64, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	if err := utils.ContextIsDone(ctx); err != nil {
		return 0, err
	}

	val, ok := s.values[key].(*sortedSetValue)
	if !ok {
		return 0, MemoryNil
	}

	data, ok := val.mapping[member]
	if !ok {
		return 0, MemoryNil
	}

	return int64(val.rankList[len(val.rankList)-1].Rank - data.Rank), nil
}

func (s *sortedSetStore) ZScore(ctx context.Context, key, member string) (float64, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	if err := utils.ContextIsDone(ctx); err != nil {
		return 0, err
	}

	val, ok := s.values[key].(*sortedSetValue)
	if !ok {
		return 0, MemoryNil
	}

	data, ok := val.mapping[member]
	if !ok {
		return 0, MemoryNil
	}
	return data.Score, nil
}
