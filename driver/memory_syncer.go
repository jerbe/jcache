package driver

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jerbe/jcache/v2/driver/proto"
	"github.com/jerbe/jcache/v2/utils"

	v3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"google.golang.org/grpc"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/22 07:42
  @describe :
*/

const (
	etcdPrefix = "/jcache"
)

var (
	_ proto.SyncerServer = new(syncerServer)
)

type syncerServer struct {
	syncer *memorySyncer
}

func (s *syncerServer) sync(ctx context.Context, in *proto.SyncRequest) (*proto.SyncResponse, error) {
	if err := utils.ContextIsDone(ctx); err != nil {
		return nil, err
	}

	if s.syncer == nil {
		return nil, errors.New("syncerServer: syncer is nil")
	}

	rsp := new(proto.SyncResponse)

	memory := s.syncer.memory
	if memory == nil {
		return nil, errors.New("syncerServer: syncer.memory is nil")
	}

	var err error
	switch in.Action {
	case proto.Action_Del:
		val := memory.del(context.Background(), in.Values...)
		rsp.Value, _ = marshalData(val)
	case proto.Action_Expire:
		var i int64
		i, err = strconv.ParseInt(in.Values[1], 10, 64)
		if err == nil {
			var b bool
			b, err = memory.expire(context.Background(), in.Values[0], time.Duration(i))
			if b {
				rsp.Value = "1"
			}
		}
	case proto.Action_ExpireAt:
		var t time.Time
		t, err = time.Parse(time.RFC3339Nano, in.Values[1])
		if err == nil {
			var b bool
			b, err = memory.expireAt(context.Background(), in.Values[0], &t)
			if b {
				rsp.Value = "1"
			}
		}
	case proto.Action_Persist:
		var b bool
		b, err = memory.persist(context.Background(), in.Values[0])
		if b {
			rsp.Value = "1"
		}
	case proto.Action_Set:
		var i int64
		i, err = strconv.ParseInt(in.Values[2], 10, 64)
		if err == nil {
			err = memory.set(context.Background(), in.Values[0], in.Values[1], time.Duration(i))
			if err == nil {
				rsp.Value = "OK"
			}
		}
	case proto.Action_SetNX:
		var i int64
		i, err = strconv.ParseInt(in.Values[2], 10, 64)
		if err == nil {
			var b bool
			b, err = memory.setNX(context.Background(), in.Values[0], in.Values[1], time.Duration(i))
			if b {
				rsp.Value = "1"
			}
		}
	case proto.Action_HDel:
		var i int64
		i, err = memory.hDel(context.Background(), in.Values[0], in.Values[1:]...)
		if err == nil {
			rsp.Value, _ = marshalData(i)
		}
	case proto.Action_HSet:
		var i int64
		i, err = memory.hSet(context.Background(), in.Values[0], in.Values[1:]...)
		if err == nil {
			rsp.Value, _ = marshalData(i)
		}
	case proto.Action_HSetNx:
		var b bool
		b, err = memory.hSetNX(context.Background(), in.Values[0], in.Values[1], in.Values[2])
		if b {
			rsp.Value = "1"
		}
	case proto.Action_LPush:
		var i int64
		i, err = memory.lPush(context.Background(), in.Values[0], in.Values[1:]...)
		if err == nil {
			rsp.Value, _ = marshalData(i)
		}
	case proto.Action_LPop:
		var v string
		v, err = memory.lPop(context.Background(), in.Values[0])
		if err == nil {
			rsp.Value = v
		}
	case proto.Action_LShift:
		var v string
		v, err = memory.lShift(context.Background(), in.Values[0])
		if err == nil {
			rsp.Value = v
		}
	case proto.Action_LTrim:
		var start, stop int64
		start, err = strconv.ParseInt(in.Values[1], 10, 64)
		stop, err = strconv.ParseInt(in.Values[2], 10, 64)
		err = memory.lTrim(context.Background(), in.Values[0], start, stop)
		if err == nil {
			rsp.Value = "OK"
		}
	default:
		err = errors.New("unknown action")
	}

	if err != nil {
		var statusCode codes.Code
		switch err {
		case MemoryNil:
			statusCode = codes.NotFound
		default:
			statusCode = codes.InvalidArgument
		}
		err = status.New(statusCode, err.Error()).Err()
	}

	return rsp, err
}

