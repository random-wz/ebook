> 在 open-falcon 监控系统中，judge 模块负责告警判断，它会定期向 hbs 请求告警策略信息，然后将告警信息推送到 redis ，本文向大家介绍 judge 是如何进行告警诊断的，希望对你有帮助。

#### 阅读前需要了解的概念：

##### 1. 监控数据分类

- GAUGE：实测值，直接使用采集的原始数值，比如气温；
- COUNTER：记录连续增长的数据，只增不减。比如汽车行驶里程，网卡流出流量，cpu_idle等；
- DERIVE：变化率，类似COUNTER ，但是可增可减；

GAUGE 类型的值在进行告警诊断的时候直接与告警值比较即可，另外两种则需要根据差值进行判断。



##### 2. 数据模型

open-falcon的"监控项"模型如下:

```go
{
    metric: cpu.busy,               // 监控项名称
    endpoint: open-falcon-host,     // 目标服务器的主机名
    tags: srv=falcon,group=az1,     // tag标签，作用是聚合和归类，在配报警策略时会比较方便。
    value: 10,                      // 监控项数值
    timestamp: `date +%s`,          // 采集时间
    counterType: GAUGE,             // 监控项类型。 只能是COUNTER或者GAUGE二选一，前者表示该数据采集项为计时器类型，后者表示其为原值 (注意大小写)
    step: 60                        // 采集间隔。 秒。
}
```

这种模型有两个好处：

- 一是方便从多个维度来配置告警

  ```
  比如tag的使用起到了给机器进行归类的作用，比如有3台机器：host1、host2和host3，如果tags依次配置为"group=java", "group=java"和"group=erlang"，那么配置报警策略"metric=cpu/group=java“时就只会对java标签的机器（即host1，host2)生效。
  ```

- 二是可以灵活的进行自主数据采集。

  ```
  由于agent会自发现的采集很多基本的系统指标，但是对业务应用等需要研发人员自己写脚本收集和上报。这里openfalcon定义了监控项模型，相当于定义了一个规范，当研发人员需要监控某个对象（比如mysql、redis等），只需采集数据，并遵守规范包装成监控项模型，再上报即可。
  ```



##### 3. 告警模板参数说明

```go
type Strategy struct {
	Id         int               `json:"id"`
	Metric     string            `json:"metric"`     // 告警适用的 metric
	Tags       map[string]string `json:"tags"`       // 告警适用的 tag
	Func       string            `json:"func"`       // e.g. max(#3) all(#3)
	Operator   string            `json:"operator"`   // e.g. < !=
	RightValue float64           `json:"rightValue"` // 告警阈值
	MaxStep    int               `json:"maxStep"`    // 表示该数据采集项的汇报周期，这对于后续的配置监控策略很重要，必须明确指定。
	Priority   int               `json:"priority"`   // 比如P0/P1/P2等等，每个及别的报警都会对应不同的redis队列
	Note       string            `json:"note"`       // 告警信息
	Tpl        *Template         `json:"tpl"`
}
```

需要重点说的是 func 参数，下面是 func 参数可以配置的类型：

- max  根据最大值判断是否触发告警。

  ```
  max(#3): 对于最新的3个点，其最大值满足阈值条件则报警
  ```

- min 根据最小值判断是否触发告警。

  ```
  min(#3): 对于最新的3个点，其最小值满足阈值条件则报警
  ```

- all 

  ```
  all(#3): 最新的3个点都满足阈值条件则报警
  ```

- sum 根据和判断是否触发告警。

  ```
  sum(#3): 对于最新的3个点，其和满足阈值条件则报警
  ```

- avg 根据平均值判断是否触发告警。

  ```
  avg(#3): 对于最新的3个点，其平均值满足阈值条件则报警
  ```

- diff  只要有一个点的diff触发阈值，就报警，diff是当前值减去历史值

  ```
  diff(#3): 拿最新push上来的点（被减数），与历史最新的3个点（3个减数）相减，得到3个差，只要有一个差满足阈值条件则报警
  ```

- pdiff

  ```
  pdiff(#3): 拿最新push上来的点，与历史最新的3个点相减，得到3个差，再将3个差值分别除以减数，得到3个商值，只要有一个商值满足阈值则报警
  ```

- lookup

  ```
  lookup(#2,3): 最新的3个点中有2个满足条件则报警；
  ```

- stddev

  ```
  离群点检测函数，更多请参考3-sigma算法：https://en.wikipedia.org/wiki/68%E2%80%9395%E2%80%9399.7_rule
  	stddev(#10) = 3 //取最新 **10** 个点的数据分别计算得到他们的标准差和均值，分别计为 σ 和 μ，其中当前值计为 X，那么当 X 落在区间 [μ-3σ, μ+3σ] 之外时则报警。
  ```

  

