> 在分布式系统开发中，我们经常会遇到服务器负载均衡的问题，我们需要将用户的请求均匀的分摊到每一台服务器，从而保证系统资源的有效利用。那么如何将请求均匀的进行分配呢？比较常见的就是 Hash 算法了，但是普通的余数hash（hash(比如用户id)%服务器机器数）算法伸缩性很差，当新增或者下线服务器机器时候，用户id与服务器的映射关系会大量失效。一致性hash则利用hash环对其进行了改进。

#### 1. Hash 算法

Hash，一般翻译做“散列”，也有直接音译为“哈希”的，就是把任意长度的输入（又叫做预映射， pre-image），通过散列算法，变换成固定长度的输出，该输出就是散列值。这种转换是一种压缩映射，也就是，散列值的空间通常远小于输入的空间，不同的输入可能会散列成相同的输出，而不可能从散列值来唯一的确定输入值。简单的说就是一种将任意长度的消息压缩到某一固定长度的消息摘要的函数。

常用的 Hash 算法如下：

（1）MD4

> MD4（RFC 1320）是 MIT 的Ronald L. Rivest在 1990 年设计的，MD 是 Message Digest（消息摘要） 的缩写。它适用在32位字长的处理器上用高速软件实现——它是基于 32位操作数的位操作来实现的。

（2）MD5

> MD5（RFC 1321）是 Rivest 于1991年对MD4的改进版本。它对输入仍以512位分组，其输出是4个32位字的级联，与 MD4 相同。MD5比MD4来得复杂，并且速度较之要慢一点，但更安全，在抗分析和抗差分方面表现更好。

（3）SHA-1及其他

> SHA1是由NIST NSA设计为同DSA一起使用的，它对长度小于264的输入，产生长度为160bit的散列值，因此抗穷举（brute-force）性更好。SHA-1 设计时基于和MD4相同原理，并且模仿了该算法。

#### 2. 一致性 Hash 算法

> 一致性哈希算法在1997年由[麻省理工学院](https://baike.baidu.com/item/麻省理工学院/117999)提出，是一种特殊的哈希算法，目的是解决分布式缓存的问题。在移除或者添加一个服务器时，能够尽可能小地改变已存在的服务请求与处理请求服务器之间的映射关系。一致性哈希解决了简单哈希算法在分布式[哈希表](https://baike.baidu.com/item/哈希表/5981869)( Distributed Hash Table，DHT) 中存在的动态伸缩等问题  。

我们通过算法来实现一个一致性 Hash:

```mermaid
graph LR
A(初始化哈希环) --> B(对哈希环的key进行排序)
B --> C(通过hash值匹配node节点)
```

```go
package main

import (
	"fmt"
	"hash/crc32"
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"time"
)

// 一致性 Hash 节点副本数量
const DEFAULT_REPLICAS = 100

// SortKeys 存储一致性 Hash 的值
type SortKeys []uint32

// Len 一致性哈希数量
func (sk SortKeys) Len() int {
	return len(sk)
}

// Less Hash 值的比较
func (sk SortKeys) Less(i, j int) bool {
	return sk[i] < sk[j]
}

// Swap 交换两个 Hash 值
func (sk SortKeys) Swap(i, j int) {
	sk[i], sk[j] = sk[j], sk[i]
}

// Hash 环，存储每一个节点的信息
type HashRing struct {
	Nodes map[uint32]string
	Keys  SortKeys
	sync.RWMutex
}

// New 根据node新建一个Hash环
func (hr *HashRing) New(nodes []string) {
	if nodes == nil {
		return
	}
	hr.Nodes = make(map[uint32]string)
	hr.Keys = SortKeys{}
	for _, node := range nodes {
		// Hash 通过 node 节点名称生成哈希值，该Hash值指向对应 node 节点
		hr.Nodes[hr.Hash(str)] = node
		// 将哈希值保存在key列表中
		hr.Keys = append(hr.Keys, hr.Hash(str))
	}
	// 对 Hash 值进行排序，后面计算的 Hash 值与 Keys 进行比较，取大于等于计算所得的Hash值所对应的node节点
	sort.Sort(hr.Keys)
}

// hashStr 根据Key计算Hash值
func (hr *HashRing) Hash(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

// GetNode 根据key找出对应的node节点
func (hr *HashRing) GetNode(key string) string {
	hr.RLock()
	defer hr.RUnlock()
	hash := hr.Hash(key)
	i := hr.get_position(hash)
	return hr.Nodes[hr.Keys[i]]
}

// get_position 找出第一个大于等于 hash 的key
func (hr *HashRing) get_position(hash uint32) (i int) {
	i = sort.Search(len(hr.Keys), func(i int) bool {
		return hr.Keys[i] >= hash
	})
	if i >= len(hr.Keys) {
		return 0
	}
	return
}

func main() {
	var nodes []string
	// 初始化 node 节点
	for i := 1; i < 6; i++ {
		nodes = append(nodes, fmt.Sprintf("Server%d", i))
	}
	hashR := new(HashRing)
	// 生成 hash 环
	hashR.New(nodes)
	rand.Seed(time.Now().Unix())
	// 寻找匹配的 node 节点
	fmt.Println("random1", " 发送到 ", hashR.GetNode("random1"))
}
```



#### 3. 一致性 Hash 算法的改进

在服务器数量比较少的情况下，一致性 Hash 非常容易出现数据倾斜的问题，为了解决这个问题，人们引入了一致性 Hash 来解决这个问题。

```go
// New 根据node新建一个Hash环
func (hr *HashRing) New(nodes []string) {
	if nodes == nil {
		return
	}
	hr.Nodes = make(map[uint32]string)
	hr.Keys = SortKeys{}
	for _, node := range nodes {
		// 每个节点生成 DEFAULT_REPLICAS 个虚拟节点
		for i := 0; i < DEFAULT_REPLICAS; i++ {
			str := node + strconv.Itoa(i)
			// Hash 通过 node 虚拟节点名称生成哈希值，该Hash值指向对应 node 节点
			hr.Nodes[hr.Hash(str)] = node
			// 将哈希值保存在key列表中
			hr.Keys = append(hr.Keys, hr.Hash(str))
		}
	}
	// 对 Hash 值进行排序，后面计算的 Hash 值与 Keys 进行比较，取大于等于计算所得的Hash值所对应的node节点
	sort.Sort(hr.Keys)
}
// Hash 根据Key计算Hash值
func (hr *HashRing) Hash(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}
```