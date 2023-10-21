package driver

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/jerbe/jcache/v2/driver/proto"

	utils "github.com/jerbe/go-utils"

	"github.com/redis/go-redis/v9"
	v3 "go.etcd.io/etcd/client/v3"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/8 15:58
  @describe :
*/

var MemoryNil = errors.New("memory cache: nil")

var (
	_ Cache = new(Memory)
)

// Memory 内存驱动
type Memory struct {
	rwMutex sync.RWMutex

	storeList []baseStoreer

	ss *stringStore

	hs *hashStore

	ls *listStore

	sts *sortedSetStore

	syncer *memorySyncer
}

/*
内存驱动器
*/

// NewMemory 实例化一个内存核心的缓存驱动
func NewMemory() Cache {
	ss := newStringStore()
	hs := newHashStore()
	ls := newListStore()
	sts := newSortSetStore()

	storeList := []baseStoreer{ss, hs, ls, sts}

	return &Memory{
		rwMutex:   sync.RWMutex{},
		storeList: storeList,
		ss:        ss,
		hs:        hs,
		ls:        ls,
		sts:       sts,
	}
}

// EtcdConfig 使用ETCD服务的配置
type EtcdConfig = v3.Config

// MemoryConfig 内存配置
type MemoryConfig struct {
	// Prefix 业务名前缀,如果用于隔离不同业务
	Prefix string

	// Port 如果打算启用多个驱动,请分别设置多个不冲突的IP用于启动服务
	Port int

	// Username 用户名
	Username string

	// Password 密码
	Password string

	// EtcdConfig 用于启用ETCD的服务
	EtcdConfig EtcdConfig

	// Context 上下文
	Context context.Context
}

// NewMemoryWithConfig 实例化一个分布式的内存核心的缓存驱动
func NewMemoryWithConfig(cfg MemoryConfig) (Cache, error) {
	mem := NewMemory().(*Memory)
	syncer, err := newMemorySyncer(&cfg)
	if err != nil {
		return nil, err
	}
	syncer.setMemory(mem)
	return mem, nil
}

// NewStringMemory 实例化一个仅带字符串存储功能的内存核心缓存驱动
func NewStringMemory() String {
	return NewMemory()
}

// ================================================================================================
// ====================================== PRIVATE ==================================================
// ================================================================================================

// checkKeyExists 检测Key是否已经存在,并放回存储该key的容器
func (m *Memory) checkKeyExists(key string) (baseStoreer, bool) {
	for _, store := range m.storeList {
		if store.KeyExists(key) {
			return store, true
		}
	}
	return nil, false
}

// checkKeyAble 检测Key是否可以使用
func (m *Memory) checkKeyAble(key string, storeType driverStoreType) (bool, error) {
	if s, b := m.checkKeyExists(key); b && s.Type() != storeType {
		return false, fmt.Errorf("the key is already of type '%s' and cannot be set to type'%s'", s.Type(), storeType)
	}
	return true, nil
}

// checkKeysAble 检测 Keys 是否是合法的key
func (m *Memory) checkKeysAble(keys []string, storeType driverStoreType) (bool, error) {
	for _, k := range keys {
		if s, b := m.checkKeyExists(k); b && s.Type() != storeType {
			return false, fmt.Errorf("the key[%s] is already of type '%s' and cannot be set to type'%s'", k, s.Type(), storeType)
		}
	}
	return true, nil
}

// syncToSlave 同步数据到各个终端
func (m *Memory) syncToSlave(action proto.Action, values ...string) {
	// 没有同步器就退出
	if m.syncer == nil {
		return
	}

	// 不是主机的话就退出
	if !m.syncer.isMaster {
		return
	}
	m.syncer.syncToSlaves(action, values...)
}

// syncToMaster 同步数据到主节点
func (m *Memory) syncToMaster(action proto.Action, values ...string) ([]string, error) {
	empty := make([]string, 1)
	if m.syncer == nil {
		return empty, errors.New("Memory: no syncer")
	}

	if m.syncer.isMaster {
		return empty, errors.New("Memory: syncer no a slave node")
	}

	return m.syncer.syncToMaster(action, values...)
}

// ================================================================================================
// ====================================== COMMON ==================================================
// ================================================================================================

// Exists 判断某个Key是否存在
func (m *Memory) Exists(ctx context.Context, keys ...string) IntValuer {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()
	result := &redis.IntCmd{}

	cnt := int64(0)
	for _, store := range m.storeList {
		i, err := store.Exists(ctx, keys...)
		if err != nil {
			result.SetErr(err)
			return result
		}
		cnt += i
	}

	result.SetVal(cnt)
	return result
}

