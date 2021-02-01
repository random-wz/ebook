> RPC（Remote Procedure Call）远程过程调用，它可以使一台主机上的进程调用另一台主机的进程，由以访为其他若干个主机提供服务，也就是我们常说的C/S服务，Server与Client之间通过rpc方式进行通信。下面向大叫刨析以下net/rpc标准库，希望对你有帮助。

## 一、Server和Client

#### 1. server

##### （1）Server对象

在Server对象中定义了互斥锁用来保护请求数据，另外还包含请求信息和返回的信息以及注册的服务。

```go
// Server represents an RPC Server.
type Server struct {
	serviceMap sync.Map   // map[string]*service
	reqLock    sync.Mutex // protects freeReq
	freeReq    *Request
	respLock   sync.Mutex // protects freeResp
	freeResp   *Response
}
```

我们可以通过NewServer初始化一个Server对象：

```go
// NewServer returns a new Server.
func NewServer() *Server {
	return &Server{}
}
```

Server对象有8个方法，下面进行介绍。

##### （2）func (server *Server) Register(rcvr interface{}) error

Register用来向Server注册rpc服务，rpc服务必须满足下面五种要求：

- 函数必须是导出的
- 必须有两个导出类型参数
- 第一个参数是接收参数
- 第二个参数是返回给客户端参数，必须是指针类型
- 函数还要有一个返回值error

注册之后的服务，Client可以进行远程调用。

还有一个和Registry类似的方法：

```go
// RegisterName类似Register，但使用提供的name代替rcvr的具体类型名作为服务名。
func (server *Server) RegisterName(name string, rcvr interface{}) error {
	return server.register(rcvr, name, true)
}
```

##### （3）监听器

```go
// Accept接收监听器l获取的连接，然后服务每一个连接。Accept会阻塞，调用者应另开线程："go server.Accept(l)"
func (server *Server) Accept(lis net.Listener)
```

##### （4）服务端处理请求的相关方法

- ServeConn方法

  ```go
  // ServeConn在单个连接上执行server。ServeConn会阻塞，服务该连接直到客户端挂起。
  // 调用者一般应另开线程调用本函数："go server.ServeConn(conn)"。ServeConn在该连接使用gob（参见encoding/gob包）有线格式。
  // 要使用其他的编解码器，可调用ServeCodec方法。
  func (server *Server) ServeConn(conn io.ReadWriteCloser)
  ```

- ServeCodec方法

  ```go
  // ServeCodec类似ServeConn，但使用指定的编解码器，以编码请求主体和解码回复主体。
  func (server *Server) ServeCodec(codec ServerCodec)
  ```

- ServeRequest方法

  ```go
  // ServeRequest类似ServeCodec，但异步的服务单个请求。它不会在调用结束后关闭codec。
  func (server *Server) ServeRequest(codec ServerCodec) error
  ```

- ServeHTTP方法

  ```go
  // ServeHTTP实现了回应RPC请求的http.Handler接口。
  func (server *Server) ServeHTTP(w http.ResponseWriter, req *http.Request)
  ```

- HandleHTTP方法

  ```go
  // HandleHTTP注册server的RPC信息HTTP处理器对应到rpcPath，注册server的debug信息HTTP处理器对应到debugPath。
  // HandleHTTP会注册到http.DefaultServeMux。之后，仍需要调用http.Serve()，一般会另开线程："go http.Serve(l, nil)"
  func (server *Server) HandleHTTP(rpcPath, debugPath string)
  ```

  

#### 2. client

##### （1）Client对象

```go
// Client类型代表RPC客户端。同一个客户端可能有多个未返回的调用，也可能被多个go程同时使用。
type Client struct {
	codec ClientCodec

	reqMutex sync.Mutex // protects following
	request  Request

	mutex    sync.Mutex // protects following
	seq      uint64
	pending  map[uint64]*Call
	closing  bool // user has called Close
	shutdown bool // server has told us to stop
}
```



##### （2）新建一个Client

- 初始化一个Client

  ```go
  // NewClient返回一个新的Client，以管理对连接另一端的服务的请求。它添加缓冲到连接的写入侧，以便将回复的头域和有效负载作为一个单元发送。
  func NewClient(conn io.ReadWriteCloser) *Client
  ```

- 初始化一个Client并指定编码器

  ```go
  // 另外还有一个NewClientWithCodec方法，NewClientWithCodec类似NewClient，但使用指定的编解码器，以编码请求主体和解码回复主体。
  func NewClientWithCodec(codec ClientCodec) *Client
  ```



