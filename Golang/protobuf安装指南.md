> 今天在学习 protocol buffer 的时候遇到了 protoc-gen-go not find 的问题，我是在 windows10 系统中安装的 protoc ，在官网下载了安装包（包虽然很小，但是下载是真的慢），然后安装，以为万事大吉，结果报依赖的问题，这里记录下解决方法，希望对大叫有帮助。

#### 1. 先来看问题

```bash
$ protoc --go_out=. hello.proto
'protoc-gen-go' 不是内部或外部命令，也不是可运行的程序
或批处理文件。
--go_out: protoc-gen-go: Plugin failed with status code 1.
```

#### 2. 遇到问题第一时间问度娘

找到一篇博客说，是因为缺少包，那就安装包。

试了网上的教程：

```bash
go get -u github.com/golang/protobuf/protoc-gen-go
```

但是因为网络原因失败了：

```bash
$ go get -u github.com/golang/protobuf/protoc-gen-go
# cd C:\GOPATH\src\google.golang.org\protobuf; git pull --ff-only
fatal: unable to access 'https://go.googlesource.com/protobuf/': Failed to connect to go.googlesource.com port 443: Timed out
package google.golang.org/protobuf/compiler/protogen: exit status 1
```

网上说可以使用 gopm（类似于 npm ）命令可以解决此类问题，但是。。。

```bash
$ gopm get -g  github.com/golang/protobuf/protoc-gen-go
[GOPM] ?[36m08-24 10:08:30?[0m [?[31mERROR?[0m] github.com/golang/protobuf/protoc-gen-go: fail to make request: Get https://gopm.io/api/v1/revision?pkgname=github.com/golang/protobuf: dial tcp: lookup gopm.
io: no such host
```

我决定放弃了，直接找源码吧。

#### 3. 编译源码

直接访问 Github 上面 [protobuf](https://github.com/golang/protobuf) 的源码，使用 git 下载到本地：

```bash
git clone git@github.com:golang/protobuf.git
```

源码使用 go module 模式进行的包管理，初始化包名为 github.com/golang/protobuf ，为了方便调用我改成了 myprotobuf ，接着就是 `go mod tidy && go mod vendor`下载依赖包了，进入 protoc-gen-go 目录，很奇怪`github.com/golang/protobuf/internal/gengogrpc `这个包没有下载下来，但是我发现源码里面就有这个包，那就改一下导入包的路径吧：`myprotobuf/internal/gengogrpc`，问题成功解决。

接下来就是进入 `protoc-gen-go` 编译源码生成可执行文件了：

```go
random@random-wz MINGW64 /d/GOCODE/protobuf/protoc-gen-go (master)
$ go install
```

很顺利，接下来就是测试了。

#### 4. 测试

创建一个 proto 文件：

```protobuf
syntax = "proto3";

package main;

message String {
    string value = 1;
}
```

生成 go 代码：

```bash
$ protoc --go_out=. hello.proto
$ ls
hello.pb.go  hello.proto
```

大功告成！！！！