// Del 删除一个或多个key
func (m *Memory) Del(ctx context.Context, keys ...string) IntValuer {

	result := &redis.IntCmd{}
	select {
	case <-ctx.Done():
		result.SetErr(ctx.Err())
		return result
	default:
	}

	if m.syncer == nil || (m.syncer != nil && m.syncer.isMaster) {
		result.SetVal(m.del(context.Background(), keys...))

		if m.syncer != nil {
			m.syncToSlave(proto.Action_Del, keys...)
		}
		return result
	}

	// 同步到主节点
	rsp, err := m.syncToMaster(proto.Action_Del, keys...)
	if err != nil {
		result.SetErr(err)
		return result
	}

	i, err := strconv.ParseInt(rsp[0], 10, 64)
	result.SetErr(err)
	result.SetVal(i)
	return result
}

func (m *Memory) del(ctx context.Context, keys ...string) int64 {
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()
	cnt := int64(0)
	for _, store := range m.storeList {
		i, _ := store.Del(ctx, keys...)
		cnt += i
	}
	return cnt
}

// Expire 设置某个key的存活时间
func (m *Memory) Expire(ctx context.Context, key string, ttl time.Duration) BoolValuer {
	result := &redis.BoolCmd{}
	select {
	case <-ctx.Done():
		result.SetErr(ctx.Err())
		return result
	default:
	}

	// 同步到本地并同步到从节点
	if m.syncer == nil || (m.syncer != nil && m.syncer.isMaster) {
		b, err := m.expire(ctx, key, ttl)
		result.SetVal(b)
		result.SetErr(err)

		if m.syncer != nil {
			dur, _ := marshalData(ttl)
			m.syncToSlave(proto.Action_Expire, key, dur)
		}
		return result
	}

	// 同步到主节点
	dur, _ := marshalData(ttl)
	rsp, err := m.syncToMaster(proto.Action_Expire, key, dur)
	if err != nil {
		result.SetErr(err)
		return result
	}
	if rsp[0] == "1" {
		result.SetVal(true)
	}
	return result
}

func (m *Memory) expire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()
	for _, store := range m.storeList {
		// Todo 错误收集
		b, err := store.Expire(ctx, key, ttl)
		if err != nil {
			return false, err
		}
		if b {
			return b, nil
		}
	}
	return false, nil
}

// ExpireAt 设置某个key在指定时间内到期
func (m *Memory) ExpireAt(ctx context.Context, key string, at time.Time) BoolValuer {
	result := &redis.BoolCmd{}
	select {
	case <-ctx.Done():
		result.SetErr(ctx.Err())
		return result
	default:
	}

	// 设置到本地并同步到从节点
	if m.syncer == nil || (m.syncer != nil && m.syncer.isMaster) {
		b, err := m.expireAt(ctx, key, &at)
		result.SetVal(b)
		result.SetErr(err)

		if m.syncer != nil {
			tm, _ := marshalData(at)
			m.syncToSlave(proto.Action_ExpireAt, key, tm)
		}
		return result
	}

	// 同步到主节点
	tm, _ := marshalData(at)
	rsp, err := m.syncToMaster(proto.Action_ExpireAt, key, tm)
	if err != nil {
		result.SetErr(err)
		return result
	}

	if rsp[0] == "1" {
		result.SetVal(true)
	}
	return result
}

func (m *Memory) expireAt(ctx context.Context, key string, at *time.Time) (bool, error) {
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()
	for _, store := range m.storeList {
		// Todo 错误收集
		b, err := store.ExpireAt(ctx, key, *at)
		if err != nil {
			return false, err
		}

		if b {
			return true, nil
		}
	}
	return false, nil
}

// Persist 删除key的过期时间,并设置成持久性
// 注意,持久化后最长也也不会超过 ValueMaxTTL
func (m *Memory) Persist(ctx context.Context, key string) BoolValuer {
	result := &redis.BoolCmd{}
	// 设置到本地,并同步到从节点
	if m.syncer == nil || (m.syncer != nil && m.syncer.isMaster) {
		b, err := m.persist(ctx, key)
		result.SetErr(err)
		result.SetVal(b)

		if m.syncer != nil {
			m.syncToSlave(proto.Action_Persist, key)
		}

		return result
	}

	// 同步到主节点
	rsp, err := m.syncToMaster(proto.Action_Persist, key)
	if err != nil {
		result.SetErr(err)
		return result
	}

	if rsp[0] == "1" {
		result.SetVal(true)
	}
	return result
}

func (m *Memory) persist(ctx context.Context, key string) (bool, error) {
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()
	for _, store := range m.storeList {
		// Todo 错误收集
		b, err := store.Persist(ctx, key)
		if err != nil {
			return false, err
		}
		if b {
			return b, nil
		}
	}
	return false, nil
}

// ================================================================================================
// ====================================== STRING ==================================================
// ================================================================================================

