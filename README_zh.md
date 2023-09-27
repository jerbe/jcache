# JCache - Golang实现的高可用的多级缓存集成方案库


![](https://img.shields.io/github/issues/jerbe/jcache?color=green)
![](https://img.shields.io/github/stars/jerbe/jcache?color=yellow)
![](https://img.shields.io/github/forks/jerbe/jcache?color=orange)
![](https://img.shields.io/github/license/jerbe/jcache?color=ff69b4)
![](https://img.shields.io/badge/language-go-blue)
[![](https://img.shields.io/badge/doc-go-blue)](https://pkg.go.dev/github.com/jerbe/jcache@v1.1.9)
![](https://img.shields.io/github/languages/code-size/jerbe/jcache?color=blueviolet)

## 说明

* 该项目为golang库，golang项目引用该库即可使用。`go get github.com/jerbe/jcache/v2`。
* 此项目是模拟Redis开发的一个轻量的多级缓存集成方案。
* 可同时运行多种驱动，互不干扰。
* 支持Redis驱动，或自定义驱动，只要实现 `driver.Cache` 即可。
* 内置内存缓存驱动。
* 内存驱动支持分布式，基于`ETCD`的服务发现跟选举策略，会选出其中一台实例当做主节点，其余的为从节点。主节点的每次操作都会使用`gRPC`接口同步到其他从节点上；从节点的写操作会使用`gRPC`请求到主节点上再同步到其他从节点上。以尽量达到高可用和数据的一致性。


## 基本架构

现行阶段优先实现功能，未来可能会根据driver的权重指定优先获取顺序。
当前版本的优先顺序按实例化client时指定的driver顺序。
```go
// 实例化一个以redis驱动为优先获取，内存驱动为后取的客户端
client := jcache.NewClient(driver.NewRedis(), driver.NewMemory())


// 实例化一个以内存驱动为优先获取，redis驱动为后取的客户端
client := jcache.NewClient(driver.NewMemory(), driver.NewRedis())
```
### 基本架构图
![](./assets/架构图.jpeg)
## 进度

- [x] Redis驱动支持
- [x] 本机内存驱动支持
  - [x] 单机模式支持
  - [x] 分布式模式支持
    - [x] 从节点的增量同步
    - [ ] 节点的全量同步

## 案例
```shell
  go get github.com/jerbe/jcache/v2
```

```go
import (
    "time"
	
    "github.com/jerbe/jcache/v2"
    "github.com/jerbe/jcache/v2/driver"
)

func main(){
	// 实例化一个以内存作为驱动的缓存客户端
    client := jcache.NewClient()

	// 实例化一个分布式的内存驱动缓存客户端
    cfg := driver.DistributeMemoryConfig{
		Port: 10080,         // 用于启动grpc服务端口,同机器请设置不同端口
        Prefix: "/prefix",   //根据自己的业务设置对应前缀
		// EtcdCfg 根据自己部署的ETCD服务设置对应配置
        EtcdCfg: clientv3.Config{
			Endpoints: []string{"127.0.0.1:2379"}
		}
    }
	client := driver.NewDistributeMemory(cfg)
	
	
	
    // 支持所有操作的客户端,包括 String,Hash,List 
	client := jcache.NewClient(driver.NewMemory())
	client.Set(context.Background(),"hello","world", time.Hour)
	client.Get(context.Background(),"hello")
	client.MGet(context.Background(),"hello","hi")
	...
		
	// 仅支持 String 操作的客户端 
	stringClient := jcache.NewStringClien(driver.NewMemory()); 
	stringClient.Set(context.Background(),"hello","world", time.Hour)
	stringClient.Get(context.Background(),"hello")
	stringClient.MGet(context.Background(),"hello","hi")
	...
	
	// 仅支持 Hash 操作的客户端
	hashClient := jcache.NewHashClient(driver.NewMemory()); 
	hashClient.HSet(context.Background(),"hello","world","boom")
	hashClient.HGet(context.Background(),"hello","world")
	...
	
	// 仅支持 List 操作的客户端 
	listClient := jcache.NewListClient(driver.NewMemory());
	listClient.Push(context.Background(),"hello","world")
	listClient.Pop(context.Background(),"hello")
	listClient.Shift(context.Background(),"hello")
}
```