> 在open-falcon监控系统中，agent采集到数据后就会将数据上报的 transfer 组件，那transfer都做了什么呢？带着这个问题，我们一起学习transfer组件，希望对你有帮助。

## 一、 transfer 启动的时候都做了什么

同样我们从main函数入手，下面时main函数的执行流程：

```mermaid
graph LR
A[加载配置文件] --> B[接收命令行参数]
B --> C[启动proc]
C --> D[开启接收器]
D --> E[开启发送器]
E --> F[开启http服务]
```

#### 1. 加载配置文件

```js
{
    "debug": true,                   // 如果为true，日志中会打印debug信息
    "minStep": 30,                   // //最小上报周期,单位sec
    "http": {
        "enabled": true,             // 表示是否开启该http端口，该端口为控制端口，主要用来对transfer发送控制命令、统计命令、debug命令等
        "listen": "0.0.0.0:6060"     // 表示监听的http端口
    },
    "rpc": {
        "enabled": true,             // 是否开启rpc服务
        "listen": "0.0.0.0:8433"     // rpc监听地址
    },
    "socket": {
        "enabled": true,             // 是否开启socket服务
        "listen": "0.0.0.0:4444",    // socket服务监听地址
        "timeout": 3600              // 连接超时时间
    },
    "judge": {
        "enabled": true,             // 表示是否开启向judge发送数据
        "batch": 200,                // 数据转发的批量大小，可以加快发送速度，建议保持默认值
        "connTimeout": 1000,         // 位是毫秒，与后端建立连接的超时时间，可以根据网络质量微调，建议保持默认
        "callTimeout": 5000,         // 单位是毫秒，发送数据给后端的超时时间，可以根据网络质量微调，建议保持默认
        "maxConns": 32,              // 连接池相关配置，最大连接数，建议保持默认
        "maxIdle": 32,               // 连接池相关配置，最大空闲连接数，建议保持默认
        "replicas": 500,             // 这是一致性hash算法需要的节点副本数量，建议不要变更，保持默认即可
        "cluster": {                 // key-value形式的字典，表示后端的judge列表，其中key代表后端judge名字，value代表的是具体的ip:port
            "judge-00" : "127.0.0.1:6080"
        }
    },
    "graph": {
        "enabled": true,                  // 表示是否开启向graph发送数据
        "batch": 200,                     // 数据转发的批量大小，可以加快发送速度，建议保持默认值
        "connTimeout": 1000,              // 位是毫秒，与后端建立连接的超时时间，可以根据网络质量微调，建议保持默认
        "callTimeout": 5000,              // 单位是毫秒，发送数据给后端的超时时间，可以根据网络质量微调，建议保持默认
        "maxConns": 32,                   // 连接池相关配置，最大连接数，建议保持默认
        "maxIdle": 32,                    // 连接池相关配置，最大空闲连接数，建议保持默认
        "replicas": 500,                  // 这是一致性hash算法需要的节点副本数量，建议不要变更，保持默认即可
        "cluster": {                      // key-value形式的字典，表示后端的graph列表，其中key代表后端graph名字，value代表的是具体的ip:port(多个地址用逗号隔开, transfer会将同一份数据发送至各个地址，利用这个特性可以实现数据的多重备份)
            "graph-00" : "127.0.0.1:6070"
        }
    },
    "tsdb": {
        "enabled": false,             // 表示是否开启向open tsdb发送数据
        "batch": 200,                 // 数据转发的批量大小，可以加快发送速度
        "connTimeout": 1000,          // 位是毫秒，与后端建立连接的超时时间，可以根据网络质量微调，建议保持默认
        "callTimeout": 5000,          // 单位是毫秒，发送数据给后端的超时时间，可以根据网络质量微调，建议保持默认
        "maxConns": 32,               // 连接池相关配置，最大连接数，建议保持默认
        "maxIdle": 32,                // 连接池相关配置，最大空闲连接数，建议保持默认
        "retry": 3,                   // 连接后端的重试次数和发送数据的重试次数
        "address": "127.0.0.1:8088"   // tsdb地址或者tsdb集群vip地址, 通过tcp连接tsdb.
    }
}

```



#### 2. 接收命令行参数

这里和agent一样，提供了三个命令行选项，查看版本，指定配置文件和查看git 版本。

```go
cfg := flag.String("c", "cfg.json", "configuration file")
version := flag.Bool("v", false, "show version")
versionGit := flag.Bool("vg", false, "show version")
flag.Parse()

if *version {
	fmt.Println(g.VERSION)
	os.Exit(0)
}
if *versionGit {
	fmt.Println(g.VERSION, g.COMMIT)
	os.Exit(0)
}
```



#### 3. 启动proc

这里很简单，只是告诉用户进行启动了。

```go
func Start() {
	log.Println("proc.Start, ok")
}
```



#### 4. 初始化

