# JCache - High-Availability Multi-Level Cache Integration Library in Golang [[中文版]](./README_zh.md)


![](https://img.shields.io/github/issues/jerbe/jcache?color=green)
![](https://img.shields.io/github/stars/jerbe/jcache?color=yellow)
![](https://img.shields.io/github/forks/jerbe/jcache?color=orange)
![](https://img.shields.io/github/license/jerbe/jcache?color=ff69b4)
![](https://img.shields.io/badge/language-go-blue)
[![](https://img.shields.io/badge/doc-go-blue)](https://pkg.go.dev/github.com/jerbe/jcache@v1.1.9)
![](https://img.shields.io/github/languages/code-size/jerbe/jcache?color=blueviolet)

## Introduction

* This project is a Golang library that can be used by Golang projects. You can use it by running `go get github.com/jerbe/jcache/v2`.
* This project is a lightweight multi-level cache integration solution inspired by Redis.
* It can run multiple drivers simultaneously without interference.
* Supports the Redis driver or custom drivers, as long as they implement the `driver.Cache` interface.
* Built-in memory cache driver.
* The memory driver supports distribution, based on ETCD service discovery and election strategy. It selects one instance as the master node, and the rest as slave nodes. Every operation on the master node is synchronized to the other slave nodes via a `gRPC` interface, and write operations on slave nodes are first sent to the master node via `gRPC` and then synchronized to the other slave nodes to achieve high availability and data consistency.



## Basic Architecture

At the current stage, functionality is prioritized, and in the future, the priority retrieval order may be specified based on the driver's weight. The current version's priority order is determined by the driver order specified when instantiating the client.
```go
// Instantiate a client with Redis as the preferred driver and memory driver as the fallback.
client := jcache.NewClient(driver.NewRedis(), driver.NewMemory())


// Instantiate a client with memory driver as the preferred driver and Redis as the fallback.
client := jcache.NewClient(driver.NewMemory(), driver.NewRedis())
```
### Basic Architecture Diagram

![](./assets/架构图.jpeg)
## Progress


- [x] Redis driver support
- [x] Native memory driver support
- [x] Standalone mode support
  - [x] Distributed mode support
    - [x] Incremental synchronization of slave nodes
    - [ ] Full synchronization of nodes

## Examples
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
// Instantiate a cache client with memory as the driver
    client := jcache.NewClient()

// Instantiate a distributed memory cache client
    cfg := driver.DistributeMemoryConfig{
		Port: 10080,         // Port for starting the gRPC server, please set different ports for the same machine
        Prefix: "/prefix",   // Set the corresponding prefix according to your business needs
		// EtcdCfg is set according to the configuration of your deployed ETCD service
        EtcdCfg: clientv3.Config{
			Endpoints: []string{"127.0.0.1:2379"}
		}
    }
	client := driver.NewDistributeMemory(cfg)



// Client that supports all operations, including String, Hash, List, etc.
	client := jcache.NewClient(driver.NewMemory())
	client.Set(context.Background(),"hello","world", time.Hour)
	client.Get(context.Background(),"hello")
	client.MGet(context.Background(),"hello","hi")
	...

    // Client that supports only String operations
	stringClient := jcache.NewStringClien(driver.NewMemory()); 
	stringClient.Set(context.Background(),"hello","world", time.Hour)
	stringClient.Get(context.Background(),"hello")
	stringClient.MGet(context.Background(),"hello","hi")
	...

    // Client that supports only Hash operations
	hashClient := jcache.NewHashClient(driver.NewMemory()); 
	hashClient.HSet(context.Background(),"hello","world","boom")
	hashClient.HGet(context.Background(),"hello","world")
	...

// Client that supports only List operations
	listClient := jcache.NewListClient(driver.NewMemory());
	listClient.Push(context.Background(),"hello","world")
	listClient.Pop(context.Background(),"hello")
	listClient.Shift(context.Background(),"hello")
}
```