> 在go语言标准库中，net包提供了可移植的网络I/O接口，包括TCP/IP、UDP、域名解析和Unix域socket。本文向大家介绍net标准库的使用，希望对你有帮助。

## 一、 服务端

#### 1.  解析地址

在TCP服务端我们需要监听一个TCP地址，因此建立服务端前我们需要生成一个正确的TCP地址，这就需要用到下面的函数了。

```go
// ResolveTCPAddr函数会输出一个TCP连接地址和一个错误信息
func ResolveTCPAddr(network, address string) (*TCPAddr, error)
// 解析IP地址
func ResolveIPAddr(net, addr string) (*IPAddr, error)
// 解析UDP地址
func ResolveUDPAddr(net, addr string) (*UDPAddr, error)
// 解析Unix地址
func ResolveUnixAddr(net, addr string) (*UnixAddr, error)
```



#### 2. 监听请求

我们可以通过 Listen方法监听我们解析后的网络地址。

```go
// 监听net类型，地址为laddr的地址
func Listen(net, laddr string) (Listener, error)
// 监听TCP地址
func ListenTCP(network string, laddr *TCPAddr) (*TCPListener, error) 
// 监听IP地址
func ListenIP(netProto string, laddr *IPAddr) (*IPConn, error)
// 监听UDP地址
func ListenMulticastUDP(net string, ifi *Interface, gaddr *UDPAddr) (*UDPConn, error)
func ListenUDP(net string, laddr *UDPAddr) (*UDPConn, error)
// 监听Unix地址
func ListenUnixgram(net string, laddr *UnixAddr) (*UnixConn, error)
func ListenUnix(net string, laddr *UnixAddr) (*UnixListener, error)
```



#### 3. 接收请求

TCPAddr 实现了两个接受请求的方法，两者代码实现其实是一样的，唯一的区别是第一种返回了一个对象，第二种返回了一个接口。

```go
func (l *TCPListener) AcceptTCP() (*TCPConn, error)
func (l *TCPListener) Accept() (Conn, error) 
```

其他类型也有类似的方法，具体请参考go语言标准库文档。

#### 4. 连接配置

- 配置监听器超时时间

  ```go
  // 超过t之后监听器自动关闭，0表示不设置超时时间
  func (l *TCPListener) SetDeadline(t time.Time) error
  ```

- 关闭监听器

  ```go
  // 关闭监听器
  func (l *TCPListener) Close() error
  ```

  

#### 5. 编写一个服务端

```go
func main() {
	// 解析服务端监听地址，本例以tcp为例
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Panic(err)
	}
	// 创建监听器
	listen, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Panic(err)
	}
	for {
		// 监听客户端连接请求
		conn, err := listen.AcceptTCP()
		if err != nil {
			continue
		}
		// 处理客户端请求 这个函数可以自己编写
		go HandleConnectionForServer(conn)
	}
}
```



## 二、 TCP客户端

#### 1.  解析TCP地址

在TCP服务端我们需要监听一个TCP地址，因此建立服务端前我们需要生成一个正确的TCP地址，这就需要用到下面的函数了。

```go
// ResolveTCPAddr函数会输出一个TCP连接地址和一个错误信息
func ResolveTCPAddr(network, address string) (*TCPAddr, error)
```

#### 2.  发送连接请求

net包提供了多种连接方法

```go
// DialIP的作用类似于IP网络的拨号
func DialIP(network string, laddr, raddr *IPAddr) (*IPConn, error)
// Dial 连接到指定网络上的地址，涵盖
func Dial(network, address string) (Conn, error)
// 这个方法只是在Dial上面设置了超时时间
func DialTimeout(network, address string, timeout time.Duration) (Conn, error)
// DialTCP 专门用来进行TCP通信的
func DialTCP(network string, laddr, raddr *TCPAddr) (*TCPConn, error)
// DialUDP 专门用来进行UDP通信的
func DialUDP(network string, laddr, raddr *UDPAddr) (*UDPConn, error)
// DialUnix 专门用来进行 Unix 通信
func DialUnix(network string, laddr, raddr *UnixAddr) (*UnixConn, error)
```



