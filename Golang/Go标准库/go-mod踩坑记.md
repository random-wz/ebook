> 这篇博客向大家介绍，我再使用Go mod进行包管理中遇到的问题，持续更新，希望对你有帮助。

#### 1.  Go mod 下载包报unrecognized的问题

今天在一个开源项目上使用go mod进行包管理的时候遇到了下面的问题：

```bash
golang.org/x/sys@v0.0.0-20200615200032-f1bc736245b1: unrecognized import path "golang.org/x/sys"
```

遇到类似的问题，我们可以通过匹配值proxy的方式进行处理（类似的方法，应该也可以，不一定必须是sys包报错）

```bash
$ export GOPROXY="https://goproxy.io"
```

我是windows从操作系统，因此修改环境变量是export，这种只是临时设置，如果要永久设置就需要该环境变量了，比较简自行百度吧。

#### 2.  unexpected module path

我在指定go mod tidy的时候报了个unexpected module path错误，如下：

```bash
go: github.com/go-resty/resty@v1.12.0: parsing go.mod: unexpected module path "gopkg.in/resty.v1"
```

最后在stackoverflow上面看到一个类似的问题，这种问题一般是因为新版本改了包名导致的，因此我们需要在go.mod文件中使用replace修改包名：

```bash
replace github.com/go-resty/resty v1.12.0 => gopkg.in/resty.v1 v1.12.0
```

修改完成再次`go mod tidy`就ok了。