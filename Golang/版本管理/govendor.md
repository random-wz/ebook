> 在 Go 语言开发工具中，有许多包管理器，之前介绍了 go module 的使用，这也是我个人比较推荐的包管理器，但是在项目上碰到了使用govendor做包管理器，于是学习了一下，希望对你有帮助。

#### 1. govendor的安装

执行下面的命令即可安装：

```bash
go get -u -v github.com/kardianos/govendor
```

安装后会生成一个个执行文件，该文件在**GOPATH/bin/govendor**，建议将**bin**目录加入到环境变量中，这样方便使用安装的各种工具，关于如何配置，比较简单大家自己百度吧。

#### 2. govendor命令

govendor命令语法很简单，如下：

```bash
govendor COMMAND
```

下面是常用的功能列表：

| 命令           | 功能                                                         |
| :------------- | ------------------------------------------------------------ |
| `init`         | 初始化 vendor 目录                                           |
| `list`         | 列出所有的依赖包                                             |
| `add`          | 添加包到 vendor 目录，如 govendor add +external 添加所有外部包 |
| `add PKG_PATH` | 添加指定的依赖包到 vendor 目录                               |
| `update`       | 从 $GOPATH 更新依赖包到 vendor 目录                          |
| `remove`       | 从 vendor 管理中删除依赖                                     |
| `status`       | 列出所有缺失、过期和修改过的包                               |
| `fetch`        | 添加或更新包到本地 vendor 目录                               |
| `sync`         | 本地存在 vendor.json 时候拉去依赖包，匹配所记录的版本        |
| `get`          | 类似 `go get` 目录，拉取依赖包到 vendor 目录                 |

govendor管理的包类型列表：

| 状态      | 缩写状态 | 含义                                                   |
| --------- | -------- | ------------------------------------------------------ |
| +missing  | +m       | 代码引入了该包，但是没有找到该包。                     |
| +program  | +p       | 主程序包，意味着可以编译为可执行文件。                 |
| +local    | +l       | 本地包，即项目中自己写的包。                           |
| +external | +e       | 外部包，即被$GOPATH管理，但不在 vender 目录下的包。    |
| +vender   | +v       | 已被 govendor管理，即在 vendor 目录下的包。            |
| +std      | +s       | 标准库中的包。                                         |
| +unused   | +u       | 未使用的包，也就是说包在vendor目录下，但项目没有用到。 |
| +outside  |          | 外部包和缺失的包。                                     |
| +all      |          | 所有的包。                                             |

#### 3. 实践

创建一个新项目`govendor-test`，项目中只有一个`main.go`文件，为了演示govendor的包管理功能，只匿名导入包：

```go
package main

import (
	"fmt"
	_ "github.com/gin-gonic/gin"
	_ "github.com/sirupsen/logrus"
	_ "github.com/spf13/viper"
)

func main() {
	fmt.Println("Hello World!")
}
```

执行`govendor init`初始化，初始化之后会在项目目录下面生成一个vendor目录：

```bash
$ govendor init
$ ls
main.go  vendor/
```

在 verdor 目录下有一个 vendor.json 文件，这个文件用来包的版本信息：

```json
{
	"comment": "",
	"ignore": "test",
	"package": [],
	"rootPath": "govendor-test"
}
```

将外部包导入 vendor 目录：

```bash
$ govendor add +e
```

此时 vendor.json 文件中已经记录了所有的外部包信息了：

```json
{
	"comment": "",
	"ignore": "test",
	"package": [
		{
			"checksumSHA1": "auxm2wsouwNJiA8vZNEmuFh4ICQ=",
			"path": "github.com/fsnotify/fsnotify",
			"revision": "7f4cf4dd2b522a984eaca51d1ccee54101d3414a",
			"revisionTime": "2020-04-17T21:56:12Z"
		},
		{
			"checksumSHA1": "mKall8xfKBPL/a5Ji//i968rLsY=",
			"path": "github.com/gin-contrib/sse",
			"revision": "54d8467d122d380a14768b6b4e5cd7ca4755938f",
			"revisionTime": "2019-06-02T15:02:53Z"
		},
		{
			"checksumSHA1": "9NnSQExpNJSbqIV7Ta0juBLEl7Q=",
			"path": "github.com/gin-gonic/gin",
			"revision": "e602d524cccad90261e10bbb5ca41e9a81e467d4",
			"revisionTime": "2019-07-03T23:57:52Z"
		},
		{
			"checksumSHA1": "7B85He9fYRLpGHDd3R3q1qQhP0Y=",
			"path": "github.com/gin-gonic/gin/binding",
			"revision": "e602d524cccad90261e10bbb5ca41e9a81e467d4",
			"revisionTime": "2019-07-03T23:57:52Z"
		},
        ...
   ]
}
```

下来我们来看如何让别人使用我们的包，首先删除 vendor 目录下已经下载的依赖包，这时就只有一个 vendor.json 文件。执行`govendor sync同步包`：

```bash
$ govendor sync -v
```

执行完成后可以看到项目需要的包都同步到 vendor 目录下了。