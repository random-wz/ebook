> os/signal包实现了对输入信号的访问。这个包只有两个重要方法，这里向大家介绍一下，希望对你有帮助。



#### 1. 信号的转发

> Notify函数让signal包将输入信号转发到c。如果没有列出要传递的信号，会将所有输入信号传递到c；否则只传递列出的输入信号。
>
> signal包不会为了向c发送信息而阻塞（就是说如果发送时c阻塞了，signal包会直接放弃）：调用者应该保证c有足够的缓存空间可以跟上期望的信号频率。对使用单一信号用于通知的通道，缓存为1就足够了。
>
> 可以使用同一通道多次调用Notify：每一次都会扩展该通道接收的信号集。唯一从信号集去除信号的方法是调用Stop。可以使用同一信号和不同通道多次调用Notify：每一个通道都会独立接收到该信号的一个拷贝。

```go
func Notify(c chan<- os.Signal, sig ...os.Signal)
```

Example:

```go
package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	// 初始化一个os.Signal类型的channel
	// 我们必须使用缓冲通道，否则在信号发送时如果还没有准备好接收信号，就有丢失信号的风险。
	c := make(chan os.Signal, 1)
	// notify用于监听信号
	// 参数1表示接收信号的channel
	// 参数2及后面的表示要监听的信号
	// os.Interrupt 表示中断
	// os.Kill 杀死退出进程
	signal.Notify(c, os.Interrupt, os.Kill)
	// 阻塞直到接收到信息
	s := <-c
	fmt.Println("Got signal:", s)
}
```

Output:

```bash
$ go run main.go   // 注意：运行代码后按`Ctrl +C`发送信号结果为:
Got signal: interrupt
```

#### 2. 停止转发信号

> Stop函数让signal包停止向c转发信号。它会取消之前使用c调用的所有Notify的效果。当Stop返回后，会保证c不再接收到任何信号。

```go
func Stop(c chan<- os.Signal)
```



Example:

```go
package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	// 初始化一个os.Signal类型的channel
	// 我们必须使用缓冲通道，否则在信号发送时如果还没有准备好接收信号，就有丢失信号的风险。
	ch := make(chan os.Signal)
	// notify用于监听信号 默认是所有信号
	signal.Notify(ch)

	//停止向ch转发信号，ch将不再收到任何信号
	signal.Stop(ch)
	fmt.Println("signal.Stop")
	//ch将一直阻塞在这里，因为它将收不到任何信号
	//所以下面的exit输出也无法执行
	s := <-ch
	fmt.Println("Got signal:", s)
}
```

Output:

```bash
$ go run main.go
signal.Stop
exit status 2
```

#### 3. 优雅的退出守护进程

我们先来了解一下Golang中的信号类型：

- 在POSIX.1-1990标准中定义的信号列表

| 信号值  | 值       | 动作 | 说明                                                         |
| ------- | -------- | ---- | ------------------------------------------------------------ |
| SIGHUP  | 1        | Term | 终端控制进程结束(终端连接断开)                               |
| SIGINT  | 2        | Term | 用户发送INTR字符(Ctrl+C)触发                                 |
| SIGQUIT | 3        | Core | 用户发送QUIT字符(Ctrl+/)触发                                 |
| SIGILL  | 4        | Core | 非法指令(程序错误、试图执行数据段、栈溢出等)                 |
| SIGABRT | 6        | Core | 调用abort函数触发                                            |
| SIGFPE  | 8        | Core | 算术运行错误(浮点运算错误、除数为零等)                       |
| SIGKILL | 9        | Term | 无条件结束程序(不能被捕获、阻塞或忽略)                       |
| SIGSEGV | 11       | Core | 无效内存引用(试图访问不属于自己的内存空间、对只读内存空间进行写操作) |
| SIGPIPE | 13       | Term | 消息管道损坏(FIFO/Socket通信时，管道未打开而进行写操作)      |
| SIGALRM | 14       | Term | 时钟定时信号                                                 |
| SIGTERM | 15       | Term | 结束程序(可以被捕获、阻塞或忽略)                             |
| SIGUSR1 | 30,10,16 | Term | 用户保留                                                     |
| SIGUSR2 | 31,12,17 | Term | 用户保留                                                     |
| SIGCHLD | 20,17,18 | Ign  | 子进程结束(由父进程接收)                                     |
| SIGCONT | 19,18,25 | Cont | 继续执行已经停止的进程(不能被阻塞)                           |
| SIGSTOP | 17,19,23 | Stop | 停止进程(不能被捕获、阻塞或忽略)                             |
| SIGTSTP | 18,20,24 | Stop | 停止进程(可以被捕获、阻塞或忽略)                             |
| SIGTTIN | 21,21,26 | Stop | 后台程序从终端中读取数据时触发                               |
| SIGTTOU | 22,22,27 | Stop | 后台程序向终端中写数据时触发                                 |