这里初始化了 transfer 的 rpc 连接池、数据发送队列、发送任务等。

```go
// 初始化数据发送服务, 在main函数中调用
func Start() {
	// 初始化默认参数
	MinStep = g.Config().MinStep
	if MinStep < 1 {
		MinStep = 30 //默认30s
	}
	//
	initConnPools()
	initSendQueues()
	initNodeRings()
	// SendTasks依赖基础组件的初始化,要最后启动
	startSendTasks()
	startSenderCron()
	log.Println("send.Start, ok")
}
```





#### 5. 开启 rpc 和 socket 服务

```go
func Start() {
	go rpc.StartRpc()
	go socket.StartSocket()
}
```

#### 6. 开启http服务

这里的 http 接口提供了查看系统性能的一系列接口，这里就不介绍了，感兴趣的同学可以看一下源码。



## 二、transfer 中的缓存机制

> 这里我们沿着 transfer 向 graph 发送数据这条线来看 transfer 是如何设置缓存机制的，以及 transfer 是如何发送数据的。

#### 1.  一致性 Hash

在 transfer 中通过一致性 Hash 决定数据发送给那个node，并通过设置虚拟节点的方式达到负载均衡的效果。

#### 2.  连接池

这里的连接池如下：

```go
// ConnPools Manager
type SafeRpcConnPools struct {
	sync.RWMutex
	M           map[string]*connp.ConnPool 
	MaxConns    int
	MaxIdle     int
	ConnTimeout int
	CallTimeout int
}
```

系统根据 cfg.json 文件的配置创建 rpc 连接池，后面会使用连接池进行数据的发送。

#### 3. 发送队列

transfer 共有三种类型的发送队列：

```go
// 发送缓存队列
// node -> queue_of_data
var (
	TsdbQueue   *nlist.SafeListLimited
	JudgeQueues = make(map[string]*nlist.SafeListLimited)
	GraphQueues = make(map[string]*nlist.SafeListLimited)
)
```

底层是通过加锁的双向链表实现的。

以 graph 为例：

```go
DefaultSendQueueMaxSize = 102400 //10.24w 最大发送队列
for node, nitem := range cfg.Graph.ClusterList {
	for _, addr := range nitem.Addrs {
		Q := nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
		GraphQueues[node+addr] = Q
	}
}
```

可以看到这里初始化了 GraphQueues ，每一台 graph 设备都初始化了一个长度为102400的队列。

#### 4. 将数据打入缓存队列

```go
// 当接收到 agnet 发送的数据时，调用RecvMetricValues将数据进行封装，并打入缓存队列
func (t *Transfer) Update(args []*cmodel.MetricValue, reply *cmodel.TransferResponse) error {
	return RecvMetricValues(args, reply, "rpc")
}

// process new metric values
func RecvMetricValues(args []*cmodel.MetricValue, reply *cmodel.TransferResponse, from string) error 

// 将数据 打入 某个Graph的发送缓存队列, 具体是哪一个Graph 由一致性哈希 决定
func Push2GraphSendQueue(items []*cmodel.MetaData)
// 将数据 打入 某个Judge的发送缓存队列, 具体是哪一个Judge 由一致性哈希 决定
func Push2JudgeSendQueue(items []*cmodel.MetaData)
// 将原始数据入到tsdb发送缓存队列
func Push2TsdbSendQueue(items []*cmodel.MetaData)
```



#### 5. 发送数据

系统会给每一台 graph 设备创建一个协程用来发送数据：

```go
for node, nitem := range cfg.Graph.ClusterList {
	for _, addr := range nitem.Addrs {
		queue := GraphQueues[node+addr]
		go forward2GraphTask(queue, node, addr, graphConcurrent)
	}
}
```

我们来看一下协程中都做了什么：

【1】控制发送数据的大小

```go
batch := g.Config().Graph.Batch // 一次发送,最多batch条数据
sema := nsema.NewSemaphore(concurrent)
```

这里通过 nsema.NewSemaphore 来控制发送数据的大小，这底层其实就是一个有缓存的 channel，通过channel 的缓存大小限制并发量。

【2】从队列中取数据

这里其实就是一个从链表拿数据的过程，去问数据后，对应的节点被删除。

```go
items := Q.PopBackBy(batch)
```

【3】通过协程发送数据

```go
for i := 0; i < 3; i++ { //最多重试3次
	err = GraphConnPools.Call(addr, "Graph.Send", graphItems, resp)
	if err == nil {
		sendOk = true
		break
	}
	time.Sleep(time.Millisecond * 10)
}
```

从 GraphConnPools 中获取 rpc 连接，然后发送数据， 每次数据发送有三次机会。

#### 6. 总结

transfer 通过初始化连接池和发送队列，减小了系统发送时频繁创建连接的开销，大大提高系统性能，通过控制缓存队列的大小，避免系统资源被大量消耗。