// Slave 同步到从节点
func (s *syncerServer) Slave(ctx context.Context, in *proto.SyncRequest) (*proto.SyncResponse, error) {
	var rsp *proto.SyncResponse
	if s.syncer.isMaster {
		rsp = new(proto.SyncResponse)
		err := errors.New("is master")
		return rsp, err
	}
	rsp, err := s.sync(ctx, in)
	return rsp, err
}

// Master 同步到主节点
func (s *syncerServer) Master(ctx context.Context, in *proto.SyncRequest) (*proto.SyncResponse, error) {
	var rsp *proto.SyncResponse
	if !s.syncer.isMaster {
		rsp = new(proto.SyncResponse)
		err := errors.New("not master")
		return rsp, err
	}

	rsp, err := s.sync(ctx, in)
	// 如果是服务端接收到同步数据,需要同步到其他从节点
	if err == nil && s.syncer.isMaster {
		s.syncer.syncToSlaves(in.Action, in.Values...)
	}
	return rsp, err
}

// syncerEndpoint 同步器终端
type syncerEndpoint struct {
	key string
	proto.SyncerClient
	conn     *grpc.ClientConn
	isMaster bool
}

func (e *syncerEndpoint) Close() error {
	return e.conn.Close()
}

// memorySyncer 内存驱动分布式同步器
type memorySyncer struct {
	grpcSvr proto.SyncerServer

	etcdCli *v3.Client

	etcdServerPrefix   string
	etcdElectionPrefix string
	etcdServerKey      string

	isMaster bool

	servername string
	port       int

	endpoints      map[string]*syncerEndpoint
	masterEndpoint *syncerEndpoint

	memory *Memory

	sync.RWMutex

	closed bool
}

func newMemorySyncer(cfg *DistributeMemoryConfig) (*memorySyncer, error) {
	prefix := strings.TrimPrefix(cfg.Prefix, "/")
	port := cfg.Port

	if prefix == "" {
		return nil, errors.New("prefix nil")
	}

	if port <= 0 {
		return nil, errors.New("listen port zero")
	}

	//grpc.
	etcdCli, err := v3.New(cfg.EtcdCfg)
	if err != nil {
		return nil, err
	}

	// 初始化同步服务器
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	grpcSvr := grpc.NewServer()
	svr := &syncerServer{}
	proto.RegisterSyncerServer(grpcSvr, svr)

	// 运行同步服务器
	runGrpcSvr := func() error {
		ch := make(chan error)
		go func() {
			err = grpcSvr.Serve(listen)
			if err != nil {
				ch <- err
				close(ch)
			}
		}()
		time.Sleep(time.Millisecond * 500)
		go func() {
			ch <- nil
		}()
		return <-ch
	}

	err = runGrpcSvr()
	if err != nil {
		return nil, err
	}

	servername := utils.Hostname()
	serverPrefix := fmt.Sprintf("%s/%s/server", etcdPrefix, prefix)
	electionPrefix := fmt.Sprintf("%s/%s/election", etcdPrefix, prefix)

	rand.Seed(time.Now().UnixNano())
	syncer := &memorySyncer{
		grpcSvr:            svr,
		etcdCli:            etcdCli,
		servername:         servername,
		port:               port,
		etcdServerPrefix:   serverPrefix,
		etcdElectionPrefix: electionPrefix,
		etcdServerKey:      fmt.Sprintf("%s/%s/%d", serverPrefix, servername, rand.Int63()),
		endpoints:          make(map[string]*syncerEndpoint),
	}

	svr.syncer = syncer

	var ctx context.Context
	if cfg.Context == nil {
		ctx = context.TODO()
	} else {
		ctx = cfg.Context
	}

	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		if err != nil {
			cancel()
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				syncer.Close()
				return
			}
		}
	}()

	// 进行etcd注册
	err = syncer.register(ctx)
	if err != nil {
		return nil, err
	}

	// 进行etcd监听
	err = syncer.watchingEndpoints(ctx)
	if err != nil {
		return nil, err
	}

	// 进行选主监控
	err = syncer.election(ctx)
	if err != nil {
		return nil, err
	}

	return syncer, nil
}

