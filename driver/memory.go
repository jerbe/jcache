package driver

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/jerbe/jcache/v2/driver/proto"
	"github.com/jerbe/jcache/v2/utils"

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

	storeList := []baseStoreer{ss, hs, ls}

	return &Memory{
		rwMutex:   sync.RWMutex{},
		storeList: storeList,
		ss:        ss,
		hs:        hs,
		ls:        ls,
	}
}

type DistributeMemoryConfig struct {
	// Prefix 业务名前缀,如果用于隔离不同业务
	Prefix string

	// Port 如果打算启用多个驱动,请分别设置多个不冲突的IP用于启动服务
	Port int

	// Username 用户名
	Username string

	// Password 密码
	Password string

	// EtcdCfg 用于启用ETCD的服务
	EtcdCfg v3.Config

	// Context 上下文
	Context context.Context
}

// NewDistributeMemory 实例化一个分布式的内存核心的缓存驱动
func NewDistributeMemory(cfg DistributeMemoryConfig) (Cache, error) {
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
func (m *Memory) syncToMaster(action proto.Action, values ...string) (string, error) {
	if m.syncer == nil {
		return "", errors.New("Memory: no syncer")
	}

	if m.syncer.isMaster {
		return "", errors.New("Memory: syncer no master")
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
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()

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

	i, err := strconv.ParseInt(rsp, 10, 64)
	result.SetErr(err)
	result.SetVal(i)
	return result
}

func (m *Memory) del(ctx context.Context, keys ...string) int64 {
	cnt := int64(0)
	for _, store := range m.storeList {
		i, _ := store.Del(ctx, keys...)
		cnt += i
	}
	return cnt
}

// Expire 设置某个key的存活时间
func (m *Memory) Expire(ctx context.Context, key string, ttl time.Duration) BoolValuer {
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()

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
	if rsp == "1" {
		result.SetVal(true)
	}
	return result
}

func (m *Memory) expire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
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
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()

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

	if rsp == "1" {
		result.SetVal(true)
	}
	return result
}

func (m *Memory) expireAt(ctx context.Context, key string, at *time.Time) (bool, error) {
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
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()
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

	if rsp == "1" {
		result.SetVal(true)
	}
	return result
}

func (m *Memory) persist(ctx context.Context, key string) (bool, error) {
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
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()
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
	val.SetVal(rsp)
	val.SetErr(err)
	return val
}

func (m *Memory) set(ctx context.Context, key, value string, expiration time.Duration) error {
	// 检测该Key是否被其他类型用了
	if _, err := m.checkKeyAble(key, driverStoreTypeString); err != nil {
		return err
	}
	return m.ss.Set(ctx, key, value, expiration)
}

// SetNX 设置数据,如果key不存在的话
func (m *Memory) SetNX(ctx context.Context, key string, data interface{}, expiration time.Duration) BoolValuer {
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()
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
	val.SetVal(rsp == "1")
	val.SetErr(err)
	return val
}

func (m *Memory) setNX(ctx context.Context, key, data string, expiration time.Duration) (bool, error) {
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
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()
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
		i, _ := strconv.ParseInt(rsp, 10, 64)
		val.SetVal(i)
	}

	val.SetErr(err)
	return val
}

func (m *Memory) hDel(ctx context.Context, key string, fields ...string) (int64, error) {
	// 检测该Key是否被其他类型用了
	if _, err := m.checkKeyAble(key, driverStoreTypeHash); err != nil {
		return 0, err
	}
	return m.hs.HDel(ctx, key, fields...)
}

// HSet 哈希表设置数据
func (m *Memory) HSet(ctx context.Context, key string, data ...interface{}) IntValuer {
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()
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
		i, _ := strconv.ParseInt(rsp, 10, 64)
		val.SetVal(i)
	}
	val.SetErr(err)
	return val
}

func (m *Memory) hSet(ctx context.Context, key string, data ...string) (int64, error) {
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
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()
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
	if err == nil && rsp == "1" {
		val.SetVal(true)
	}
	val.SetErr(err)
	return val
}

func (m *Memory) hSetNX(ctx context.Context, key, field, value string) (bool, error) {
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
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()
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
	val.SetVal(rsp)
	val.SetErr(err)
	return val
}

func (m *Memory) lTrim(ctx context.Context, key string, start, stop int64) error {
	// 检测该Key是否被其他类型用了
	if _, err := m.checkKeyAble(key, driverStoreTypeList); err != nil {
		return err
	}
	return m.ls.LTrim(ctx, key, start, stop)
}

// LPush 将数据推入到列表中
func (m *Memory) LPush(ctx context.Context, key string, data ...interface{}) IntValuer {
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()

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
		cnt, _ := strconv.ParseInt(rsp, 10, 64)
		val.SetVal(cnt)
	}
	return val
}

func (m *Memory) lPush(ctx context.Context, key string, values ...string) (int64, error) {
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
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()
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
	val.SetVal(rsp)
	val.SetErr(err)
	return val
}

func (m *Memory) lPop(ctx context.Context, key string) (string, error) {
	// 检测该Key是否被其他类型用了
	if _, err := m.checkKeyAble(key, driverStoreTypeList); err != nil {
		return "", err
	}

	return m.ls.LPop(ctx, key)
}

// LShift 推出列表头的第一个数据
func (m *Memory) LShift(ctx context.Context, key string) StringValuer {
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()
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
	val.SetVal(rsp)
	val.SetErr(err)
	return val
}

func (m *Memory) lShift(ctx context.Context, key string) (string, error) {
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