// Set 设置数据
func (m *Memory) Set(ctx context.Context, key string, data interface{}, expiration time.Duration) StatusValuer {

	val := new(redis.StatusCmd)

	select {
	case <-ctx.Done():
		val.SetErr(ctx.Err())
		return val
	default:
	}

	value, err := marshalData(data)
	if err != nil {
		val.SetErr(err)
		return val
	}
	ttl, _ := marshalData(expiration)
	// 设置到本地并同步到从节点
	if m.syncer == nil || (m.syncer != nil && m.syncer.isMaster) {
		err = m.set(ctx, key, value, expiration)
		if err == nil {
			val.SetVal("OK")
		}
		val.SetErr(err)

		if m.syncer != nil && err == nil {
			m.syncToSlave(proto.Action_Set, key, value, ttl)
		}

		return val
	}
	// 同步到主节点
	rsp, err := m.syncToMaster(proto.Action_Set, key, value, ttl)
	val.SetVal(rsp[0])
	val.SetErr(err)
	return val
}

func (m *Memory) set(ctx context.Context, key, value string, expiration time.Duration) error {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()
	// 检测该Key是否被其他类型用了
	if _, err := m.checkKeyAble(key, driverStoreTypeString); err != nil {
		return err
	}

	return m.ss.Set(ctx, key, value, expiration)
}

// SetNX 设置数据,如果key不存在的话
func (m *Memory) SetNX(ctx context.Context, key string, data interface{}, expiration time.Duration) BoolValuer {

	val := new(redis.BoolCmd)

	select {
	case <-ctx.Done():
		val.SetErr(ctx.Err())
		return val
	default:
	}

	value, err := marshalData(data)
	if err != nil {
		val.SetErr(err)
		return val
	}
	ttl, _ := marshalData(expiration)

	if m.syncer == nil || (m.syncer != nil && m.syncer.isMaster) {
		nx, err := m.setNX(ctx, key, value, expiration)
		val.SetVal(nx)
		val.SetErr(err)

		if m.syncer != nil && err == nil {
			m.syncToSlave(proto.Action_SetNX, key, value, ttl)
		}

		return val
	}

	// 同步到主节点
	rsp, err := m.syncToMaster(proto.Action_SetNX, key, value, ttl)
	val.SetVal(rsp[0] == "1")
	val.SetErr(err)
	return val
}

func (m *Memory) setNX(ctx context.Context, key, data string, expiration time.Duration) (bool, error) {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()
	if _, b := m.checkKeyExists(key); b {
		return false, nil
	}

	return m.ss.SetNX(ctx, key, data, expiration)
}

// Get 获取数据
func (m *Memory) Get(ctx context.Context, key string) StringValuer {
	val := new(redis.StringCmd)
	// 检测该Key是否被其他类型用了
	if _, err := m.checkKeyAble(key, driverStoreTypeString); err != nil {
		val.SetErr(err)
		return val
	}

	v, err := m.ss.Get(ctx, key)
	val.SetVal(v)
	val.SetErr(translateErr(err))
	return val
}

// MGet 获取数据
func (m *Memory) MGet(ctx context.Context, keys ...string) SliceValuer {
	val := new(redis.SliceCmd)
	anies, err := m.ss.MGet(ctx, keys...)
	val.SetVal(anies)
	val.SetErr(err)
	return val
}

// ================================================================================================
// ======================================== HASH ==================================================
// ================================================================================================

// HExists 检测field是否存在哈希表中
func (m *Memory) HExists(ctx context.Context, key, field string) BoolValuer {
	val := new(redis.BoolCmd)
	cnt, err := m.hs.HExists(ctx, key, field)
	val.SetVal(cnt)
	val.SetErr(err)
	return val
}

// HDel 哈希表删除指定字段(fields)
func (m *Memory) HDel(ctx context.Context, key string, fields ...string) IntValuer {

	val := new(redis.IntCmd)

	if err := utils.ContextIsDone(ctx); err != nil {
		val.SetErr(err)
		return val
	}

	values := make([]string, 0, len(fields)+1)
	values = append(values, key)
	values = append(values, fields...)

	// 设置到本地并同步到从节点
	if m.syncer == nil || (m.syncer != nil && m.syncer.isMaster) {
		cnt, err := m.hDel(ctx, key, fields...)
		val.SetVal(cnt)
		val.SetErr(err)
		if m.syncer != nil && err == nil {
			m.syncToSlave(proto.Action_HDel, values...)
		}
		return val
	}

	// 同步到主节点
	rsp, err := m.syncToMaster(proto.Action_HDel, values...)
	if err == nil {
		i, _ := strconv.ParseInt(rsp[0], 10, 64)
		val.SetVal(i)
	}

	val.SetErr(err)
	return val
}

func (m *Memory) hDel(ctx context.Context, key string, fields ...string) (int64, error) {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()
	// 检测该Key是否被其他类型用了
	if _, err := m.checkKeyAble(key, driverStoreTypeHash); err != nil {
		return 0, err
	}

	return m.hs.HDel(ctx, key, fields...)
}

