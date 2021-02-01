###  一、gorilla/websocket

#### 1. 安装

```bash
go get github.com/gorilla/websocket
```

#### 2. upgrader 

upgrader 用于升级 http 请求，把 http 请求升级为长连接的 WebSocket。结构如下：

```go
type Upgrader struct {
    // 指定升级 websocket 握手完成的超时时间
    HandshakeTimeout time.Duration

    // 指定 io 操作的缓存大小，如果不指定就会自动分配。
    ReadBufferSize, WriteBufferSize int

    // 写数据操作的缓存池，如果没有设置值，write buffers 将会分配到链接生命周期里。
    WriteBufferPool BufferPool

    //按顺序指定服务支持的协议，如值存在，则服务会从第一个开始匹配客户端的协议。
    Subprotocols []string

    // 指定 http 的错误响应函数，如果没有设置 Error 则，会生成 http.Error 的错误响应。
    Error func(w http.ResponseWriter, r *http.Request, status int, reason error)

    // 请求检查函数，用于统一的链接检查，以防止跨站点请求伪造。如果不检查，就设置一个返回值为true的函数。
    // 如果请求Origin标头可以接受，CheckOrigin将返回true。 如果CheckOrigin为nil，则使用安全默认值：如果Origin请求头存在且原始主机不等于请求主机头，则返回false
    CheckOrigin func(r *http.Request) bool

    // EnableCompression 指定服务器是否应尝试协商每个邮件压缩（RFC 7692）。 
    // 将此值设置为true并不能保证将支持压缩。 
    // 目前仅支持“无上下文接管”模式
    EnableCompression bool
}
```

Upgrade 函数将 http 升级到 WebSocket 协议。函数签名如下：

```go
// responseHeader包含在对客户端升级请求的响应中。 
// 使用responseHeader指定cookie（Set-Cookie）和应用程序协商的子协议（Sec-WebSocket-Protocol）。
// 如果升级失败，则升级将使用HTTP错误响应回复客户端
// 返回一个 Conn 指针，拿到他后，可使用 Conn 读写数据与客户端通信。
func (u *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*Conn, error)
```

#### 3. 数据类型

WebSocket协议区分文本消息和二进制数据消息：

- 文本消息被解释为UTF-8编码的文本。
- 二进制消息的解释留给应用程序。

#### 4. 控制消息

WebSocket协议定义了三种**控制消息**类型：close、ping和pong。调用WriteControl，WriteMessage或NextWriter方法将控制消息发送到对等方。

- **Close 消息：**连接通过调用使用SetCloseHandler方法设置的处理函数并通过从NextReader、ReadMessage或消息Read方法返回`*CloseError`来处理收到的关闭消息。 默认的关闭处理程序将关闭消息发送到对等方。
- **Ping 消息：**连接通过调用使用SetPingHandler方法设置的处理函数来处理收到的ping消息。默认的ping处理程序将pong消息发送到对等方。
- **Pong 消息：**连接通过调用使用SetPongHandler方法设置的处理函数来处理收到的Pong消息。 默认的pong处理程序不执行任何操作。 如果应用程序发送ping消息，则该应用程序应设置一个pong处理程序以接收相应的pong。

从NextReader，ReadMessage和消息阅读器的Read方法调用控制消息处理程序函数。当处理程序写入连接时，默认的close和ping处理程序可以在短时间内阻止这些方法。

应用程序必须读取连接以处理从对等方发送的关闭，ping和pong消息。 如果应用程序对来自对等方的消息不感兴趣，则应用程序应启动goroutine来读取和丢弃来自对等方的消息。 一个简单的例子是：

```go
func readLoop(c *websocket.Conn) {
    for {
        if _, _, err := c.NextReader(); err != nil {
            c.Close()
            break
        }
    }
}
```

#### 5. 并发消息

连接支持一个并发读取器和一个并发写入器。

