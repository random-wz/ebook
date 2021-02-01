> sync.pool对象在并发比较高的系统中是非常常见的，这篇博客向大家介绍sync包创建可复用的实例池的原理以及使用方法，希望对你有帮助。

#### 1.  池的内部工作流程

我们调用sync.pool包中的方法的时候，init函数会自动注册一个清理pool对象的方法，该方法会在GC执行前被调用，所以Pool会在调用GC的时候性能较低（初始化的对象都被清理了，重新创建就会产生开销）。Pool只在两次GC之间是有效的，下面是官方的一张图，用来理解池的管理方式：

![sync.Pool workflow in Go 1.12](https://user-gold-cdn.xitu.io/2019/6/16/16b5c627c6d4d686?imageView2/0/w/1280/h/960/format/webp/ignore-error/1)

对于我们创建的每个 `sync.Pool`，go 生成一个连接到每个处理器(处理器即 Go 中调度模型 GMP 的 P，pool 里实际存储形式是 `[P]poolLocal`)的内部池 `poolLocal`。该结构由两个属性组成：

- `private`  能由其所有者访问（push 和 pop 不需要任何锁）
-  `shared`  该属性可由任何其他处理器读取，并且需要并发安全。

实际上，池不是简单的本地缓存，它可以被我们的应用程序中的任何 线程/goroutines 使用。

#### 2. 两个重要方法

- Get方法

  ```go
  // 1）尝试从本地P对应的那个本地池中获取一个对象值, 并从本地池冲删除该值。
  // 2）如果获取失败，那么从共享池中获取, 并从共享队列中删除该值。
  // 3）如果获取失败，那么从其他P的共享池中偷一个过来，并删除共享池中的该值(p.getSlow())。
  // 4）如果仍然失败，那么直接通过New()分配一个返回值，注意这个分配的值不会被放入池中。New()返回用户注册的New函数的值，如果用户未注册New，那么返回nil。
  func (p *Pool) Get() interface{} {
  	if race.Enabled {
  		race.Disable()
  	}
  	l := p.pin()
  	x := l.private
  	l.private = nil
  	runtime_procUnpin()
  	if x == nil {
  		l.Lock()
  		last := len(l.shared) - 1
  		if last >= 0 {
  			x = l.shared[last]
  			l.shared = l.shared[:last]
  		}
  		l.Unlock()
  		if x == nil {
  			x = p.getSlow()
  		}
  	}
  	if race.Enabled {
  		race.Enable()
  		if x != nil {
  			race.Acquire(poolRaceAddr(x))
  		}
  	}
  	if x == nil && p.New != nil {
  		x = p.New()
  	}
  	return x
  }
  ```
  
  
  
- Put方法

  ```go
  // 将x添加到池里面
  // 1）如果放入的值为空，直接return.
  // 2）检查当前goroutine的是否设置对象池私有值，如果没有则将x赋值给其私有成员，并将x设置为nil。
  // 3）如果当前goroutine私有值已经被设置，那么将该值追加到共享列表。
  func (p *Pool) Put(x interface{}) {
  	if x == nil {
  		return
  	}
  	if race.Enabled {
  		if fastrand()%4 == 0 {
  			// Randomly drop x on floor.
  			return
  		}
  		race.ReleaseMerge(poolRaceAddr(x))
  		race.Disable()
  	}
  	l := p.pin()
  	if l.private == nil {
  		l.private = x
  		x = nil
  	}
  	runtime_procUnpin()
  	if x != nil {
  		l.Lock()
  		l.shared = append(l.shared, x)
  		l.Unlock()
  	}
  	if race.Enabled {
  		race.Enable()
	}
  }
  ```
  

#### 3. pool缓存对象的数量和期限

在put方法中我们可以看到，源码中并没有指定pool的大小，所以sync.Pool的缓存对象数量是没有限制的（只受限于内存）。需要注意的是sync.Pool缓存对象的期限是受GC影响的，sync.Pool在init()中向runtime注册了一个cleanup方法，因此sync.Pool缓存的期限只是两次gc之间这段时间。

```go
func init() {
	runtime_registerPoolCleanup(poolCleanup)
}
```



#### 4. 基准测试

下面我们来测试一下使用Pool对程序性能的影响：

```go
package main

import (
	"sync"
	"testing"
)

type Student struct {
	Age int
}

func BenchmarkWithPool(b *testing.B) {
	var s *Student
	// 调用sync.Pool，初始化Student
	var pool = sync.Pool{
		New: func() interface{} { return new(Student) },
	}
	// 基准测试
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10000; j++ {
			// 从内存池中获取Student
			s = pool.Get().(*Student)
			s.Age = j
			// 使用完将s归还到Pool
			pool.Put(s)
		}
	}
}

// 这个例子我们不使用sync.Pool
func BenchmarkWithNoPool(b *testing.B) {
	var s *Student
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10000; j++ {
			s = &Student{Age: i}
			s.Age++
		}
	}
}

```

测试一下：

```bash
$ go test -bench=. -benchmem
goos: windows
goarch: amd64
pkg: test
BenchmarkWithPool-8                10000            182103 ns/op               0 B/op          0 allocs/op
BenchmarkWithNoPool-8              10000            160606 ns/op           80000 B/op      10000 allocs/op
PASS
ok      test    3.771s
```

这里解释一下，测试结果中各字段的含义：

- 第一列为基准测试的函数名称。
- 第二列为基准测试的迭代总次数 b.N。
- 第三列为平均每次迭代所消耗的纳秒数。
- 第四列为平均每次迭代内存所分配的字节数。
- 第五列为平均每次迭代的内存分配次数。

我们可以看到不采用Pool和采用Pool两者在内存分配上有很大的差异，使用Pool大大减少了内存的分配，从而程序性能也就大大提高了。