##### （3）连接服务端

- 通过指定的网络和地址与RPC服务端连接。

  ```go
  func Dial(network, address string) (*Client, error)
  ```

- 通过指定的网络和地址与在默认HTTP RPC路径监听的HTTP RPC服务端连接。

  ```go
  func DialHTTP(network, address string) (*Client, error)
  ```

- 通过在指定的网络、地址和路径与HTTP RPC服务端连接。

  ```go
  func DialHTTPPath(network, address, path string) (*Client, error)
  ```

  

##### （4）调用服务端的服务

- Call方法

  ```go
  // Call调用指定的方法，等待调用返回，将结果写入reply，然后返回执行的错误状态。
  func (client *Client) Call(serviceMethod string, args interface{}, reply interface{}) error
  ```

- 异步调用（Go方法）

  ```go
  // Go异步的调用函数。
  // 本方法Call结构体类型指针的返回值代表该次远程调用。通道类型的参数done会在本次调用完成时发出信号（通过返回本次Go方法的返回值）。
  // 如果done为nil，Go会申请一个新的通道（写入返回值的Done字段）；如果done非nil，done必须有缓冲，否则Go方法会故意崩溃。
  func (client *Client) Go(serviceMethod string, args interface{}, reply interface{}, done chan *Call) *Call
  ```

  

##### （5）关闭Client

```go
func (client *Client) Close() error
```



## 二、go语言标准库中的JSON-RPC

>JSON-RPC，是一个无状态且轻量级的远程过程调用（RPC）传送协议，其传递内容通过 JSON 为主。相较于一般的 REST 通过网址（如 GET /user）调用远程服务器，JSON-RPC 直接在内容中定义了欲调用的函数名称（如 {"method": "getUser"}），这也令开发者不会陷于该使用 PUT 或者 PATCH 的问题之中。
更多JSON-RPC约定参见：https://zh.wikipedia.org/wiki/JSON-RPC

#### 1. 连接Server端

```go
// Dial在指定的网络和地址连接一个JSON-RPC服务端。
func Dial(network, address string) (*rpc.Client, error)
```

#### 2. 创建CLient

- NewClient

  ```go
  // NewClient返回一个新的rpc.Client，以管理对连接另一端的服务的请求。
  func NewClient(conn io.ReadWriteCloser) *rpc.Client
  ```

- NewClientCodec

  ```go
  // NewClientCodec返回一个在连接上使用JSON-RPC的rpc.ClientCodec。
  func NewClientCodec(conn io.ReadWriteCloser) rpc.ClientCodec
  ```

- NewServerCodec

  ```go
  // NewServerCodec返回一个在连接上使用JSON-RPC的rpc. ServerCodec。
  func NewServerCodec(conn io.ReadWriteCloser) rpc.ServerCodec
  ```

#### 3. 处理客户端连接请求

```go
// ServeConn在单个连接上执行DefaultServer。ServeConn会阻塞，服务该连接直到客户端挂起。
// 调用者一般应另开线程调用本函数："go serveConn(conn)"。ServeConn在该连接使用JSON编解码格式。
func ServeConn(conn io.ReadWriteCloser)
```



## 三、rpc通信的三种方式

使用rpc方式通信需要通过下面几步才能完成：

**Server端：**

- 初始化一个Server对象
- 注册服务
- 绑定处理器
- 监听服务

**Client 端：**

- 初始化一个Client对象
- 连接RPC服务端
- 发送请求
- 接收返回值



下面介绍net/rpc的三种连接方式。

#### 1. Http方式

Server端：

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"net/rpc"
)

type Student struct {
	Name   string
	School string
}
type RpcServer struct{}

func (r *RpcServer) Introduce(student Student, words *string) error {
	fmt.Println("student: ", student)
	*words = fmt.Sprintf("Hello everyone, my name is %s, and I am from %s", student.Name, student.School)
	return nil
}