- 在SUSv2和POSIX.1-2001标准中的信号列表:

| 信号      | 值       | 动作 | 说明                                              |
| --------- | -------- | ---- | ------------------------------------------------- |
| SIGTRAP   | 5        | Core | Trap指令触发(如断点，在调试器中使用)              |
| SIGBUS    | 0,7,10   | Core | 非法地址(内存地址对齐错误)                        |
| SIGPOLL   |          | Term | Pollable event (Sys V). Synonym for SIGIO         |
| SIGPROF   | 27,27,29 | Term | 性能时钟信号(包含系统调用时间和进程占用CPU的时间) |
| SIGSYS    | 12,31,12 | Core | 无效的系统调用(SVr4)                              |
| SIGURG    | 16,23,21 | Ign  | 有紧急数据到达Socket(4.2BSD)                      |
| SIGVTALRM | 26,26,28 | Term | 虚拟时钟信号(进程占用CPU的时间)(4.2BSD)           |
| SIGXCPU   | 24,24,30 | Core | 超过CPU时间资源限制(4.2BSD)                       |
| SIGXFSZ   | 25,25,31 | Core | 超过文件大小资源限制(4.2BSD)                      |

<font color=red>注意：需要特别说明的是，SIGKILL和SIGSTOP这两个信号既不能被应用程序捕获，也不能被操作系统阻塞或忽略。</font>

通常我们在Linux系统中会使用kill命令来杀死进程，那其中的原理是什么呢？

1. kill pid 方式

   > kill pid的作用是向进程号为pid的进程发送SIGTERM（这是kill默认发送的信号），该信号是一个结束进程的信号且可以被应用程序捕获。若应用程序没有捕获并响应该信号的逻辑代码，则该信号的默认动作是kill掉进程。这是终止指定进程的推荐做法。

2. kill -9 pid 方式

   > kill -9 pid则是向进程号为pid的进程发送SIGKILL（该信号的编号为9），从本文上面的说明可知，SIGKILL既不能被应用程序捕获，也不能被阻塞或忽略，其动作是立即结束指定进程。通俗地说，应用程序根本无法“感知”SIGKILL信号，它在完全无准备的情况下，就被收到SIGKILL信号的操作系统给干掉了，显然，在这种“暴力”情况下，应用程序完全没有释放当前占用资源的机会。事实上，SIGKILL信号是直接发给init进程的，它收到该信号后，负责终止pid指定的进程。在某些情况下（如进程已经hang死，无响应正常信号），就可以使用kill -9来结束进程。

从上面的介绍不难看出，优雅退出可以通过捕获SIGTERM来实现。具体来讲，通常只需要两步动作：

- 注册SIGTERM信号的处理函数并在处理函数中做一些进程退出的准备。信号处理函数的注册可以通过signal()或sigaction()来实现，其中，推荐使用后者来实现信号响应函数的设置。信号处理函数的逻辑越简单越好，通常的做法是在该函数中设置一个bool型的flag变量以表明进程收到了SIGTERM信号，准备退出。
- 在主进程的main()中，通过类似于while(!bQuit)的逻辑来检测那个flag变量，一旦bQuit在signal handler function中被置为true，则主进程退出while()循环，接下来就是一些释放资源或dump进程当前状态或记录日志的动作，完成这些后，主进程退出。

知道了这些，我们看一下下面的例子：

```go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// 优雅退出go守护进程
func main()  {
	//创建监听退出chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Println("退出", s)
				ExitFunc()
			default:
				fmt.Println("other", s)
			}
		}
	}()

	fmt.Println("进程启动...")
	sum := 0
	for {
		sum++
		fmt.Println("sum:", sum)
		time.Sleep(time.Second)
	}
}

func ExitFunc()  {
	fmt.Println("开始退出...")
	fmt.Println("执行清理...")
	fmt.Println("结束退出...")
	os.Exit(0)
}
```

测试一下：

```bash
$ go run main.go
进程启动...
sum: 1
sum: 2
sum: 3
退出 interrupt // Ctrl + C 
开始退出...
执行清理...
结束退出...
```