// HSet 哈希表设置数据
func (m *Memory) HSet(ctx context.Context, key string, data ...interface{}) IntValuer {

	val := new(redis.IntCmd)

	if err := utils.ContextIsDone(ctx); err != nil {
		val.SetErr(err)
		return val
	}

	// 将参数列表释放成切片
	dataSlice := sliceArgs(data)
	if len(dataSlice)%2 != 0 {
		val.SetErr(errors.New("the number of parameters is incorrect"))
		return val
	}

	values := make([]string, len(dataSlice)+1, len(dataSlice)+1)
	values[0] = key
	var err error
	for i := 0; i < len(dataSlice); i++ {
		values[i+1], err = marshalData(dataSlice[i])
		if err != nil {
			val.SetErr(err)
			return val
		}
	}

	// 如果没有同步器或者同步器是一个主节点
	// 则可以直接在本地设置数据
	if m.syncer == nil || (m.syncer != nil && m.syncer.isMaster) {
		i, err := m.hSet(ctx, key, values[1:]...)
		val.SetVal(i)
		val.SetErr(err)

		// 同步到从节点
		if m.syncer != nil && err == nil {
			m.syncToSlave(proto.Action_HSet, values...)
		}
		return val
	}

	// 剩下的是从节点操作
	rsp, err := m.syncToMaster(proto.Action_HSet, values...)
	if err == nil {
		i, _ := strconv.ParseInt(rsp[0], 10, 64)
		val.SetVal(i)
	}
	val.SetErr(err)
	return val
}

func (m *Memory) hSet(ctx context.Context, key string, data ...string) (int64, error) {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()
	// 检测该Key是否被其他类型用了
	if _, err := m.checkKeyAble(key, driverStoreTypeHash); err != nil {
		return 0, err
	}

	if len(data)%2 != 0 {
		return 0, errors.New("the number of parameters is incorrect")
	}
	return m.hs.HSet(ctx, key, data...)
}

// HSetNX 如果哈希表的field不存在,则设置成功
func (m *Memory) HSetNX(ctx context.Context, key, field string, data interface{}) BoolValuer {

	val := new(redis.BoolCmd)

	if err := utils.ContextIsDone(ctx); err != nil {
		val.SetErr(err)
		return val
	}

	value, err := marshalData(data)
	if err != nil {
		val.SetErr(err)
		return val
	}

	// 设置到本地并同步到从节点
	if m.syncer == nil || (m.syncer != nil && m.syncer.isMaster) {
		cnt, err := m.hSetNX(ctx, key, field, value)
		val.SetVal(cnt)
		val.SetErr(err)
		if m.syncer != nil && err == nil {
			m.syncToSlave(proto.Action_HSet, key, field, value)
		}
		return val
	}

	// 剩下的是从节点同步到主节点
	rsp, err := m.syncToMaster(proto.Action_HSetNx, key, field, value)
	if err == nil && rsp[0] == "1" {
		val.SetVal(true)
	}
	val.SetErr(err)
	return val
}

func (m *Memory) hSetNX(ctx context.Context, key, field, value string) (bool, error) {
	m.rwMutex.RLock()

	defer m.rwMutex.RUnlock()
	// 检测该Key是否被其他类型用了
	if _, err := m.checkKeyAble(key, driverStoreTypeHash); err != nil {
		return false, err
	}

	return m.hs.HSetNX(ctx, key, field, value)
}

// HGet 哈希表获取一个数据
func (m *Memory) HGet(ctx context.Context, key string, field string) StringValuer {
	val := new(redis.StringCmd)
	v, err := m.hs.HGet(ctx, key, field)

	val.SetVal(v)
	val.SetErr(translateErr(err))
	return val
}

// HMGet 哈希表获取多个数据
func (m *Memory) HMGet(ctx context.Context, key string, fields ...string) SliceValuer {
	val := new(redis.SliceCmd)
	v, err := m.hs.HMGet(ctx, key, fields...)
	val.SetVal(v)
	val.SetErr(err)
	return val
}

// HKeys 哈希表获取某个Key的所有字段(field)
func (m *Memory) HKeys(ctx context.Context, key string) StringSliceValuer {
	val := new(redis.StringSliceCmd)
	v, err := m.hs.HKeys(ctx, key)
	val.SetVal(v)
	val.SetErr(err)
	return val
}

// HVals 哈希表获取所有值
func (m *Memory) HVals(ctx context.Context, key string) StringSliceValuer {
	val := new(redis.StringSliceCmd)
	v, err := m.hs.HVals(ctx, key)
	val.SetVal(v)
	val.SetErr(err)
	return val
}

// HGetAll 获取哈希表所有的数据,包括field跟value
func (m *Memory) HGetAll(ctx context.Context, key string) MapStringStringValuer {
	val := new(redis.MapStringStringCmd)
	v, err := m.hs.HGetAll(ctx, key)
	val.SetVal(v)
	val.SetErr(err)
	return val
}

// HLen 哈希表所有字段的数量
func (m *Memory) HLen(ctx context.Context, key string) IntValuer {
	val := new(redis.IntCmd)
	v, err := m.hs.HLen(ctx, key)
	val.SetVal(v)
	val.SetErr(err)
	return val
}

