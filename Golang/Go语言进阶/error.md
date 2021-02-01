### 一、你不知道的Error

#### 1. 异常处理

在Java、C++、Python等语言中，通常通过抛出异常及捕获异常的方式获取处理程序中的报错信息，比如在C++中引入了exception，但是无法知道被调用方会抛出什么异常，在Java中引入了 checked exception，方法的所有者必须声明，调用者必须处理。在启动时抛出大量的异常是司空见惯的事情，并在他们的调用堆栈中尽职地记录下来。Java异常不再是异常，而是变的司空见惯了。他们从良性到灾难性都有使用，异常的严重性由u函数的调用者来区分。

而Go语言的处理异常的逻辑是不引入 exception，支持多参数返回，所以可以很容易的在函数的签名中带上实现了 **error interface** 的对象，交由调用者来判定。另外 Go 语言中还提供了 **panic** 机制，panic 意味着 **fatal error**，也就是说程序挂了，不能假设调用者来解决 panic，意味着代码不能再继续运行了。

使用多个返回值和一个简单的约定，Go 解决了让程序员知道什么时候出了问题，并为真正的异常情况保留了 panic。

#### 2. Go 语言中的 Error

我们来看一下 Go 语言中的 error 内置类型：

```go
// The error built-in interface type is the conventional interface for
// representing an error condition, with the nil value representing no error.
type error interface {
	Error() string
}
```

它其实就是一个实现了 Error 方法的 interface，也就是说在 Go 语言中 **error are values**。

#### 3. Exception vs Error

关于一篇微软工程博客的博文这样评价 expection：`expection 的设计并不是糟糕的，但是想要很好的处理 expection 是非常困难的。`，而 Go 语言中通过 error ，让人们可以通过检查 error 的值来进行一些逻辑操作，对开发人员来说十分简单而且友好。

对于真正意外的情况，那些表示不可恢复的程序错误，例如索引越界、不可恢复的环境问题、栈溢出，我们才使用 panic 。因此 Go 语言中的 error 具有以下特点：

- 简单
- 考虑失败，而不是成功（Plan for failure, not suceee）
- 没有隐藏的控制流
- 完全交给你来控制 error
- Error are values

### 二、常见的几种Error类型

#### 1. Sentinel Error

**Sentinel Error** 也叫预定义的特定错误，这个名字来源于计算机编程中使用一个特定值来表示不可能进一步处理的方法。这种定义错误的方式十分常见，比如下面的例子：

```go
// ErrShortWrite means that a write accepted fewer bytes than requested
// but failed to return an explicit error.
var ErrShortWrite = errors.New("short write")

// ErrShortBuffer means that a read required a longer buffer than was provided.
var ErrShortBuffer = errors.New("short buffer")

// EOF is the error returned by Read when no more input is available.
var EOF = errors.New("EOF")

// ErrUnexpectedEOF means that EOF was encountered in the
// middle of reading a fixed-size block or data structure.
var ErrUnexpectedEOF = errors.New("unexpected EOF")

// ErrNoProgress is returned by some clients of an io.Reader when
// many calls to Read have failed to return any data or error,
// usually the sign of a broken io.Reader implementation.
var ErrNoProgress = errors.New("multiple Read calls return no data or error")
```

使用 sentinel 值是非常不灵活的错误处理策略，因为调用方必须使用 == 将结果与预先声明的值进行比较，当你想要提供更多的上下文时，这就出现了一个问题，因为返回一个不同的错误将破坏相等性的检查，甚至是一些有意义的 fmt.Errof 携带一些上下文，也会破坏调用者的检查，调用者将被迫查看 error.Error()方法的输出，以查看它是否与特定的字符串匹配。

而在 Go 语言中我们认为错误处理不应该依赖于检测 error.Error 的输出， Error 方法存在的意义主要是方便程序员使用，但不是程序（编写测试可能会依赖这个返回）。这个输出的字符串用于记录日志、输出到标准输出等。

因此引入 Sentinel Error 带来了下面几个问题：

- Sentinel errors 称为你API公共部分。

  > 如果你的公共函数或方法返回一个特定值的错误，那么该值必须是公共的，那么就需要文档记录，这回增加API的表面积。
  > 如果API定义了一个返回特定错误的 interface，则该接口的左右实现都将被限制为仅返回该错误，即使他们可以提供更具描述性的错误。比如 io.Reader，像io.Copy这类函数需要 reader的实现者比如返回 io.EOF来告诉调用者没有更多的数据了，但这又不是错误。

- Sentinel errors 在两个包之间创建了依赖。

  >sentinel errors 最糟糕的问题是他们在两个包之间创建了源代码依赖关系。比如如果在项目中有许多包导出错误值时存在耦合，项目中的其他包必须导入这些错误值才能检查特定的错误条件（in the form of an import loop）

因此我们要尽可能的避免使用 sentinel errors，虽然标准库中有一些使用他们的情况，但这不是我们应该模仿的模式。