#### 3. 编写一个客户端

通过下面的例子我们看一下如何编写一个 TCP 客户端：

```go
func main() {
	// 解析服务端地址
	RemoteAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8000")
	if err != nil {
		panic(err)
	}
	// 解析本地连接地址
	LocalAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1")
	if err != nil {
		panic(err)
	}
	// 连接服务端
	conn, err := net.DialTCP("tcp", LocalAddr, RemoteAddr)
	if err != nil {
		panic(err)
	}
	// 连接管理
	HandleConnectionForClient(conn)
}
```



## 三、 管理连接

> 这里我们来实现一个智能机器人的功能。



#### 1. 客户端

我们通过 HandleConnectionForClient(conn) 方法来处理客户端的消息，话不多说，看代码：

```go
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var sig = make(chan os.Signal)

func main() {
	// 解析服务端地址
	RemoteAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8000")
	if err != nil {
		panic(err)
	}
	// 解析本地连接地址
	LocalAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8001")
	if err != nil {
		panic(err)
	}
	// 连接服务端
	conn, err := net.DialTCP("tcp", LocalAddr, RemoteAddr)
	if err != nil {
		panic(err)
	}
	// 连接管理
	HandleConnectionForClient(conn)
}

// handleConnection 读取数据, 在这里我们可以编写自己的交互程序
func HandleConnectionForClient(conn net.Conn) {
	// 监控系统信号
	go signalMonitor(conn)
	// 初始化一个缓存区
	Stdin := bufio.NewReader(os.Stdin)
	for {
		// 接收服务端返回的消息
		getResponse(conn)
		// 读取用户输入的信息，遇到换行符结束。
		fmt.Print("[ random_w ]# ")
		input, err := Stdin.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		// 删除字符串前后的空格，主要是删除换行符。
		input = strings.TrimSpace(input)
		// 空行不做处理
		if len(input) == 0 {
			continue
		}
		// 是否接收到退出指令
		switch input {
		case "quit", "exit":
			sig <- syscall.SIGQUIT
		default:
			// 发送消息给服务端
			sendMsgToServer(conn, input)
		}
	}
}

// sendMsgToServer 发送消息给服务端
func sendMsgToServer(conn net.Conn, msg string) {
	for {
		_, err := conn.Write([]byte(msg))
		if err == nil {
			break
		}
	}
}

// getResponse 接收服务端返回的消息
func getResponse(conn net.Conn) {
	// 初始化一个1024字节的内存，用来接收服务端的消息
	respByte := make([]byte, 1024)
	// 接收服务端返回的消息
	length, err := conn.Read(respByte)
	if err != nil {
		fmt.Println("[ server ]# 接收消息失败")
	}
	for line, str := range strings.Split(string(respByte[:length]), "\n") {
		if len(str) != 0 {
			if line == 1 {
				fmt.Print(fmt.Sprintf("[ server ]# \n%s\n", str))
				continue
			}
			fmt.Println(str)
		}
	}
}

// signalMonitor 监听系统信号，如果程序收到退出到的信号通过 Goroutine 通知 server 端，关闭连接后退出。
func signalMonitor(conn net.Conn) {
	signal.Notify(sig, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGINT)
	// 接收到结束信号退出此程序
	select {
	case <-sig:
		// 通知服务端断开连接
		_, _ = conn.Write([]byte("exit"))
		fmt.Println("\nGood Bye !!!!!")
		os.Exit(0)
	}
}
```



#### 2. 服务端

我们通过 HandleConnectionForServer(conn) 方法来处理服务端的连接信息。