// ================================================================================================
// ======================================== LIST ==================================================
// ================================================================================================

// LTrim 对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除。
// 举个例子，执行命令 LTRIM list 0 2 ，表示只保留列表 list 的前三个元素，其余元素全部删除。
// 下标(index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。
// 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
func (m *Memory) LTrim(ctx context.Context, key string, start, stop int64) StatusValuer {

	val := new(redis.StatusCmd)

	if err := utils.ContextIsDone(ctx); err != nil {
		val.SetErr(err)
		return val
	}

	startStr, _ := marshalData(start)
	stopStr, _ := marshalData(stop)

	if m.syncer == nil || (m.syncer != nil && m.syncer.isMaster) {
		err := m.lTrim(ctx, key, start, stop)
		if err == nil {
			val.SetVal("OK")
		}
		val.SetErr(err)

		if m.syncer != nil && err == nil {

			m.syncToSlave(proto.Action_LTrim, key, startStr, stopStr)
		}
		return val
	}

	// 剩下的是从节点同步到主节点
	rsp, err := m.syncToMaster(proto.Action_LTrim, key, startStr, stopStr)
	val.SetVal(rsp[0])
	val.SetErr(err)
	return val
}

func (m *Memory) lTrim(ctx context.Context, key string, start, stop int64) error {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()
	// 检测该Key是否被其他类型用了
	if _, err := m.checkKeyAble(key, driverStoreTypeList); err != nil {
		return err
	}
	return m.ls.LTrim(ctx, key, start, stop)
}

// LPush 将数据推入到列表中
func (m *Memory) LPush(ctx context.Context, key string, data ...interface{}) IntValuer {

	val := new(redis.IntCmd)

	if err := utils.ContextIsDone(ctx); err != nil {
		val.SetErr(err)
		return val
	}

	values := make([]string, len(data)+1, len(data)+1)
	values[0] = key
	var err error
	for i := 0; i < len(data); i++ {
		values[i+1], err = marshalData(data[i])
		if err != nil {
			val.SetErr(err)
			return val
		}
	}

	// 设置本地,并同步到从节点
	if m.syncer == nil || (m.syncer != nil && m.syncer.isMaster) {
		cnt, err := m.lPush(ctx, key, values[1:]...)
		val.SetVal(cnt)
		val.SetErr(err)
		if m.syncer != nil && err == nil {
			m.syncToSlave(proto.Action_LPush, values...)
		}
		return val
	}

	// 访问主节点并返回数据
	rsp, err := m.syncToMaster(proto.Action_LPush, values...)
	val.SetErr(err)
	if err == nil {
		cnt, _ := strconv.ParseInt(rsp[0], 10, 64)
		val.SetVal(cnt)
	}
	return val
}

func (m *Memory) lPush(ctx context.Context, key string, values ...string) (int64, error) {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()
	// 检测该Key是否被其他类型用了
	if _, err := m.checkKeyAble(key, driverStoreTypeList); err != nil {
		return 0, err
	}
	return m.ls.LPush(ctx, key, values...)
}

// LRang 提取列表范围内的数据
func (m *Memory) LRang(ctx context.Context, key string, start, stop int64) StringSliceValuer {
	val := new(redis.StringSliceCmd)
	v, err := m.ls.LRang(ctx, key, start, stop)
	val.SetVal(v)
	val.SetErr(translateErr(err))
	return val

}

// LPop 推出列表尾的最后数据
func (m *Memory) LPop(ctx context.Context, key string) StringValuer {

	val := new(redis.StringCmd)

	if err := utils.ContextIsDone(ctx); err != nil {
		val.SetErr(err)
		return val
	}

	// 设置本地,并同步到从节点
	if m.syncer == nil || (m.syncer != nil && m.syncer.isMaster) {
		v, err := m.lPop(ctx, key)
		val.SetVal(v)
		val.SetErr(translateErr(err))

		if m.syncer != nil && err == nil {
			m.syncToSlave(proto.Action_LPop, key)
		}
		return val
	}

	// 访问主节点并返回数据
	rsp, err := m.syncToMaster(proto.Action_LPop, key)
	val.SetVal(rsp[0])
	val.SetErr(err)
	return val
}

func (m *Memory) lPop(ctx context.Context, key string) (string, error) {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()
	// 检测该Key是否被其他类型用了
	if _, err := m.checkKeyAble(key, driverStoreTypeList); err != nil {
		return "", err
	}

	return m.ls.LPop(ctx, key)
}

