> `wire`是 Google 开源的一个依赖注入工具。它是一个代码生成器，并不是一个框架。我们只需要在一个特殊的`go`文件中告诉`wire`类型之间的依赖关系，它会自动帮我们生成代码，帮助我们创建指定类型的对象，并组装它的依赖。

#### 1. 安装 wire

执行下面的命令会将 `wire` 安装到 GOPATH 下的 bin 目录：

```bash
go get github.com/google/wire/cmd/wire
```



##### Reference:

- [wire 库](https://pkg.go.dev/github.com/google/wire)

- [Compile-time Dependency Injection With Go Cloud's Wire](https://blog.golang.org/wire)

- [Go 每日一库之 wire](https://zhuanlan.zhihu.com/p/110453784)