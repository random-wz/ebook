> 因为墙的原因导致安装bee失败的解决办法，总有一个适合你。

先来看两种常见的报错：

第一种：

```bash
$ go get -u github.com/beego/bee
# cd C:\GOPATH\src\golang.org\x\text; git pull --ff-only
fatal: unable to access 'https://go.googlesource.com/text/': Failed to connect to go.googlesource.com port 443: Timed out
package golang.org/x/text/transform: exit status 1
```

第二种：

```bash
$ go get github.com/beego/bee
# github.com/gadelkareem/delve/service/debugger
..\github.com\gadelkareem\delve\service\debugger\debugger.go:129:3: cannot use logger (type *"github.com/go-delve/delve/vendor/github.com/sirupsen/logrus".Entry) as type *"github.com/gadelkareem/delve/vendo
r/github.com/sirupsen/logrus".Entry in field value
# github.com/gadelkareem/delve/service/rpccommon
..\github.com\gadelkareem\delve\service\rpccommon\server.go:83:3: cannot use logger (type *"github.com/go-delve/delve/vendor/github.com/sirupsen/logrus".Entry) as type *"github.com/gadelkareem/delve/vendor/
github.com/sirupsen/logrus".Entry in field value
```

当然这是我安装的时候的报错，可能与你的报错略有不同，但问题原因是一致的，那就是防火墙，网上有很多种解决办法这里总结一下，但是我这边环境试了都不行，最后发现了一种新的解决办法，这里向大家详细介绍一下：

#### 1. 翻墙

这里就不详述了，翻墙后，直接安装官网教程安装即可，但注意要科学上网，哈哈哈

#### 2. 设置代理

首先更改golang的配置网上有两种配置方法（这里以windows系统为例）：

```bash
go env -w GOPROXY=https://goproxy.io,direct
go env -w GO111MODULE=on
```

或者：

```bash
set GO111MODULE=on
set GOPROXY=https://goproxy.io
```

再次执行`go get -u github.com/beego/bee `命令，就可以安装成功了，测试一下：

```bash
$ bee version
______
| ___ \
| |_/ /  ___   ___
| ___ \ / _ \ / _ \
| |_/ /|  __/|  __/
\____/  \___| \___| v1.12.0

├── Beego     : 1.12.2
├── GoVersion : go1.12.5
├── GOOS      : windows
├── GOARCH    : amd64
├── NumCPU    : 8
├── GOPATH    : C:\GOPATH
├── GOROOT    : c:\Go
├── Compiler  : gc
└── Date      : Friday, 9 Oct 2020
```

很遗憾我并没有安装成功，依旧报错。



#### 3. 手动安装

首先我们要知道 bee 并不是用来在项目中实现功能的库，他是一个为了协助快速开发 beego 项目而创建的项目，通过 bee 您可以很容易的进行 beego 项目的创建、热编译、开发、测试、和部署，归根结底就是一个应用程序，我们通过改应用程序管理 beego 项目，既然使用官方的安装方式网络不同，那么我们就自己手动安装，那如何安装呢？

首先我们在 Github 上可以找到 bee 项目的源码，我们将源码下载下来：

`git clone git@github.com:beego/bee.git`

源码是通过 Go Module 进行包管理的，我们执行 `go mod vendor`下载依赖包，下载完成后，直接编译源码：

`go build -o bee`

编译完成之后测试一下：

```bash
$ bee version
______
| ___ \
| |_/ /  ___   ___
| ___ \ / _ \ / _ \
| |_/ /|  __/|  __/
\____/  \___| \___| v1.12.0

├── Beego     : 1.12.2
├── GoVersion : go1.12.5
├── GOOS      : windows
├── GOARCH    : amd64
├── NumCPU    : 8
├── GOPATH    : C:\GOPATH
├── GOROOT    : c:\Go
├── Compiler  : gc
└── Date      : Friday, 9 Oct 2020
```

发现是可以用的，然后我们将 bee 可执行程序放到 GOPATH 目录即可。