// LBPop 推出列表尾的最后数据
func (m *Memory) LBPop(ctx context.Context, timeout time.Duration, keys ...string) StringSliceValuer {

	val := new(redis.StringSliceCmd)

	if err := utils.ContextIsDone(ctx); err != nil {
		val.SetErr(err)
		return val
	}

	values := make([]string, 0, len(keys)+1)
	timeoutStr, _ := marshalData(timeout)
	values = append(values, timeoutStr)
	values = append(values, keys...)

	// 设置本地,并同步到从节点
	if m.syncer == nil || (m.syncer != nil && m.syncer.isMaster) {
		v, err := m.lBPop(ctx, timeout, keys...)
		val.SetVal(v)
		val.SetErr(translateErr(err))

		if m.syncer != nil && err == nil {
			for _, k := range keys {
				m.syncToSlave(proto.Action_LPop, k)
			}
		}
		return val
	}

	// 访问主节点并返回数据
	rsp, err := m.syncToMaster(proto.Action_LBPop, values...)
	val.SetVal(rsp)
	val.SetErr(err)
	return val
}

func (m *Memory) lBPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error) {
	m.rwMutex.RLock()
	// 检测该Key是否被其他类型用了
	if _, err := m.checkKeysAble(keys, driverStoreTypeList); err != nil {
		m.rwMutex.RUnlock()
		return nil, err
	}
	m.rwMutex.RUnlock()

	return m.ls.LBPop(ctx, timeout, keys...)
}

// LShift 推出列表头的第一个数据
func (m *Memory) LShift(ctx context.Context, key string) StringValuer {

	val := new(redis.StringCmd)

	if err := utils.ContextIsDone(ctx); err != nil {
		val.SetErr(err)
		return val
	}

	// 设置本地,并同步到从节点
	if m.syncer == nil || (m.syncer != nil && m.syncer.isMaster) {
		v, err := m.lShift(ctx, key)
		val.SetVal(v)
		val.SetErr(translateErr(err))

		if m.syncer != nil && err == nil {
			m.syncToSlave(proto.Action_LShift, key)
		}
		return val
	}

	// 访问主节点并返回数据
	rsp, err := m.syncToMaster(proto.Action_LShift, key)
	val.SetVal(rsp[0])
	val.SetErr(err)
	return val
}

func (m *Memory) lShift(ctx context.Context, key string) (string, error) {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()
	// 检测该Key是否被其他类型用了
	if _, err := m.checkKeyAble(key, driverStoreTypeList); err != nil {
		return "", err
	}

	return m.ls.LShift(ctx, key)
}

// LLen 获取列表长度
func (m *Memory) LLen(ctx context.Context, key string) IntValuer {
	val := new(redis.IntCmd)
	v, err := m.ls.LLen(ctx, key)
	val.SetVal(v)
	val.SetErr(err)
	return val
}

// ==============================================================
// ======================= Sorted Set ===========================
// ==============================================================

// ZAdd 添加有序集合的元素
func (m *Memory) ZAdd(ctx context.Context, key string, members ...Z) IntValuer {

	val := new(redis.IntCmd)

	if err := utils.ContextIsDone(ctx); err != nil {
		val.SetErr(err)
		return val
	}

	values := make([]string, 0, len(members)*2+1)
	values = append(values, key)
	for _, member := range members {
		mb, _ := marshalData(member.Member)
		score, _ := marshalData(member.Score)
		values = append(values, mb, score)
	}

	// 设置本地,并同步到从节点
	if m.syncer == nil || (m.syncer != nil && m.syncer.isMaster) {
		v, err := m.zAdd(ctx, key, values[1:]...)
		val.SetVal(v)
		val.SetErr(translateErr(err))

		if m.syncer != nil && err == nil {
			m.syncToSlave(proto.Action_LShift, values...)
		}
		return val
	}

	// 访问主节点并返回数据
	rsp, err := m.syncToMaster(proto.Action_LShift, values...)
	val.SetErr(err)
	if err == nil {
		cnt, _ := strconv.ParseInt(rsp[0], 10, 64)
		val.SetVal(cnt)
	}
	return val
}

func (m *Memory) zAdd(ctx context.Context, key string, values ...string) (int64, error) {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()
	// 检测该Key是否被其他类型用了
	if _, err := m.checkKeyAble(key, driverStoreTypeSortedSet); err != nil {
		return 0, err
	}

	if len(values)%2 > 1 {
		return 0, errors.New("params number error")
	}

	members := make([]SZ, 0, len(values)/2)
	for i := 0; i < len(values)/2; i++ {
		f, _ := strconv.ParseFloat(values[i*2+1], 64)
		members = append(members, SZ{Member: values[i*2], Score: f})
	}

	return m.sts.ZAdd(ctx, key, members...)
}

// ZCard 获取有序集合的元素数量
func (m *Memory) ZCard(ctx context.Context, key string) IntValuer {
	val := new(redis.IntCmd)
	v, err := m.sts.ZCard(ctx, key)
	val.SetVal(v)
	val.SetErr(translateErr(err))
	return val
}

// ZCount 返回有序集 key 中， score 值在 min 和 max 之间(默认包括 score 值等于 min 或 max )的成员的数量。
func (m *Memory) ZCount(ctx context.Context, key, min, max string) IntValuer {
	val := new(redis.IntCmd)
	v, err := m.sts.ZCount(ctx, key, min, max)
	val.SetVal(v)
	val.SetErr(translateErr(err))
	return val
}

