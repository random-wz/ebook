> 在编程过程中我们经常会遇到context这个单词，他的中文翻译是上下文。**所谓的上下文就是指语境，每一段程序都有很多的外部变量。只有想Add这种简单的函数才是没有外部变量的。一旦写的一段程序中有了外部变量，这段程序就是不完整的，不能独立运行，要想让他运行，就必须把所有的外部变量的值一个一个的全部传进去，这些值的集合就叫上下文**。本文向大家介绍如何在go语言中使用上下文，希望对你有帮助。

#### 1. context包的四个重要方法

> 在介绍context的四个重要方法前，我们先看一下Context接口：

```go
type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}
```

我们可以看到Context接口定义了四个方法：

- Deadline() (deadline time.Time, ok bool) 方法会返回当前job的deadline和job状态，ok=false表示没有设置deadline。
- Done() <-chan struct{} 当一个job在被取消前完成，该方法返回一个channel。

1. 在WithCancel方法中，调用cancel方法后会触发Done方法。
2. 在WithDeadline方法中，调用deadline方法后会触发Done方法。
3. 在WithTimeout方法中，调用timeout方法后会触发Done方法。

- Err() error 方法当Done已经执行了，就会返回一个error信息说明没有执行的原因，当没有执行就会返回一个nil。
- Value(key interface{}) 方法用来返回Context中保存的键值对。

##### （1）func WithCancel(parent Context) (ctx Context, cancel CancelFunc)

WithCancel方法他会返回一个context对象和一个cancel函数，我们可以通过调用cancel函数来结束goroutine的运行。下面举个例子：

```go
package main

import (
	"context"
	"time"
)

func main() {
	//调用WithCancel方法
	ctx, cancel := context.WithCancel(context.Background())
	//创建一个goroutine
	go work(ctx, "work1")
	// 让子协程有时间运行
	time.Sleep(time.Second * 2)
	//调用cancel会触发ctx.Done()方法，从而让work退出循环
	cancel()
	// 为了避免主协程退出导致子协程未执行完毕就退出，这里做一秒的延时
	time.Sleep(time.Second * 1)
}

// work函数里面是一个无限循环，当select语句检测到ctx.Done()有数据写入后会退出循环
func work(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			println(name, " get message to quit")
			return
		default:
            // 每次循环耗时一秒
			println(name, " is running", time.Now().String())
			time.Sleep(time.Second)
		}
	}
}
```

我们测试一下：

```bash
$ go run main.go
work1  is running 2020-07-15 11:17:52.5416195 +0800 CST m=+0.001958801
work1  is running 2020-07-15 11:17:53.553897 +0800 CST m=+1.014236301
work1  get message to quit
```

我们可以看到通过调用cancel函数结束了work协程。

##### (2) func WithValue(parent Context, key, val interface{}) Context

WithValue可以设置一个key/value的键值对，可以在下游任何一个嵌套的context中通过key获取value。但是不建议使用这种来做goroutine之间的通信。 下面举个例子：

```go
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// 调用WithCancel方法
	ctx1, valueCancel := context.WithCancel(context.Background())
	// 通过WithValue方法给上下文添加数据
	valueCtx := context.WithValue(ctx1, "key", "test value context")
	// 创建协程
	go workWithValue(valueCtx, "value work", "key")
	time.Sleep(time.Second * 2)
	// 结束协程
	valueCancel()
	// 为了避免主协程退出导致子协程未执行完毕就退出，这里做一秒的延时
	time.Sleep(time.Second * 1)
}

// workWithValue函数里面是一个无限循环
// 1) 当select语句检测到ctx.Done()有数据写入后会退出循环
// 2) 退出循环前会输出context中保存的值。
func workWithValue(ctx context.Context, name string, key string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Key:", ctx.Value(key))
			println(name, " get message to quit")
			return
		default:
			println(name, " is running", time.Now().String())
			time.Sleep(time.Second)
		}
	}
}
```

我们测试一下：

```bash
$ go run main.go
value work  is running 2020-07-15 11:23:49.4355245 +0800 CST m=+0.002996401
value work  is running 2020-07-15 11:23:50.4467395 +0800 CST m=+1.014211401
Key: test value context
value work  get message to quit
```

我们可以看到goroutine获取到了上下文中的值。

##### （3）func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)

WithTimeout 函数可以设置一个time.Duration，到了这个时间则会cancel这个context。下面举个例子：