#### 2. Error types

在 Go 语言中 error 类型实际上是一个实现了 Error()方法的接口，因此我们可以自定义一个类型，然后让他实现 error 接口，这就是一个 **Error type** 了。如下面的例子，我们定义一个 MyError 对象，并实现 error 接口：

```go
type MyError struct {
	Msg  string
	Line string
	File string
}

func (e *MyError) Error() string {
	return fmt.Sprintf("File: %s\tLine:%s\tMsg:%s", e.File, e.Line, e.Msg)
}
```

通过定义一个实现了 error 接口的对象，我们可以传递更多的上下文信息，因为MyError是一个 type，因此调用者可以使用断言转换成这个类型来获取更多的上下文信息。如下面的例子：

```go
func ReadFile(file string) error {
	return &MyError{Msg: "File not exist", Line: "1", File: "main.go"}
}

func main() {
	err := ReadFile("file")
	switch err.(type) {
	case *MyError:
		fmt.Println(err.(*MyError).Line)
		// do something
	case nil:
		// do something
	default:
		// do something
	}
}
```

与错误值相比，错误类型的一大改进是他能够包装底层错误以提供更多的上下文信息。

调用者要使用类型断言获取它的上下文信息就要让自定义的 error 变为 public，这种模型会导致和调用者产生强耦合，从而导致 API 变的脆弱，所以要尽量避免使用 error types，或者至少避免将它们作为公共 API 的一部分，虽然错误类型比 sentinel errors 更好，因为它们可以捕获关于出错的更多上下文，但是 error types 共享 error values 许多相同的问题。

#### 3. Opaque errors

**Opaque errors**（不透明的错误处理）要求代码和调用者之间的耦合最少，作为调用者，关于操作的结果你只需要知道的就是它是否起作用了（成功还是失败），只返回错误而不假设其内容，如下面的例子：

```go
func QpaqueError() error{
	if _, err := os.Getwd(); err != nil {
		return err
	}
	return nil
}
```

对于错误处理我们应该``断言行为错误，而不是类型错误``，比如与进程外的世界进行交互(如网络活动)，需要调用方调查错误的性质，以确定重试该操作是否合理。在这种情况下，我们可以断言错误实现了特定的行为，而不是断言错误是特定的类型或值。如下面的例子：

```go
package main

import (
	"fmt"
)

type MyError struct {
	Msg  string
	Line string
	File string
}

type Error interface {
	error
	Timeout() bool
}

func (e *MyError) Error() string {
	return fmt.Sprintf("File: %s\tLine:%s\tMsg:%s", e.File, e.Line, e.Msg)
}

func (e *MyError) Timeout() bool {
	return false
}

func DoSomething() Error {
	return &MyError{}
}
func main() {
	err := DoSomething()
	if e, ok := err.(Error); ok && e.Timeout() {
		// do something
	}
}
```

这里的关键是，这个逻辑可以在不导入定义错误的包或者实际上不了解 err 的底层类型的情况下实现——我们只对它的行为感兴趣。

### 三、恰如其分的处理Error

#### 1. 处理error的常见方式

在 Go 语言程序中，正常的代码流程为一条直线而不是缩进的代码，如下面的代码：

```go
f, err := os.Open("test.txt")
if err != nil {
	// do something
}

if err := f.Close(); err != nil {
	// do something
}
```

有时候我们在开发的时候会调用一个方法，该方法执行后的返回值与当前方法一致，如下面的代码：

```go
func Auth(username, password string) error{
	if err := auth(username, password); err != nil {
		return err
	}
	return nil
}
```

我们可能会将它简化为：

```go
func Auth(username, password string) error{
	return auth(username, password)
}
```

代码是简洁了，但这样做调用者可能会一层层的将错误返回给程序的下一级，直到程序的顶级，程序的主体会把错误打印标准输出或日志文件中，打印出来的只有一条报错信息，并没有记录没有导致错误的调用堆栈的堆栈跟踪信息。于是你可能会这样做：

```go
func Auth(username, password string) error{
    if err := auth(username, password); err != nil {
		return fmt.Errorf("authenticate fail: %v", err)
	}
	return nil
}
```

但是正如我们前面看到的，这种模式与 sentinel errors 或 type assertions 的使用不兼容，因为将错误值转换为字符串，将其与另一个字符串合并，然后将其转换回 fmt.Errorf 破坏了原始错误，导致等值判定失败。而对于错误信息你应该只处理一次错误，处理错误意味着检查错误值，并做出一个单独的决定。

#### 2.  wrap error

对于错误的处理我们应该遵守下面的规则：

- 报错信息要被日志记录

- 应用程序处理错误，保证**100%**完整性。

- 之后不再报告当前错误，也就是说错误只被处理一次。

而  **wrap error** 可以帮我们记录错误信息的堆栈信息，让我们可以更好的处理和定位程序中的错误信息。例如下面的程序：