// ZIncrBy 为有序集 key 的成员 member 的 score 值加上增量 increment 。
// 可以通过传递一个负数值 increment ，让 score 减去相应的值，比如 ZINCRBY key -5 member ，就是让 member 的 score 值减去 5
// @return member 成员的新 score 值
func (m *Memory) ZIncrBy(ctx context.Context, key string, increment float64, member string) FloatValuer {

	val := new(redis.FloatCmd)

	if err := utils.ContextIsDone(ctx); err != nil {
		val.SetErr(err)
		return val
	}

	incrementStr, _ := marshalData(increment)

	// 设置本地,并同步到从节点
	if m.syncer == nil || (m.syncer != nil && m.syncer.isMaster) {
		v, err := m.zIncrBy(ctx, key, increment, member)
		val.SetVal(v)
		val.SetErr(translateErr(err))

		if m.syncer != nil && err == nil {
			m.syncToSlave(proto.Action_LShift, key, incrementStr, member)
		}
		return val
	}

	// 访问主节点并返回数据
	rsp, err := m.syncToMaster(proto.Action_LShift, key, incrementStr, member)
	val.SetErr(err)
	if err == nil {
		cnt, _ := strconv.ParseFloat(rsp[0], 64)
		val.SetVal(cnt)
	}
	return val
}

func (m *Memory) zIncrBy(ctx context.Context, key string, increment float64, member string) (float64, error) {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()
	// 检测该Key是否被其他类型用了
	if _, err := m.checkKeyAble(key, driverStoreTypeSortedSet); err != nil {
		return 0, err
	}

	return m.sts.ZIncrBy(ctx, key, increment, member)
}

// ZRange 返回有序集 key 中，指定区间内的成员。
// 其中成员的位置按 score 值递增(从小到大)来排序。
// 具有相同 score 值的成员按字典序(lexicographical order )来排列。
// 如果你需要成员按 score 值递减(从大到小)来排列，请使用 ZREVRANGE 命令。
// 下标参数 start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。
// 你也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。
func (m *Memory) ZRange(ctx context.Context, key string, start, stop int64) StringSliceValuer {
	val := new(redis.StringSliceCmd)
	v, err := m.sts.ZRange(ctx, key, start, stop)
	val.SetVal(v)
	val.SetErr(translateErr(err))
	return val
}

// ZRangeByScore 返回有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。有序集成员按 score 值递增(从小到大)次序排列。
// 具有相同 score 值的成员按字典序(lexicographical order)来排列(该属性是有序集提供的，不需要额外的计算)。
// 可选的 LIMIT 参数指定返回结果的数量及区间(就像SQL中的 SELECT LIMIT offset, count )，注意当 offset 很大时，定位 offset 的操作可能需要遍历整个有序集，此过程最坏复杂度为 O(N) 时间。
func (m *Memory) ZRangeByScore(ctx context.Context, key string, opt *ZRangeBy) StringSliceValuer {
	o := &ZRangeBy{
		Min:    opt.Min,
		Max:    opt.Max,
		Offset: opt.Offset,
		Count:  opt.Count,
	}
	val := new(redis.StringSliceCmd)
	v, err := m.sts.ZRangeByScore(ctx, key, o)
	val.SetVal(v)
	val.SetErr(translateErr(err))
	return val
}

// ZRank 返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递增(从小到大)顺序排列。
// 排名以 0 为底，也就是说， score 值最小的成员排名为 0 。
func (m *Memory) ZRank(ctx context.Context, key, member string) IntValuer {
	val := new(redis.IntCmd)
	v, err := m.sts.ZRank(ctx, key, member)
	val.SetVal(v)
	val.SetErr(translateErr(err))
	return val
}

// ZRem 移除有序集 key 中的一个或多个成员，不存在的成员将被忽略。
// @return 被成功移除的成员的数量，不包括被忽略的成员
func (m *Memory) ZRem(ctx context.Context, key string, members ...interface{}) IntValuer {

	val := new(redis.IntCmd)

	if err := utils.ContextIsDone(ctx); err != nil {
		val.SetErr(err)
		return val
	}

	values := make([]string, 0, len(members)+1)
	values = append(values, key)
	for i := 0; i < len(members); i++ {
		member, err := marshalData(members[i])
		if err != nil {
			val.SetErr(err)
			return val
		}
		values = append(values, member)
	}

	// 设置本地,并同步到从节点
	if m.syncer == nil || (m.syncer != nil && m.syncer.isMaster) {
		v, err := m.zRem(ctx, key, values[1:]...)
		val.SetVal(v)
		val.SetErr(translateErr(err))

		if m.syncer != nil && err == nil {
			m.syncToSlave(proto.Action_LShift, values...)
		}
		return val
	}

	// 访问主节点并返回数据
	rsp, err := m.syncToMaster(proto.Action_LShift, values...)
	val.SetErr(err)
	if err == nil {
		cnt, _ := strconv.ParseInt(rsp[0], 10, 64)
		val.SetVal(cnt)
	}
	return val
}