##### 4. falcon_protal 数据库表结构

![falcon_portal](C:\Users\random\Documents\Good Good Study, Day Day Up\博客\open-falcon\falcon_portal.png)



## 一、 judge 模块的启动流程

#### 1. 加载配置文件

```json
{
    "debug": true,
    "debugHost": "nil",
    "remain": 11,
    "http": {
        "enabled": true,
        "listen": "0.0.0.0:6081"
    },
    "rpc": {
        "enabled": true,
        "listen": "0.0.0.0:6080"
    },
    "hbs": {
        "servers": ["127.0.0.1:6030"], # hbs最好放到lvs vip后面，所以此处最好配置为vip:port
        "timeout": 300,
        "interval": 60
    },
    "alarm": {
        "enabled": true,
        "minInterval": 300, # 连续两个报警之间至少相隔的秒数，维持默认即可
        "queuePattern": "event:p%v",
        "redis": {
            "dsn": "127.0.0.1:6379", # 与alarm、sender使用一个redis
            "maxIdle": 5,
            "connTimeout": 5000,
            "readTimeout": 5000,
            "writeTimeout": 5000
        }
    }
}
```



#### 2. 接收命令行参数

这里提供了两个参数，-c 参数用来指定配置文件，-v 参数用来查看 judge 组件版本。

```bash
[root@localhost judge]# ./falcon-judge --help
Usage of ./falcon-judge:
  -c string
        configuration file (default "cfg.json")
  -v    show version
```

同样使用 flag 标准库接收命令行参数：

```go
cfg := flag.String("c", "cfg.json", "configuration file")
version := flag.Bool("v", false, "show version")
flag.Parse()
```



#### 3. 初始化 redis 连接池

这里使用配置文件中指定的 redis 信息，创建 redis 连接池，Judge 模块会将告警信息推送到 redis 中，alarm 模块从 redis 中定期读取告警信息，并进行推送。



#### 4. 初始化 hbs 连接池

通过 Rpc 的方式与 hbs 进行通信，这里主要用来同步告警策略。

```go
func InitHbsClient() {
	HbsClient = &SingleConnRpcClient{
		RpcServers:  Config().Hbs.Servers,
		Timeout:     time.Duration(Config().Hbs.Timeout) * time.Millisecond,
		CallTimeout: time.Duration(3000) * time.Millisecond,
	}
}
```



#### 5. 初始化监控数据池

源码中创建了一个Map:

```go
type JudgeItemMap struct {
	sync.RWMutex
	M map[string]*SafeLinkedList
}

var HistoryBigMap = make(map[string]*JudgeItemMap)

// 初始化HistoryBigMap，为什么这样初始化，需要参考 graph 模块数据存储方式，我会出一片博客讲。
func InitHistoryBigMap() {
	arr := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			HistoryBigMap[arr[i]+arr[j]] = NewJudgeItemMap()
		}
	}
}
// 初始化 HistoryBigMap 的值，本质上其实是一个双向链表
func NewJudgeItemMap() *JudgeItemMap {
	return &JudgeItemMap{M: make(map[string]*SafeLinkedList)}
}
```

在 map 中每个键对应的值是一个双向的安全链表，链表的值保存了 transfer 模块发送的监控数据：

```go
type JudgeItem struct {
	Endpoint  string            `json:"endpoint"`    // 告警的服务器名称，一般是 hostname
	Metric    string            `json:"metric"`      // 监控项
	Value     float64           `json:"value"`       // 采集到的值
	Timestamp int64             `json:"timestamp"`   // 时间戳，根据这个值判断是否清理该节点，有效期为一周
	JudgeType string            `json:"judgeType"`   // 类型
	Tags      map[string]string `json:"tags"`
} 
```



#### 6. 启动 rpc 进程以及 http 服务

**（1）rpc**

judge 模块 rpc 只提供了两个方法：

```go
type Judge int
// 测试连通性
func (this *Judge) Ping(req model.NullRpcRequest, resp *model.SimpleRpcResponse) error {
	return nil
}
// transfer 通过 Send 方法将监控数据发送给 Judge 模块，Judge 将 监控数据保存到 HistoryBigMap 中
func (this *Judge) Send(items []*model.JudgeItem, resp *model.SimpleRpcResponse) error {
	remain := g.Config().Remain
	// 把当前时间的计算放在最外层，是为了减少获取时间时的系统调用开销
	now := time.Now().Unix()
	for _, item := range items {
		exists := g.FilterMap.Exists(item.Metric)
		if !exists {
			continue
		}
		pk := item.PrimaryKey()
		store.HistoryBigMap[pk[0:2]].PushFrontAndMaintain(pk, item, remain, now)
	}
	return nil
}
```

**（2）http**

这里主要定义了一些查看系统参数以及监控信息到的接口：

