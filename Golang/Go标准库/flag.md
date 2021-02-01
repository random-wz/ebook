> 和其它语言一样，go语言也提供了用来接收命令行参数的功能（flag包），我们可以使用flag包很方便的写一些自定义的命令，就像cobra包一样，这里向大家介绍flag包的使用方法，如果项目需要复杂或更高级的命令行解析方式，可以使用 [https://github.com/urfave/cli](https://link.jianshu.com/?t=https%3A%2F%2Fgithub.com%2Furfave%2Fcli) 或者 [https://github.com/spf13/cobra](https://link.jianshu.com/?t=https%3A%2F%2Fgithub.com%2Fspf13%2Fcobra) 这两个强大的库。希望对你有帮助。

##  一、flag自定义命令的步骤

flag包自定义的步骤很简单，一般只需要两步：

- 接收参数

- 解析参数
- 自定义方法处理命令

下面看个例子：

```go
package main

import (
	"flag"
	"fmt"
)

// 定义OK为bool类型，用来接收传入的参数
var OK bool

func main() {
    // 下面为获取命令行的参数并将值赋值给OK
    // 1) 参数名称为ok,默认值为false
    // 2) usage为print ok
    // 当使用-ok参数时，OK的值会变为true
	flag.BoolVar(&OK, "ok", false, "print ok")
    // 解析参数
	flag.Parse()
    // 自定义方法处理命令
	if OK {
		fmt.Println("OK")
	} else {
		fmt.Println("Not OK")
	}
}
```

我们来看一下测试结果：

```go
$ go run main.go -ok
OK
$ go run main.go
Not OK
```

## 二、使用flag获取命令行参数

#### 1. 获取命令行参数

**flag.Xxx()，其中 Xxx 可以是 Int、String，Bool 等；返回一个相应类型的指针，如：**

`var config =  flag.String("config", "config.cfg", "config file")`

这种方式会从命令行读取参数后，并将参数值返回，config会接收返回的值，String方法有三个参数（其他类型也一样）：

- 第一个参数：定义参数名称，命令行可以通过指定config的值来传递数据。
- 第二个参数：定义参数的默认值，如果命令行没有接收到config的值，则使用默认值。
- 第三个参数：定义参数的实用信息，也就是help信息。

**flag.XxxVar()，将 flag 绑定到一个变量上，如：**

```go
var config string
flag.StringVar(&config, "config", "config.cfg", "config file")
```

这种方式有四种参数，后面三种和上一种获取flag参数的含义相同，唯一的不同点是，在StringVar中，第一个参数为一个字符串类型的指针（其他类型类似），我们只需要提前定义好变量，并引入变量即可，实现的效果和上面那种方式是一样的。

#### 2. 自定义Value

**flag命令除了可以接收官方提供的参数类型外，我们可以自定义flag，只要实现 flag.Value 接口即可（要求 `receiver` 是指针），这时候可以通过如下方式定义该 flag：**

`flag.Var(&MyFlag, "name", "help message for flagname")`

**flag.Value接口：**

```go
type Value interface {
	String() string
	Set(string) error
}
```

flag.Value接口只有两个方法，我们只需要创建一个实例并实现这两个方法即可。

**自定义Value：**

在flag中对`Duration`这种非基本类型的支持，就是使用的类似的方式，我们来看一下他是怎么实现的：

（1）首先定义了一个time.Duration类型

```go
// -- time.Duration Value
type durationValue time.Duration
```

（2）通过newDurationValue函数new一个存放参数值的指针

```go
func newDurationValue(val time.Duration, p *time.Duration) *durationValue {
	*p = val
	return (*durationValue)(p)
}
```

（3）实现flag.Getter接口，flag.Getter接口继承了flag,Value接口

```go
func (d *durationValue) Set(s string) error {
	v, err := time.ParseDuration(s)
	if err != nil {
		err = errParse
	}
	*d = durationValue(v)
	return err
}

func (d *durationValue) Get() interface{} { return time.Duration(*d) }

func (d *durationValue) String() string { return (*time.Duration)(d).String() }

// Value is the interface to the dynamic value stored in a flag.
// (The default value is represented as a string.)
//
// If a Value has an IsBoolFlag() bool method returning true,
// the command-line parser makes -name equivalent to -name=true
// rather than using the next command-line argument.
//
// Set is called once, in command line order, for each flag present.
// The flag package may call the String method with a zero-valued receiver,
// such as a nil pointer.
type Value interface {
	String() string
	Set(string) error
}

// Getter is an interface that allows the contents of a Value to be retrieved.
// It wraps the Value interface, rather than being part of it, because it
// appeared after Go 1 and its compatibility rules. All Value types provided
// by this package satisfy the Getter interface.
type Getter interface {
	Value
	Get() interface{}
}
```

（4）通过CommandLine.Var实现参数的接收和传递

```go
// DurationVar defines a time.Duration flag with specified name, default value, and usage string.
// The argument p points to a time.Duration variable in which to store the value of the flag.
// The flag accepts a value acceptable to time.ParseDuration.
func DurationVar(p *time.Duration, name string, value time.Duration, usage string) {
	CommandLine.Var(newDurationValue(value, p), name, usage)
}

// Duration defines a time.Duration flag with specified name, default value, and usage string.
// The return value is the address of a time.Duration variable that stores the value of the flag.
// The flag accepts a value acceptable to time.ParseDuration.
func (f *FlagSet) Duration(name string, value time.Duration, usage string) *time.Duration {
	p := new(time.Duration)
	f.DurationVar(p, name, value, usage)
	return p
}
```

其实我们只需要创建一个实例并让他实现flag.Value接口就可以达到自定义类型的目的。

下面我们创建一个结构体类型的变量，并让他实现flag.Value接口：

```go
package main

import (
	"encoding/json"
	"flag"
	"fmt"
)

type student struct {
	Name string
	Age  int
}
type Student student

/*
flag.Value接口:
type Value interface {
	String() string
	Set(string) error
}
*/
// 这里可以设置自定义类型的默认值
func (s *Student) String() string {
	return fmt.Sprintf("%v", student{Name: "none"})
}

// Set方法用来将接收到的值解析出来
func (s *Student) Set(v string) error {
	fmt.Println("val:", v)
	return json.Unmarshal([]byte(v), &s)
}

// new一个Student的方法
func newStudentValue(value student, vPointer *student) *Student {
	*vPointer = value
	return (*Student)(vPointer)
}
func main() {
	var S student
	flag.Var(newStudentValue(student{}, &S), "student", "input student msg")
	flag.Parse()
	fmt.Println(S)
}
```

测试一下：

```bash
$ go run main.go --student='{"Name":"wang","Age":18}'
val: {"Name":"wang","Age":18}
{wang 18}
```

我们可以看到Student类型的值也可以正常解析。



## 三、解析flag参数

#### 1. 解析参数（Parse）

前面我们举例的时候已经用到了Parse方法，他会从参数列表中解析定义好的flag。

flag.Parse会调用CommandLine.Parse方法并使用命令行传入的参数并进行解析：

```go
// Parse parses the command-line flags from os.Args[1:]. Must be called
// after all flags are defined and before flags are accessed by the program.
func Parse() {
	// Ignore errors; CommandLine is set for ExitOnError.
	CommandLine.Parse(os.Args[1:])
}
```

<font color=red>注意：该方法应该在 flag 参数定义后而具体参数值被访问前调用。</font>

如果提供了 `-help` 参数（命令中给了）但没有定义（代码中没有），该方法返回 `ErrHelp` 错误。默认的 CommandLine，在 Parse 出错时会退出程序（ExitOnError）。

我们看一下Parse方法的源码：

```go
// Parse parses flag definitions from the argument list, which should not
// include the command name. Must be called after all flags in the FlagSet
// are defined and before flags are accessed by the program.
// The return value will be ErrHelp if -help or -h were set but not defined.
func (f *FlagSet) Parse(arguments []string) error {
	f.parsed = true
	f.args = arguments
	for {
		seen, err := f.parseOne()
		if seen {
			continue
		}
		if err == nil {
			break
		}
		switch f.errorHandling {
		case ContinueOnError:
			return err
		case ExitOnError:
			os.Exit(2)
		case PanicOnError:
			panic(err)
		}
	}
	return nil
}
```

我们可以看到，他其实执行的是不可导出的parseOne方法。

```go
func (f *FlagSet) parseOne() (bool, error)
```

parseOne会返回参数解析结果。

#### 2. 解析停止条件

通过查看parseOne方法的源码我们可以看到解析终止的条件有下面三种：

（1）参数列表长度为0

```go
	if len(f.args) == 0 {
		return false, nil
	}
```

在parseOne中每执行成功一次 parseOne，f.args 会少一个。所以，FlagSet 中的 args 最后留下来的就是所有 `non-flag` 参数。

```go
// it's a flag. does it have an argument?
f.args = f.args[1:]
```

（2）第一个 non-flag 参数

```go
if len(name) == 0 || name[0] == '-' || name[0] == '=' {
	return false, f.failf("bad flag syntax: %s", s)
}
```

也就是说遇到下面三种情况会停止解析：

- 无flag参数
- 接收到单独的`-`参数
- 接收到单独的`=`参数

（3）两个连续的"--"

```go
	numMinuses := 1
	if s[1] == '-' {
		numMinuses++
		if len(s) == 2 { // "--" terminates the flags
			f.args = f.args[1:]
			return false, nil
		}
	}
```

也就是，当遇到连续的两个"-"时，解析停止。但是有特殊情况。

如下面会正常解析：

```bash
$ command -c --
```

command命令的c参数接收到的值为"--"。

## 四、自定义Usage

在Linux系统中，所有的命令都会提供usage信息，那么在flag包中我们该如何定义呢？

其实很简单，flag包里面有一个Usage变量：

```go
var Usage = func() {
	fmt.Fprintf(CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	PrintDefaults()
}
```

我们可以重写Usage来定义Usage信息，下面举个例子：

```go
package main

import (
	"flag"
	"fmt"
)

// 自定义usage
func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", "This is a test.")
	flag.PrintDefaults()
}
func main() {
	var Config string
	flag.StringVar(&Config, "", "config.cfg", "config file")
	// 重写usage
	flag.Usage = usage
	flag.Parse()
}
```

测试一下：

```bash
$ go run main.go --help
Usage of This is a test.:
  - string
        config file (default "config.cfg")
exit status 2
```

我们可以看到输出了我们定义的usage信息。

## 五、其他方法

#### 1. Arg(i int) 和 Args()、NArg()、NFlag()

```go
// Arg返回第i个命令行参数。Arg（0）是剩余的第一个参数
// 在处理标志之后,如果请求的元素不存在，则Arg返回空字符串。
func Arg(i int) string {
	return CommandLine.Arg(i)
}

// Args返回非标志命令行参数。
func Args() []string { return CommandLine.args }

// NArg是处理标志后剩余的参数个数。
func NArg() int { return len(CommandLine.args) }

// NFlag返回已设置的命令行标志数。
func NFlag() int { return len(CommandLine.actual) }
```



#### 2. Visit/VisitAll

这两个函数分别用于访问 FlatSet 的 actual（存放参数值实际Flag的map） 和 formal（存放参数名默认Flag的map） 中的 Flag，而具体的访问方式由调用者决定。

```go
// VisitAll visits the flags in lexicographical order, calling fn for each.
// It visits all flags, even those not set.
func (f *FlagSet) VisitAll(fn func(*Flag)) {
	for _, flag := range sortFlags(f.formal) {
		fn(flag)
	}
}
```

```go
// Visit visits the flags in lexicographical order, calling fn for each.
// It visits only those flags that have been set.
func (f *FlagSet) Visit(fn func(*Flag)) {
	for _, flag := range sortFlags(f.actual) {
		fn(flag)
	}
}
```



#### 3. PrintDefaults()

看到这个名称是不是很熟悉，我们在重写usage函数的时候就有用到这个方法，他的作用是打印所有已定义参数的默认值（调用 VisitAll 实现），默认输出到标准错误，除非指定了 FlagSet 的 output（通过SetOutput() 设置）。

#### 4. SetOutput()

```go
// SetOutput 用来设置命令的usage和错误信息。
// 如果output为空, 则输出错误信息。
func (f *FlagSet) SetOutput(output io.Writer) {
	f.output = output
}
```

#### 5. Set(name, value string)

```go
// 将名称为name的flag的值设置为value, 成功返回nil。
func (f *FlagSet) Set(name, value string) error {
	flag, ok := f.formal[name]
	if !ok {
		return fmt.Errorf("no such flag -%v", name)
	}
	err := flag.Value.Set(value)
	if err != nil {
		return err
	}
	if f.actual == nil {
		f.actual = make(map[string]*Flag)
	}
	f.actual[name] = flag
	return nil
}
```