func (m *Memory) zRem(ctx context.Context, key string, members ...string) (int64, error) {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()
	// 检测该Key是否被其他类型用了
	if _, err := m.checkKeyAble(key, driverStoreTypeSortedSet); err != nil {
		return 0, err
	}

	return m.sts.ZRem(ctx, key, members...)
}

// ZRemRangeByRank 移除有序集 key 中，指定排名(rank)区间内的所有成员。
// 区间分别以下标参数 start 和 stop 指出，包含 start 和 stop 在内。
// 下标参数 start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。
// 你也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。
func (m *Memory) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) IntValuer {

	val := new(redis.IntCmd)

	if err := utils.ContextIsDone(ctx); err != nil {
		val.SetErr(err)
		return val
	}
	startStr, _ := marshalData(start)
	stopStr, _ := marshalData(stop)

	// 设置本地,并同步到从节点
	if m.syncer == nil || (m.syncer != nil && m.syncer.isMaster) {
		v, err := m.zRemRangeByRank(ctx, key, start, stop)
		val.SetVal(v)
		val.SetErr(translateErr(err))

		if m.syncer != nil && err == nil {
			m.syncToSlave(proto.Action_LShift, key, startStr, stopStr)
		}
		return val
	}

	// 访问主节点并返回数据
	rsp, err := m.syncToMaster(proto.Action_LShift, key, startStr, stopStr)
	val.SetErr(err)
	if err == nil {
		cnt, _ := strconv.ParseInt(rsp[0], 10, 64)
		val.SetVal(cnt)
	}
	return val
}

func (m *Memory) zRemRangeByRank(ctx context.Context, key string, start, stop int64) (int64, error) {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()
	// 检测该Key是否被其他类型用了
	if _, err := m.checkKeyAble(key, driverStoreTypeSortedSet); err != nil {
		return 0, err
	}

	return m.sts.ZRemRangeByRank(ctx, key, start, stop)
}

// ZRemRangeByScore 返回有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。
// 有序集成员按 score 值递增(从小到大)次序排列。
func (m *Memory) ZRemRangeByScore(ctx context.Context, key, min, max string) IntValuer {

	val := new(redis.IntCmd)

	if err := utils.ContextIsDone(ctx); err != nil {
		val.SetErr(err)
		return val
	}

	// 设置本地,并同步到从节点
	if m.syncer == nil || (m.syncer != nil && m.syncer.isMaster) {
		v, err := m.zRemRangeByScore(ctx, key, min, max)
		val.SetVal(v)
		val.SetErr(translateErr(err))

		if m.syncer != nil && err == nil {
			m.syncToSlave(proto.Action_LShift, key, min, max)
		}
		return val
	}

	// 访问主节点并返回数据
	rsp, err := m.syncToMaster(proto.Action_LShift, key, min, max)
	val.SetErr(err)
	if err == nil {
		cnt, _ := strconv.ParseInt(rsp[0], 10, 64)
		val.SetVal(cnt)
	}
	return val
}

func (m *Memory) zRemRangeByScore(ctx context.Context, key, min, max string) (int64, error) {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()
	// 检测该Key是否被其他类型用了
	if _, err := m.checkKeyAble(key, driverStoreTypeSortedSet); err != nil {
		return 0, err
	}

	return m.sts.ZRemRangeByScore(ctx, key, min, max)
}

// ZRevRange 返回有序集 key 中，指定区间内的成员。
// 其中成员的位置按 score 值递减(从大到小)来排列。
// 具有相同 score 值的成员按字典序的逆序(reverse lexicographical order)排列。
func (m *Memory) ZRevRange(ctx context.Context, key string, start, stop int64) StringSliceValuer {
	val := new(redis.StringSliceCmd)
	v, err := m.sts.ZRevRange(ctx, key, start, stop)
	val.SetVal(v)
	val.SetErr(translateErr(err))
	return val
}

// ZRevRank 返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递减(从大到小)排序。
// 排名以 0 为底，也就是说， score 值最大的成员排名为 0 。
func (m *Memory) ZRevRank(ctx context.Context, key, member string) IntValuer {
	val := new(redis.IntCmd)
	v, err := m.sts.ZRevRank(ctx, key, member)
	val.SetVal(v)
	val.SetErr(translateErr(err))
	return val
}

// ZScore 返回有序集 key 中，成员 member 的 score 值。
// 如果 member 元素不是有序集 key 的成员，或 key 不存在，返回 nil 。
func (m *Memory) ZScore(ctx context.Context, key, member string) FloatValuer {
	val := new(redis.FloatCmd)
	v, err := m.sts.ZScore(ctx, key, member)
	val.SetVal(v)
	val.SetErr(translateErr(err))
	return val
}