// tryElection 尝试进行选举,当异常退出选举时有用
func (s *memorySyncer) tryElection(ctx context.Context) {
	var err error
	for err != nil {
		if utils.ContextIsDone(ctx) != nil {
			return
		}

		time.Sleep(time.Second * 5)
		err = s.election(ctx)
		if err == nil {
			return
		}
	}
}

// election 选举
func (s *memorySyncer) election(ctx context.Context) error {
	session, err := concurrency.NewSession(s.etcdCli, concurrency.WithTTL(10))
	if err != nil {
		return err
	}
	election := concurrency.NewElection(session, s.etcdElectionPrefix)
	go func() {
		var r bool
		observerCh := election.Observe(ctx)
		defer func() {
			if o := recover(); o != nil {
				log.Printf("election is panic. reason:[%v]", o)
			}
			session.Close()
			if r && ctx.Err() == nil {
				log.Printf("election was exit, but not closed, register again.  Local:[%s]", s.etcdServerKey)
				go s.tryElection(ctx)
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case <-session.Done():
				return
			case rsp, ok := <-observerCh:
				if ok {
					v := string(rsp.Kvs[0].Value)
					s.Lock()

					// 如果新主终端在目标终端列表中,则需要在列表中删除
					// 并且需要设置主终端
					if e, ok := s.endpoints[v]; ok {
						delete(s.endpoints, v)

						// 如果该key不是本地节点的key,则设置成主节点
						if s.etcdServerKey != v {
							s.masterEndpoint = e
						}
					}

					if v == s.etcdServerKey {
						s.isMaster = true
					}
					s.Unlock()
				} else {
					r = !s.closed
					return
				}
			}
		}
	}()

	errCh := make(chan error)
	go func() {
		e := election.Campaign(ctx, s.etcdServerKey)
		defer func() {
			close(errCh)
		}()
		if e != nil {
			errCh <- e
			session.Close()
			return
		}

		s.isMaster = true
		log.Printf("[Leader] 成为主节点. local:[%s]", s.etcdServerKey)
	}()

	time.Sleep(time.Second)
	select {
	case err = <-errCh:
		return err
	default:
		return nil
	}
}

// tryRegister 尝试注册服务,当异常退出时有用
func (s *memorySyncer) tryRegister(ctx context.Context) {
	var err error
	for err != nil {
		if utils.ContextIsDone(ctx) != nil {
			return
		}

		time.Sleep(time.Second * 5)
		err = s.register(ctx)
		if err == nil {
			return
		}
	}
}

// register 注册服务
func (s *memorySyncer) register(ctx context.Context) error {
	lease, err := s.etcdCli.Grant(ctx, 10)
	if err != nil {
		return err
	}

	kv := v3.NewKV(s.etcdCli)
	_, err = kv.Put(ctx, s.etcdServerKey, fmt.Sprintf("%s:%d", utils.Hostname(), s.port), v3.WithLease(lease.ID))
	if err != nil {
		return err
	}

	var aliveResp <-chan *v3.LeaseKeepAliveResponse
	aliveResp, err = s.etcdCli.KeepAlive(ctx, lease.ID)
	if err != nil {
		return err
	}

	go func(alive <-chan *v3.LeaseKeepAliveResponse) {
		var r bool
		defer func() {
			if o := recover(); o != nil {
				log.Printf("election is panic. reason:[%v]", o)
			}
			if r && ctx.Err() == nil {
				log.Printf("keepalive chan was close, but not closed, register again.  Local:[%s]", s.etcdServerKey)
				go s.tryRegister(ctx)
			}
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case _, ok := <-alive:
				if !ok {
					r = !s.closed
					return
				}
			}
		}
	}(aliveResp)

	return nil
}