应用程序负责确保不超过一个goroutine同时调用写方法`（NextWriter，SetWriteDeadline，WriteMessage，WriteMessage，WriteJSON，EnableWriteCompression，SetCompressionLevel）`，并确保不超过一个goroutine调用读取方法`（NextReader，SetReadDeadline，ReadMessage，ReadJSON，SetPongHandler，SetPingHandler）`并发。 可以与所有其他方法同时调用 `Close` 和 `WriteControl` 方法。

#### 6. 跨域问题

Web浏览器允许Javascript应用程序打开与任何主机的WebSocket连接。由服务器决定是否使用浏览器发送的Origin请求标头来实施Origin策略。

Upgrader 调用CheckOrigin字段中指定的函数以检查源：

- 如果CheckOrigin函数返回false，则Upgrade方法使HTTP状态为403的WebSocket握手失败。
- 如果CheckOrigin字段为nil，则升级程序使用安全默认值：如果存在Origin请求头且Origin主机不等于Host request头，则握手失败。

不建议使用包级别 Upgrade 功能不执行跨域检查。该应用程序负责在调用 Upgrade 功能之前检查Origin标头。

#### 7. 缓冲区

连接可缓冲网络输入和输出，以减少读取或写入消息时的系统调用次数。

写缓冲区还用于构造WebSocket框架。 有关消息帧的讨论，请参见[RFC 6455, Section 5](http://tools.ietf.org/html/rfc6455#section-5) 。 每次将写缓冲区刷新到网络时，都会将WebSocket帧标头写入网络。 减小写缓冲区的大小可能会增加连接上的成帧开销。

缓冲区大小（以字节为单位）由 `Dialer` 和 `Upgrader` 中的 `ReadBufferSize` 和 `WriteBufferSize` 字段指定。 当缓冲区大小字段设置为零时，拨号程序使用默认大小4096。 当缓冲区大小字段设置为零时，升级程序将重用HTTP服务器创建的缓冲区。 在撰写本文时，HTTP服务器缓冲区的大小为4096。

缓冲区大小不限制连接可以读取或写入的消息的大小。

默认情况下，缓冲区在连接的整个生命周期内保持不变。如果设置了 `Dialer` 或 `Upgrader WriteBufferPool` 字段，则连接仅在写消息时才保留写缓冲区。

应用程序应调整缓冲区大小以平衡内存使用和性能。 增加缓冲区大小会占用更多内存，但可以减少读取或写入网络的系统调用次数。 在写入的情况下，增加缓冲区大小可以减少写入网络的帧头的数量。设置缓冲区参数的一些准则是：

- 将缓冲区大小限制为最大预期消息大小。大于最大消息的缓冲区不会带来任何好处。
- 根据消息大小的分布，将缓冲区大小设置为小于最大期望消息大小的值可以大大减少内存使用，而对性能的影响很小。 这是一个示例：如果99％的消息小于256字节，并且最大消息大小为512字节，则256字节的缓冲区大小将导致比512字节的缓冲区大小多1.01个系统调用。 内存节省为50％。
- 当应用程序通过大量连接进行适度的写操作时，写缓冲池很有用。 当缓冲池合并时，较大的缓冲区大小将减少对总内存使用的影响，并具有减少系统调用和帧开销的好处。

#### 8. 消息压缩

每个消息压缩扩展 ([RFC 7692](http://tools.ietf.org/html/rfc7692))在有限的能力下都受到该程序包的实验性支持。在 `Dialer` 或 `Upgrader` 中将 `EnableCompression` 选项设置为`true` 将尝试协商每个消息支持。

```go
var upgrader = websocket.Upgrader{
    EnableCompression: true,
}
```

如果与连接的对等方成功协商了压缩，则以压缩形式收到的任何消息都将自动解压缩。所有 `Read` 方法将返回未压缩的字节。通过调用相应的 `Conn` 方法，可以启用或禁用写入连接的消息的每个消息压缩：

```go
conn.EnableWriteCompression(false)
```

当前改包不支持通过上下文控制压缩，这意味着必须隔离地压缩和解压缩消息，而不能在消息之间保留滑动窗口或字典状态。更多信息可以参考[RFC 7692](http://tools.ietf.org/html/rfc7692).

需要注意的是压缩的使用是实验性的，可能会导致性能下降。

#### 9. 举个例子

> 我们实现一个监听文件内容的功能

 **（1）建立连接**

```go
// 定义 ping、pong 和 监听时间间隔
const (
	writeWait = 1 * time.Second
	pingWait = 50 * time.Second
	pongWait = 60 * time.Second
)

func WebSocketServer(w http.ResponseWriter, r *http.Request) {
	var upgrade = websocket.Upgrader{
		WriteBufferSize: 1024,
		ReadBufferSize: 1024,
	}
    
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
    
    var lastMod time.Time
	if n, err := strconv.ParseInt(r.FormValue("lastMod"), 16, 64); err == nil {
		lastMod = time.Unix(0, n)
	}
	
	// 写数据
	go Writer(conn, lastMod)
    // 读pong消息
	Reader(conn)
}

func Reader(conn *websocket.Conn) {
	
}

func Writer(conn *websocket.Conn) {
	
}
```



**（2）发送消息**

这一步完成 `Writer` 函数：

```go
func Writer(conn *websocket.Conn, lastMod time.Time) {
	var lastError string
	pingTicker := time.NewTicker(pingWait)
	writeTicker := time.NewTicker(writeWait)
	defer func() {
		pingTicker.Stop()
		writeTicker.Stop()
		if err := conn.Close(); err != nil {
			log.Println("conn.Close() Error: ", err.Error())
		}
	}()
	conn.SetPingHandler(PingHandler)
	select {
	case <-writeTicker.C:
		var p []byte
		var err error
		p, lastMod, err = readFileIfModified(lastMod)
		if err != nil {
			if s := err.Error(); s != lastError {
				lastError = s
				p = []byte(lastError)
			}
		} else {
			lastError = ""
		}
		if err := conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
			log.Println("conn.SetWriteDeadline Fail:  ", err.Error())
		}
		if err := conn.WriteMessage(websocket.TextMessage, p); err != nil {
			log.Println("conn.WriteMessage(websocket.BinaryMessage) Fail:  ", err.Error())
			return
		}
	case <-pingTicker.C:
		if err := conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
			log.Println("conn.SetWriteDeadline Fail:  ", err.Error())
		}
		if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
			log.Println("conn.WriteMessage(websocket.PingMessage) Fail:  ", err.Error())
			return
		}
	}
}

func readFileIfModified(lastMod time.Time) ([]byte, time.Time, error) {

}
```



**（3）完成 `readFileIfModified` 函数**

```go
func readFileIfModified(lastMod time.Time) ([]byte, time.Time, error) {
	// filename 通过 flag 库从全局变量中获取
	fi, err := os.Stat(filename)
	if err != nil {
		return nil, lastMod, err
	}
	if !fi.ModTime().After(lastMod) {
		return nil, lastMod, nil
	}
	p, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fi.ModTime(), err
	}
	return p, fi.ModTime(), nil
}
```



**（4）`ping` 消息处理函数**

```go

```



**（5）`Reader `函数监听控制消息**

```go
func Reader(conn *websocket.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("conn.Close() Error: ", err.Error())
		}
	}()
	
	conn.SetReadLimit(512)
	
	if err := conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println("conn.SetReadDeadline(time.Now().Add(pongWait)) Error: ", err.Error())
	}
	conn.SetPongHandler(PongHandle)
	for {
		_, _ , err := conn.ReadMessage()
		if err != nil {
			log.Println("conn.ReadMessage() fail Error:", err.Error())
		}
		return
	}
}

func PongHandle(appData string) error {
	return nil
}
```



**（6）`pong` 消息处理函数**

```go

```



**（7）编写 `main` 函数**

```go

```



**测试一下：**

```bash

```



#### Reference:

- https://www.jianshu.com/p/11ded5e80cdf
- https://godoc.org/github.com/gorilla/websocket
- https://github.com/gorilla/websocket

