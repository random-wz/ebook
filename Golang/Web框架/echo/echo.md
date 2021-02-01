> echo

# 第一部分

## 一、安装

#### 1. 下载标准库

echo框架的安装其实很简单，下载官方的库就可以了：

`go get github.com/labstack/echo/v4`

当然你也可以选择其他版本进行下载，只需要将上面命令后面的v4改成你要下载的版本即可：

`go get github.com/labstack/echo/{version}`

Notice ：如果 go get 的时候报网络问题，如` unrecognized import path "golang.org/x/crypto" (https fetch: Get https://golang.org/x/crypto?go-get=1: dial tcp 216.239.37.1:443: connectex: A con
nection attempt failed because the connected party did not properly respond after a period of time, or established connection failed because connected host has failed to respond.)`

执行下面命令，设置proxy就可以了：

`export GOPROXY="https://goproxy.io"`

<font color=red>注意：这是在 git bash 中的命令，cmd 中使用 set 命令。</font>

#### 2. Hello World

按照惯例，完成包的安装后，我们实现一个 Hello World 程序。

```go
package main


import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
```

运行程序：

```bash
$ go run main.go

   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v4.1.16
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
⇨ http server started on [::]:1323
```

访问127.0.0.1:1323 

```bash
$ curl 127.0.0.1:1323
Hello, World!
```