```go
package main

import (
	"log"
	"net"
)

func main() {
	// 解析服务端监听地址
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Panic(err)
	}
	// 创建监听器
	listen, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Panic(err)
	}
	for {
		// 监听客户端连接请求
		conn, err := listen.AcceptTCP()
		if err != nil {
			continue
		}
		// 处理客户端请求
		go handleConnectionForServer(conn)
	}
}

// handleConnection 读取数据, 在这里我们可以编写自己的交互程序
func handleConnectionForServer(conn net.Conn) {
	for flag := false; ; {
		// 设置消息长度为1024比特
		buf := make([]byte, 1024)
		if !flag {
			// 客户端连接成功，提示可以操作的内容
			if _, err := conn.Write([]byte(Usage())); err != nil {
				log.Println("Error: ", err)
			}
			flag = true
			continue
		}
		/* 读取客户端发送的数据，数据会保存到buf
		这里有一个知识点:
		conn.Read会返回接收到的值的长度，如果不指定长度，通过string转换的时候你会活得一个1024字节的字符串
		但我们不需要后面的初始化的值，因此通过buf[:length]提取我们想要的值。
		*/
		if length, err := conn.Read(buf); err != nil {
			// 读取失败
			writeResponse(parseRequest(""), conn)
		} else {
			// 读取成功
			req := string(buf[:length])
			if req == "exit" {
				break
			}
			writeResponse(parseRequest(req), conn)
		}
	}
}
func Usage() string {
	return `
---------------------------------------------------------------
Hello, my name is randow_w, I'm glad to serve you.
I can provide you with the following services:
1.查工资
2.猜年龄
3.查天气
----------------------------------------------------------------`
}

// writeResponse 返回信息给客户端
func writeResponse(resp string, conn net.Conn) {
	if _, err := conn.Write([]byte(resp)); err != nil {
		log.Println("Error: ", err)
	}
}

// parseRequest 解析客户端输入的信息
func parseRequest(req string) (resp string) {
	switch req {
	case "查工资":
		resp = checkSalary()
	case "猜年龄":
		resp = guessAge()
	case "查天气":
		resp = chat()
	default:
		resp = "对不起，我爸爸还没有教我怎么回答你，能不能换一个问题(*^_^*)"
	}
	return
}

// 查工资
func checkSalary() string {
	return "据权威机构推测，你未来有机会冲刺福布斯排行榜，加油哦(ง •_•)ง"
}

// 猜年龄
func guessAge() string {
	return "永远18岁"
}

// 聊天
func chat() string {
	return "你好，主人，今天是晴天，空气质量优，适合去爬山。"
}
```

<font color=red>注意：服务端里面你自己也可以定义一些方法用来处理客户端的请求，这里只写了几个简单的例子。</font>

#### 3. 测试

启动服务端：

```bash
$ go run server.go
```

启动客户端：

```bash
$ go run client.go
[ server ]#
---------------------------------------------------------------
Hello, my name is randow_w, I'm glad to serve you.
I can provide you with the following services:
1.查工资
2.猜年龄
3.查天气
----------------------------------------------------------------
[ random_w ]# 查工资
据权威机构推测，你未来有机会冲刺福布斯排行榜，加油哦(ง •_•)ง
[ random_w ]# 猜年龄
永远18岁
[ random_w ]# 查天气
你好，主人，今天是晴天，空气质量优，适合去爬山。
[ random_w ]# 你好
对不起，我还在爸爸没有教我怎么回答你，能不能换一个问题(*^_^*)
[ random_w ]# quit

Good Bye !!!!!
```



## 四、 UDP

> 通过net包我们还可以创建一个UDP连接，下面我们通过代码学习如何创建UDP通信的客户端和服务端。

#### 1. UDP 服务端

```go
package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	// 解析服务端监听地址
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8000")
	if err != nil {
		log.Panic(err)
	}
	// 创建监听器
	listen, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Panic(err)
	}
	for {
		// 设置消息长度为1024比特
		buf := make([]byte, 1024)
		// 读取消息，UDP不是面向连接的因此不需要等待连接
		length, udpAddr, err := listen.ReadFromUDP(buf)
		if err != nil {
			log.Println("Error: ", err)
			continue
		}
		fmt.Println("[ server ]# UdpAddr: ", udpAddr, "Data: ", string(buf[:length]))
	}
}

```



#### 2. UDP 客户端

