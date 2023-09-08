# JCache 一个简单的缓存集成方案


![](https://img.shields.io/github/issues/jerbe/jcache?color=green)
![](https://img.shields.io/github/stars/jerbe/jcache?color=yellow)
![](https://img.shields.io/github/forks/jerbe/jcache?color=orange)
![](https://img.shields.io/github/license/jerbe/jcache?color=ff69b4)
![](https://img.shields.io/github/languages/count/jerbe/jcache)
![](https://img.shields.io/github/languages/code-size/jerbe/jcache?color=blueviolet)
## 项目由来
我们在开发项目中，少不了需要用到缓存，甚至是分布式缓存。我们用的最多的就是Redis，它是一个非常优秀的分布式缓存数据库。

但是如果在生产环境中，Redis挂了，导致某些业务无法再进行，甚至缓存雪崩，导致所有业务都无法进行。

所以一般情况下，当Redis坏掉了，可以再降级使用其他缓存方案，我们利用服务器的内存开发缓存系统是最优的选择.

如果每次开发一个项目都需要写一套缓存系统出来，那得多累人，所以，当前项目就是将缓存操作集成起来，进行了封装，减少重复开发，以免浪费时间。

## 架构
    * 本项目方案采用Redis优先,当Redis无法获取到数据时降级成直接使用本地内存.
    * golang是一个牛逼的语言,使用map可以写出一大堆很优秀的本地缓存框架

## 进度

- [x] Redis缓存支持
- [ ] 其他分布式缓存支持
- [ ] 本机内存支持
  - [x] 支持进行中...

## 问题
    
1. Q:怎么保持数据的一致性?
   A:各个服务之间可以通过服务发现方案建立之间的联系,如果有数据更新,其他订阅的服务实例也能同时接收到数据,并进行更新
2. Q:如果有数据需要更新,是否所有在线的实例都需要更新?
   A:现行方案是这样的，所以可能会导致内存爆棚的问题，这个可能需要更加详细的方案设计。
3. Q:这个方案是独立的服务还是嵌入到代码里面的？
   A:我们这个方案是嵌入到代码里面的，如果应用我们这个项目并使用，那么该应用也会成为一个缓存服务实例。
4. Q:我们的这个方案是最牛逼的吗？
   A:不是的
5. Q:这个库适合所有项目吗?
   A:并不，该库只是进行了简单的封装，其中逻辑较为简单，可能不适合一些较为大型的项目使用。如果需要使用，可以进行`fork`并执行对应的业务修改。