```go
package main

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
)

func OpenFile(file string) (*os.File, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.Wrapf(err, "open file %s fail", file)
	}
	return f, nil
}

func ReadFile(file string) (interface{}, error) {
	f, err := OpenFile(file)
	if err != nil {
		return nil, errors.Wrap(err, "read file fail")
	}
	var buf []byte
	_, err = f.Read(buf)
	if err != nil {
		return nil, errors.Wrapf(err, "f.Read file")
	}
	return string(buf), nil
}

func main() {
	// test.txt文件不存在所以会报错
	if _, err := ReadFile("test.txt"); err != nil {
		fmt.Printf("original error: %T %v\n", errors.Cause(err), errors.Cause(err))
        // 注意这里使用%+v可以打印出堆栈信息
		fmt.Printf("stack trace: \n%+v\n", err)
		return
	}
	return
}
```

运行之后，我们可以看到报错信息的堆栈信息，很容易定位错误：

```go
$ go run main.go
original error: *os.PathError open test.txt: The system cannot find the file specified.
stack trace:
open test.txt: The system cannot find the file specified.
open file test.txt fail
main.OpenFile
        F:/A_StudyDocument/Code/Golang/Study/study/main.go:12
main.ReadFile
        F:/A_StudyDocument/Code/Golang/Study/study/main.go:18
main.main
        F:/A_StudyDocument/Code/Golang/Study/study/main.go:31
runtime.main
        c:/Go/src/runtime/proc.go:200
runtime.goexit
        c:/Go/src/runtime/asm_amd64.s:1337
read file fail
main.ReadFile
        F:/A_StudyDocument/Code/Golang/Study/study/main.go:20
main.main
        F:/A_StudyDocument/Code/Golang/Study/study/main.go:31
runtime.main
        c:/Go/src/runtime/proc.go:200
runtime.goexit
        c:/Go/src/runtime/asm_amd64.s:1337
```

#### 3. 介绍一个好用的Error包

`github.com/pkg/errors` 包是一个专门用来处理错误信息的第三方包，第二步的例子中我也是使用了这个包，安装这个包到的方式很简单：

```bash
go get "github.com/pkg/errors"
```

`github.com/pkg/errors` 包的八个方法：

- `New(message string) error `

  ```go
  // 返回一个error,并附加堆栈信息
  ```

- ` Wrap(err error, message string) error`

  ```go
  // Wrap返回一个错误，在调用Wrap时使用堆栈跟踪对err进行注释，并提供消息。如果err为nil，Wrap返回nil。
  ```

- `Wrapf(err error, format string, args ...interface{}) error`

  ```go
  // 与Wrap类似
  ```

- `Cause(err error) error`

  ```go
  // 如果错误对象实现了下面到的接口，Cause返回错误的根本原因
  //     type causer interface {
  //            Cause() error
  //     }
  // 如果错误没有实现Cause，将返回原始错误。如果错误为nil，则返回nil，无需进一步调查。
  ```

- `WithMessage(err error, message string) error`

  ```go
  // WithMessage使用message注释报错信息。如果err为nil，则WithMessage返回nil。
  ```

- `WithMessagef(err error, format string, args ...interface{}) error`

  ```go
  // 与WithMessage类似
  ```

- `Errorf(format string, args ...interface{}) error`

  ```go
  // Errorf根据格式说明符格式化，并将字符串作为满足error的值返回。
  // Errorf还记录调用时的堆栈跟踪。
  ```

- `WithStack(err error) error`

  ```go
  // WithStack在调用WithStack的点处使用堆栈跟踪对err进行注释。
  // 如果err为nil，则WithStack返回nil。
  ```

####  4. Go1.13处理Error的新方法

Go 1.13版本 errors 包提供了三个新的方法用来支持 wrap error，`github.com/pkg/errors` 也提供了下面三个方法：

- `Is(err, target error) bool`

  ```go
  // Is报告err链中是否有任何错误与目标匹配。
  // 这个链由err本身和通过反复调用Unwrap获得的错误序列组成。
  // 如果一个错误等于某个目标，或者它实现了一个方法Is（error）bool，如果Is（target）返回true，则认为该错误与该目标匹配。
  ```

- `As(err error, target interface{}) bool`

  ```go
  // As查找err链中与target匹配的第一个错误，如果是，则将target设置为该错误值并返回true。这个链由err本身和通过反复调用Unwrap获得的错误序列组成。
  // 如果错误的具体值可分配给target所指向的值，或者如果错误的方法为（interface{}）bool，例如As（target）返回true，则错误与target匹配。在后一种情况下，As方法负责设置目标。
  // 如果目标不是一个指向实现错误的类型或任何接口类型的非nil指针，则As将panic。如果err为nil，As返回false。
  ```

- `Unwrap(err error) error`

  ```go
  // 如果err的类型包含返回错误的Unwrap方法，则Unwrap返回对err调用Unwrap方法的结果。
  // 否则，Unwrap返回nil。
  ```