func main() {
	rpcServer := new(RpcServer)
	// 注册rpc服务
	_ = rpc.Register(rpcServer)
	//把服务处理绑定到http协议上
	rpc.HandleHTTP()
	log.Println("http rpc service start success addr:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
```

Client端：

```go
package main

import (
	"fmt"
	"net/rpc"
)

func main() {
	type Student struct {
		Name   string
		School string
	}
    // 连接RPC服务端 Dial会调用NewClient初始化一个Client
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	defer client.Close()
	// 发送请求
	var reply string
	err = client.Call("RpcServer.Introduce", &Student{Name: "random_w", School: "Secret"}, &reply)
	if err != nil {
		panic(err)
	}
	fmt.Println(reply)
}
```

测试：

启动服务端：

```bash
$ go run httpRPC.go
2020/07/15 16:28:30 http rpc service start success addr:8080
student:  {random_w Secret}
```

启动客户端：

```bash
$ go run httpRPCClient.go
Hello everyone, my name is random_w, and I am from Secret
```

从客户端的日志可以看到我们成功执行了服务端的Introduce方法，服务端的日志中也显示了接收到的Student信息。



#### 2. TCP 方式

Server端：

```go
package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

type Student struct {
	Name   string
	School string
}
type RpcServer struct{}

func (r *RpcServer) Introduce(student Student, words *string) error {
	fmt.Println("student: ", student)
	*words = fmt.Sprintf("Hello everyone, my name is %s, and I am from %s", student.Name, student.School)
	return nil
}

func main() {
	rpcServer := new(RpcServer)
	// 注册rpc服务
	_ = rpc.Register(rpcServer)
	// 指定rpc模式为TCP模式，地址为127.0.0.1:8081
	tcpAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8081")
	tcpListen, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("tcp rpc service start success addr:8081")
	for {
		// 监听Client发送的请求
		conn, err3 := tcpListen.Accept()
		if err3 != nil {
			continue
		}
		// 创建一个goroutine处理请求
		go rpc.ServeConn(conn)
	}
}
```

Client端：

```go
package main

import (
	"fmt"
	"net/rpc"
)

func main() {
	type Student struct {
		Name   string
		School string
	}
	// 连接RPC服务端 Dial会调用NewClient初始化一个Client
	client, err := rpc.Dial("tcp", "127.0.0.1:8081")
	if err != nil {
		panic(err)
	}
	defer client.Close()
	// 发送请求
	var reply string
	err = client.Call("RpcServer.Introduce", &Student{Name: "random_w", School: "Secret"}, &reply)
	if err != nil {
		panic(err)
	}
	fmt.Println(reply)
}
```

测试：

启动服务端：

```bash
$ go run tcpRPC.go
2020/07/15 16:13:21 tcp rpc service start success addr:8081
student:  {random_w Secret}
```

运行客户端：

```bash
$ go run tcpRPCClient.go
Hello everyone, my name is random_w, and I am from Secret
```

从客户端的日志可以看到我们成功执行了服务端的Introduce方法，服务端的日志中也显示了接收到的Student信息。

#### 3. jsonrpc方式

jsonrpc方式支持跨语言调用。

Server端：

```go
package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Student struct {
	Name   string
	School string
}

type RpcServer struct{}

func (r *RpcServer) Introduce(student Student, words *string) error {
	fmt.Println("student: ", student)
	*words = fmt.Sprintf("Hello everyone, my name is %s, and I am from %s", student.Name, student.School)
	return nil
}

func main() {
	rpcServer := new(RpcServer)
	// 注册rpc服务
	_ = rpc.Register(rpcServer)
    //jsonrpc是基于TCP协议的，现在他还不支持http协议
	tcpAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8082")
	tcpListen, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	log.Println("tcp json-rpc service start success addr:8082")
	for {
		// 监听客户端请求
		conn, err3 := tcpListen.Accept()
		if err3 != nil {
			continue
		}
		go jsonrpc.ServeConn(conn)
	}
}
```

Client端：

```go
package main

import (
	"fmt"
	"net/rpc/jsonrpc"
)

func main() {
	type Student struct {
		Name   string
		School string
	}
	client, err := jsonrpc.Dial("tcp", "127.0.0.1:8082")
	if err != nil {
		panic(err)
	}
	defer client.Close()
	var reply string
	// 发送json格式的数据
	err = client.Call("RpcServer.Introduce", &Student{
		Name:   "random_w",
		School: "Secret",
	}, &reply)
	if err != nil {
		panic(err)
	}
	fmt.Println(reply)
}
```

测试：

启动服务端：

```bash
$ go run jsonRPC.go
2020/07/15 17:15:42 tcp json-rpc service start success addr:8082
student:  {random_w Secret}
```

运行客户端：

```bash
$ go run jsonRPCClient.go
Hello everyone, my name is random_w, and I am from Secret
```

从客户端的日志可以看到我们成功执行了服务端的Introduce方法，服务端的日志中也显示了接收到的Student信息。