```go
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var sig = make(chan os.Signal)

func main() {
	// 解析服务端地址
	RemoteAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8000")
	if err != nil {
		panic(err)
	}
	// 解析本地连接地址
	LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8001")
	if err != nil {
		panic(err)
	}
	// 连接服务端
	conn, err := net.DialUDP("udp", LocalAddr, RemoteAddr)
	if err != nil {
		panic(err)
	}
	// 连接管理
	HandleConnectionForClient(conn)
}

// handleConnection 读取数据, 在这里我们可以编写自己的交互程序
func HandleConnectionForClient(conn net.Conn) {
	// 监控系统信号
	go signalMonitor(conn)
	// 初始化一个缓存区
	Stdin := bufio.NewReader(os.Stdin)
	for {
		// 读取用户输入的信息，遇到换行符结束。
		fmt.Print("[ random_w ]# ")
		input, err := Stdin.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		// 删除字符串前后的空格，主要是删除换行符。
		input = strings.TrimSpace(input)
		// 空行不做处理
		if len(input) == 0 {
			continue
		}
		// 是否接收到退出指令
		switch input {
		case "quit", "exit":
			sig <- syscall.SIGQUIT
		default:
			// 发送消息给服务端
			sendMsgToServer(conn, input)
		}
	}
}

// sendMsgToServer 发送消息给服务端
func sendMsgToServer(conn net.Conn, msg string) {
	for {
		_, err := conn.Write([]byte(msg))
		if err == nil {
			break
		}
	}
}

// signalMonitor 监听系统信号，如果程序收到退出到的信号通过 Goroutine 通知 server 端，关闭连接后退出。
func signalMonitor(conn net.Conn) {
	signal.Notify(sig, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGINT)
	// 接收到结束信号退出此程序
	select {
	case <-sig:
		// 通知服务端断开连接
		_, _ = conn.Write([]byte("exit"))
		fmt.Println("\nGood Bye !!!!!")
		os.Exit(0)
	}
}

```

#### 3. 测试

开启服务端：

```bash
$ go run udpserver.go
```

开启客户端并传递信息：

```bash
$ go run udpclient.go
[ random_w ]# hello world
[ random_w ]# udp test
[ random_w ]# exit
[ random_w ]#
Good Bye !!!!!
```

服务端接收到消息：

```bash
$ go run udpserver.go
[ server ]# UdpAddr:  127.0.0.1:8001 Data:  hello world
[ server ]# UdpAddr:  127.0.0.1:8001 Data:  udp test
[ server ]# UdpAddr:  127.0.0.1:8001 Data:  exit
```



## 五、 域名解析

#### 1. dns 正向解析

> CNAME 被称为规范名字。这种记录允许您将多个名字映射到同一台计算机。 通常用于同时提供WWW和MAIL服务的计算机。例如，有一台计算机名为“r0WSPFSx58.”（A记录）。 它同时提供WWW和MAIL服务，为了便于用户访问服务。可以为该计算机设置两个别名（CNAME）：WWW和MAIL。

- 域名解析到cname

  ```go
  func LookupCNAME(name string) (cname string, err error)
  ```

- 域名解析到地址

  ```go
  func LookupHost(host string) (addrs []string, err error)
  ```

- 域名解析到地址[]IP结构体.可以对具体ip进行相关操作(是否回环地址,子网,网络号等)

  ```go
  func LookupIP(host string) (addrs []IP, err error)
  ```

#### 2. dns 反向解析

```go
// 根据ip地址查找主机名地址(必须得是可以解析到的域名)[dig -x ipaddress]
func LookupAddr(addr string) (name []string, err error)
```

#### 3. 应用

```go
package main

import (
	"fmt"
	"net"
)

func main() {
    // 域名改成自己要测试的
	dns := "www.baidu.com"
	// 解析cname
	cname, _ := net.LookupCNAME(dns)
	fmt.Println("cname:", cname)
	// 解析ip地址
	ips, err := net.LookupHost(dns)
	if err != nil {
		fmt.Println("Err: ", err.Error())
		return
	}
	fmt.Println(ips)
	// 反向解析(主机必须得能解析到地址), IP地址改成你的
	dnsName, _ := net.LookupAddr("10.X.X.X")
	fmt.Println("Hostname:", dnsName)
}
```

Output:

```go
$ go run main.go
cname: www.a.shifen.com.
[14.215.177.38 14.215.177.39]
Hostname: [paas.bk.com. cmdb.bk.com. job.bk.com.]
```