### 一、使用Go Module做包管理

> 在项目上，经常需要使用外部的包来实现相应的功能，当多个人同时进行开发时，就需要一个包管理器，每个人通过包管理器将自己代码所调用包，保存在同一个目录地下，然后提交到git上面，当以通进行开发的人拉取代码后，就不需要担心本地没有你的依赖包的问题，从而确保项目开发正常进行，下面介绍golang1.11版本后的go mod使用方法。

#### 1.与go mod相关的环境变量

GO111MODULE环境变量，你可以将这个环境变量配置为：auto、on、off，默认为auto，默认情况下你通过go get下载的包会自动保存在GOPATH目录地下的pkg目录里面，这样当你import的时候只需要写github.com/xxxx，当你与同事共同开发项目的时候，如果你所更新的代码中导入了新的包，但是你同事那边没有，那他想要运行代码就必须把你go get的数据包全部go get，然后才能继续开发，相当麻烦。

#### 2.上面讲了应用场景，现在说一下解决方法

go mod有两种应用场景，一种是在GOPATH目录地下进行开发，还有一种就是在GOPATH目录以外的目录进行开发。

第一种：在GOPATH目录地下进行开发

这种情况下，你需要将GO111MODULE环境变量设置为on，然后执行下面步骤：

- 在项目目录地下执行：go mod init name 初始化modules，初始完成后会在目录下生成一个go.mod文件，里面的内容只有一行“module test”。注意：name可以为任意名称，一般为项目名称。
- 在项目目录底下执行：go mod tidy 第一次执行会生成一个.sum文件，这条命令会自动更新依赖关系，并且将包下载放入cache。
- 在项目目录地下执行：go mod vendor 这时go mod会创建一个vendor目录，然后自动将代码中所调用的包保存到vendor目录里面。

第二种：不再GOPATH目录地下进行开发

这种情况下，你需要将GO111MODULE环境变量设置为auto或off，然后执行的步骤与第一种情况相同。

然后你会发现你的代码不再GOPATH中依然可以运行和编译。

### 二、Go Module踩坑记

> 这篇博客向大家介绍，我再使用Go mod进行包管理中遇到的问题，持续更新，希望对你有帮助。

#### 1.  Go mod 下载包报unrecognized

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