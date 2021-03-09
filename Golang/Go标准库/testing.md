### 一、单元测试

> 单元测试（unit testing），是指对软件中的最小可测试单元进行检查和验证。对于单元测试中单元的含义，一般要根据实际情况去判定其具体含义，如C语言中单元指一个函数，Java 里单元指一个类，图形化的软件中可以指一个窗口或一个菜单等。总的来说，单元就是人为规定的最小的被测功能模块。
>
> 单元测试是在软件开发过程中要进行的最低级别的测试活动，软件的独立单元将在与程序的其他部分相隔离的情况下进行测试。

#### 1. go 语言中的单元测试

在 go 语言中通过文件明来标识测试文件，需要进行单元测试时，我们需要创建一个以`_test`结尾的文件，编辑完测试源码后，使用 `go test` 命令运行单元测试程序，下面是 `go test` 命令 几个常用的参数：

- -bench regexp 执行相应的 benchmarks，例如 -bench=.；
- -cover 开启测试覆盖率；
- -run regexp 只运行 regexp 匹配的函数，例如 -run=Array 那么就执行包含有 Array 开头的函数；
- -v 显示测试的详细命令。

单元测试源码文件可以由多个测试用例组成，每个测试用例函数需要以`Test`为前缀，例如：

```go
func TestXXX( t *testing.T )
```

单元测试需要注意的点：

- 测试用例文件不会参与正常源码编译，不会被包含到可执行文件中。
- 测试用例文件使用`go test`指令来执行，没有也不需要 main() 作为函数入口。所有在以`_test`结尾的源码内以`Test`开头的函数会自动被执行。
- 测试用例可以不传入 *testing.T 参数。

#### 2. 单元测试命令行

写一个简单的测试用例：

```go
// t_test.go
package main

import (
	"fmt"
	"testing"
)

func TestA(t *testing.T) {
	fmt.Println("Hello World")
}
```

使用 `go test` 命令运行单元测试程序：

```bash
$ go test -v t_test.go
=== RUN   TestA
Hello World
--- PASS: TestA (0.00s)
PASS
ok      command-line-arguments  0.281s
```

#### 3. 运行指定单元测试用例

在一个文件中写入多个单元测试函数：

```go
// t_test.go
package main

import (
	"fmt"
	"testing"
)

func TestA(t *testing.T) {
	fmt.Println("I am A")
}

func TestAB(t *testing.T) {
	fmt.Println("I am AB")
}

func TestB(t *testing.T) {
	fmt.Println("I am B")
}
```

使用 `run` 参数指定要运行的函数：

```go
$ go test -v -run TestA t_test.go
=== RUN   TestA
I am A
--- PASS: TestA (0.00s)
=== RUN   TestAB
I am AB
--- PASS: TestAB (0.00s)
PASS
ok      command-line-arguments  0.279s
```

<font color=red>注意：run 参数后面接的是正则表达式，能匹配上的所有单元测试函数都会运行，所以TestA 和 TestAB 都被执行了。</font>

#### 4. 标记单元测试结果

当需要终止当前测试用例时，可以使用 `FailNow` ，如下例：

```go
// t_test.go
package main

import (
	"fmt"
	"testing"
)

func TestA(t *testing.T) {
	fmt.Println("start")
	t.FailNow()
	fmt.Println("end")
}
```

执行单元测试程序：

```bash
$ go test -v t_test.go
=== RUN   TestA
start
--- FAIL: TestA (0.00s)
FAIL
FAIL    command-line-arguments  0.282s
```

我们也可以使用 `Fail` 只标记错误不终止测试，如下例：

```go
// t_test.go
package main

import (
	"fmt"
	"testing"
)

func TestA(t *testing.T) {
	fmt.Println("start")
	t.Fail()
	fmt.Println("end")
}
```

执行单元测试程序：

```bash
$ go test -v t_test.go
=== RUN   TestA
start
end
--- FAIL: TestA (0.00s)
FAIL
FAIL    command-line-arguments  0.298s
```



#### 5. 单元测试日志

每个测试用例可能并发执行，使用 testing.T 提供的日志输出可以保证日志跟随这个测试上下文一起打印输出。testing.T 提供了几种日志输出方法，详见下表所示。

| 方  法 | 说  明                           |
| ------ | -------------------------------- |
| Log    | 打印日志，同时结束测试           |
| Logf   | 格式化打印日志，同时结束测试     |
| Error  | 打印错误日志，同时结束测试       |
| Errorf | 格式化打印错误日志，同时结束测试 |
| Fatal  | 打印致命日志，同时结束测试       |
| Fatalf | 格式化打印致命日志，同时结束测试 |

#### 6. testing.T 更多用法



### 二、基准测试