```go
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// 调用WithTimeout方法，并设置deadline为3秒
	ctx2, _ := context.WithTimeout(context.Background(), time.Second*3)
	go work(ctx2, "time cancel")
	// 获取deadline
	deadline, ok := ctx2.Deadline()
	fmt.Println("Deadline:", deadline, "	ok:", ok)
	// 等待子协程运行结束
	time.Sleep(time.Second * 4)
}
// work函数里面是一个无限循环，当select语句检测到ctx.Done()有数据写入后会退出循环
func work(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			println(name, " get message to quit")
			return
		default:
			// 每次循环耗时一秒
			println(name, " is running", time.Now().String())
			time.Sleep(time.Second)
		}
	}
}
```

测试一下：

```bash
random@random-wz MINGW64 /d/GOCODE/Test
$ go run main.go
Deadline: 2020-07-15 11:37:09.2947631 +0800 CST m=+3.002002701  ok: true
time cancel  is running 2020-07-15 11:37:06.2947631 +0800 CST m=+0.002002701
time cancel  is running 2020-07-15 11:37:07.3070628 +0800 CST m=+1.014302401
time cancel  is running 2020-07-15 11:37:08.3085428 +0800 CST m=+2.015782401
time cancel  get message to quit
```

我们可以看到work函数进行了三次循环（每次循环延时1秒）之后就退出了。

##### （4） func WithDeadline(parent Context, d time.Time) (Context, CancelFunc)

WithDeadline WithDeadline函数跟WithTimeout很相近，只是WithDeadline设置的是一个时间点。下面举个例子：

```go
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// 调用WithDeadline方法
	ctx3, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Second*3))
	go work(ctx3, "deadline cancel")
	time.Sleep(time.Second * 4)
}
// work函数里面是一个无限循环，当select语句检测到ctx.Done()有数据写入后会退出循环
func work(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			println(name, " get message to quit")
			return
		default:
			// 每次循环耗时一秒
			println(name, " is running", time.Now().String())
			time.Sleep(time.Second)
		}
	}
}
```

测试一下：

```bash
$ go run main.go
deadline cancel  is running 2020-07-15 11:42:31.9932106 +0800 CST m=+0.001999301
deadline cancel  is running 2020-07-15 11:42:33.0054331 +0800 CST m=+1.014221801
deadline cancel  is running 2020-07-15 11:42:34.0062693 +0800 CST m=+2.015058001
deadline cancel  get message to quit
```

我们可以看到三秒后work协程自动退出。

#### 2. emptyCtx

>在上面的例子中我们可以看到函数context.Background()， 这个函数返回的就是一个emptyCtx，另外context.TODO()和context.Background()是实现同样的效果。
emptyCtx经常被用作在跟节点或者说是最上层的context，因为context是可以嵌套的。在上面的Withvalue的例子中已经看到，先用emptyCtx创建一个context，然后再使用withValue把之前创建的context传入。

```go
// An emptyCtx is never canceled, has no values, and has no deadline. It is not
// struct{}, since vars of this type must have distinct addresses.
type emptyCtx int
```



#### 3. cancelCtx

>cancelCtx是context实现里最重要的一环，context的取消几乎都是使用了这个对象。WithDeadline WithTimeout其实最终都是调用的cancel的cancel函数来实现的。

```go
// A cancelCtx can be canceled. When canceled, it also cancels any children
// that implement canceler.
type cancelCtx struct {
	Context                        // 保存parent Context

	mu       sync.Mutex            // 保护数据
	done     chan struct{}         // 用来标识是否已经被取消。
	children map[canceler]struct{} // 保存所有子canceler
	err      error                 // 已经cancel返回nil,没有cancel返回原因。
}
```



#### 4. .context使用时的注意事项

- 不要把Context存在一个结构体当中，显式地传入函数。Context变量需要作为第一个参数使用，一般命名为ctx。

- 即使方法允许，也不要传入一个nil的Context，如果你不确定你要用什么Context的时候传一个context.TODO。

- 使用context的Value相关方法只应该用于在程序和接口中传递的和请求相关的元数据，不要用它来传递一些可选的参数。

- 同样的Context可以用来传递到不同的goroutine中，Context在多个goroutine中是安全的

#### 5. context的作用

1. 保存上下文数据
2. 控制goroutine的超时