// watchingEndpoints 监控终端
func (s *memorySyncer) watchingEndpoints(ctx context.Context) (err error) {
	var response *v3.GetResponse
	response, err = s.etcdCli.Get(ctx, s.etcdServerPrefix, v3.WithPrefix())
	if err != nil {
		return err
	}

	if response.Count > 0 {
		s.Lock()
		for _, kv := range response.Kvs {
			k := string(kv.Key)
			v := string(kv.Value)

			if k == s.etcdServerKey {
				continue
			}

			endpoint, err := newSyncerEndpoint(k, v)
			if err == nil {
				s.endpoints[k] = endpoint
			}
		}
		s.Unlock()
	}

	watchCh := s.etcdCli.Watch(ctx, s.etcdServerPrefix, v3.WithPrefix())
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case w, ok := <-watchCh:
				if !ok {
					log.Println("syncer watchingEndpoints chan was closed")
					return
				}
				s.Lock()
				for _, event := range w.Events {
					k := string(event.Kv.Key)
					v := string(event.Kv.Value)
					log.Printf("[Event] local:[%s], type:[%v], key:[%s], value:[%s]", s.etcdServerKey, event.Type, k, v)
					if event.Type == v3.EventTypePut {
						// 如果Key是本机,则跳过
						if k == s.etcdServerKey {
							continue
						}

						endpoint, err := newSyncerEndpoint(k, v)
						if err != nil {
							log.Printf("new endpoint error. target:[%s]. reason:[%v]", v, err)
							continue
						}
						s.endpoints[k] = endpoint
					}

					if event.Type == v3.EventTypeDelete {
						if endpoint, ok := s.endpoints[k]; ok {
							func() {
								defer func() {
									if obj := recover(); obj != nil {
										log.Printf("close conn has fail. reason:[%v]", obj)
									}
								}()
								endpoint.conn.Close()
							}()
						}
						delete(s.endpoints, k)

						// 移除先前主节点,等待选主完成
						if s.masterEndpoint != nil && s.masterEndpoint.key == k {
							s.masterEndpoint = nil
						}
					}
				}
				s.Unlock()
			}
		}
	}()

	return nil
}

// syncToSlaves 同步数据到从节点
func (s *memorySyncer) syncToSlaves(action proto.Action, values ...string) {
	s.RLock()
	defer s.RUnlock()
	for _, endpoint := range s.endpoints {
		if !endpoint.isMaster {
			req := &proto.SyncRequest{Action: action, Values: values}
			endpoint.SyncerClient.Slave(context.TODO(), req)
		}
	}
}

// syncToMaster 同步数据到主节点
func (s *memorySyncer) syncToMaster(action proto.Action, values ...string) (string, error) {
	s.RLock()
	defer s.RUnlock()
	if !s.isMaster && s.masterEndpoint != nil {
		req := &proto.SyncRequest{Action: action, Values: values}
		rsp, err := s.masterEndpoint.SyncerClient.Master(context.TODO(), req)
		if err != nil {
			statu := status.Convert(err)
			switch statu.Code() {
			case codes.NotFound:
				err = MemoryNil
			default:
				err = errors.New(statu.Message())
			}
			return "", err
		}
		return rsp.Value, nil
	}
	return "", nil
}

func (s *memorySyncer) setMemory(memory *Memory) {
	s.memory = memory
	memory.syncer = s
}

func (s *memorySyncer) Close() error {
	s.closed = true

	if s.etcdCli != nil {
		s.etcdCli.Close()
	}

	if s.masterEndpoint != nil {
		s.masterEndpoint.Close()
	}

	for _, endpoint := range s.endpoints {
		endpoint.Close()
	}
	return nil
}

// newSyncerEndpoint 返回新终端
func newSyncerEndpoint(key, target string) (*syncerEndpoint, error) {
	retry := 3
GRPC_CLI:
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		if retry <= 0 {
			return nil, errors.New("grpc dial max retry")
		}
		retry++
		time.Sleep(time.Second)
		goto GRPC_CLI
	}
	return &syncerEndpoint{SyncerClient: proto.NewSyncerClient(conn), conn: conn, key: key}, nil
}