```go
http://IP:6081/health      查看服务健康状态
http://IP:6081/version     查看模块版本
http://IP:6081/workdir     查看工作目录
http://IP:6081/config/reload 重新加载配置文件
http://IP:6081/strategy/:endpoint/:metric 查看告警策略
http://IP:6081/expression/metric/tag 查看告警表达式
http://IP:6081/count       查看HistoryBigMap大小
http://IP:6081/history/:pk 查看graph节点的数据大小
```



#### 7. 定期更新告警策略

```go
func SyncStrategies() {
	duration := time.Duration(g.Config().Hbs.Interval) * time.Second
	for {
		// 更新告警策略
		syncStrategies()
		syncExpression()
		syncFilter()
		time.Sleep(duration)
	}
}
```



#### 8. 定期清理监控数据

监控数据最多保存一周时间，超时会被清理

```go
func CleanStale() {
	for {
		time.Sleep(time.Hour * 5)
		cleanStale()
	}
}

func cleanStale() {
    // 清理一周前的数据
	before := time.Now().Unix() - 3600*24*7
	arr := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			store.HistoryBigMap[arr[i]+arr[j]].CleanStale(before)
		}
	}
}
```



## 二、judge 模块的告警检测机制

#### 1. 配置告警策略

下图是 open-falcon 自带的 Dashboard 中策略配置的截图：

![image-20200827151356105](C:\Users\random\AppData\Roaming\Typora\typora-user-images\image-20200827151356105.png)

在图中我们可以看到，告警策略中需要配置一下几类监控信息：

- 监控项：metric、tag、最大步长、最多出现次数、等级、告警信息。

- 表达式：具体的告警规则。

- 生效时间：可以配置开始时间和结束时间。

- 告警通知：配置报警接收组以及告警方式（短信、邮件、IM等）。

  

#### 2. 策略表达式

下图为策略表达式的配置：

![image-20200827165308503](C:\Users\random\AppData\Roaming\Typora\typora-user-images\image-20200827165308503.png)

可以看到策略表达式和策略很像，唯一的区别就是，策略表达式可以针对 tag 配置监控策略，而策略模板精确到具体的 Metric 。两者相比，表达式比较简洁，在结合 tag 时可以使用策略表达式。当无法区分类别时，比如所有监控项都没有加tag，只有进行人工分类，即使用”机器分组”，然后将”策略模板”绑定到”机器分组”。



#### 3. 告警诊断

当 transfer 向 judge 模块发送监控数据时，会自动触发告警诊断功能：

```go
// PushFrontAndMaintain 将监控数据保存，并通过 Jusge 方法进行诊断。
func (this *JudgeItemMap) PushFrontAndMaintain(key string, val *model.JudgeItem, maxCount int, now int64) {
	if linkedList, exists := this.Get(key); exists {
		needJudge := linkedList.PushFrontAndMaintain(val, maxCount)
		if needJudge {
			Judge(linkedList, val, now)
		}
	} else {
		NL := list.New()
		NL.PushFront(val)
		safeList := &SafeLinkedList{L: NL}
		this.Set(key, safeList)
		Judge(safeList, val, now)
	}
}
// Judge Judge 会根据告警策略以及告警表达式进行告警判断
func Judge(L *SafeLinkedList, firstItem *model.JudgeItem, now int64) {
	CheckStrategy(L, firstItem, now)
	CheckExpression(L, firstItem, now)
}
```

**（1）告警模板诊断**

```go
func CheckStrategy(L *SafeLinkedList, firstItem *model.JudgeItem, now int64)
```

诊断流程如下：

1. 从缓存中读取告警策略，如果没有相关的告警策略直接退出。

2. 筛选告警策略，根据 tag 筛选出相匹配的告警策略。

3. 诊断，根据采集到的监控数据及监控策略计算出是否需要告警，不满足告警条件则退出。

4. 发送告警信息给 redis。

   

**（2）告警表达式诊断**

```go
func CheckExpression(L *SafeLinkedList, firstItem *model.JudgeItem, now int64)
```

告警表达式诊断过程与根据告警模板诊断的过程类似：

1. 根据监控数据的 metric、endpoint、tag 信息生成 keys 列表。
2. 从缓存中读取告警策略。
3. 筛选告警策略，根据 key 筛选出相匹配的告警策略，如果没有相关的告警策略直接退出。
4. 诊断，根据采集到的监控数据及监控策略计算出是否需要告警，不满足告警条件则退出。
5. 发送告警信息给 redis。

#### 4. 总结

Judge 模块同样通过使用缓存的方式提高诊断速度，通过链表的方式确保数据有序，系统内部提供了丰富的诊断策略，可以覆盖绝大部分监控场景，通过告警模板以及告警表达式覆盖更多监控场景。