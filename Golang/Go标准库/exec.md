> os/exec包提供了执行外部命令的方法，它包装了os.StartProcess函数以便更容易的修正输入和输出，使用管道连接I/O，以及作其它的一些调整。这里记录以下os/eexec包的学习笔记，希望对你有帮助。



#### 1. 执行外部命令

外部执行命令都是Cmd对象的方法，我们先了解一下Cmd对象：

```go
type Cmd struct {
    // Path是将要执行的命令的路径。
    //
    // 该字段不能为空，如为相对路径会相对于Dir字段。
    Path string
    // Args保管命令的参数，包括命令名作为第一个参数；如果为空切片或者nil，相当于无参数命令。
    //
    // 典型用法下，Path和Args都应被Command函数设定。
    Args []string
    // Env指定进程的环境，如为nil，则是在当前进程的环境下执行。
    Env []string
    // Dir指定命令的工作目录。如为空字符串，会在调用者的进程当前目录下执行。
    Dir string
    // Stdin指定进程的标准输入，如为nil，进程会从空设备读取（os.DevNull）
    Stdin io.Reader
    // Stdout和Stderr指定进程的标准输出和标准错误输出。
    //
    // 如果任一个为nil，Run方法会将对应的文件描述符关联到空设备（os.DevNull）
    //
    // 如果两个字段相同，同一时间最多有一个线程可以写入。
    Stdout io.Writer
    Stderr io.Writer
    // ExtraFiles指定额外被新进程继承的已打开文件流，不包括标准输入、标准输出、标准错误输出。
    // 如果本字段非nil，entry i会变成文件描述符3+i。
    //
    // BUG: 在OS X 10.6系统中，子进程可能会继承不期望的文件描述符。
    // http://golang.org/issue/2603
    ExtraFiles []*os.File
    // SysProcAttr保管可选的、各操作系统特定的sys执行属性。
    // Run方法会将它作为os.ProcAttr的Sys字段传递给os.StartProcess函数。
    SysProcAttr *syscall.SysProcAttr
    // Process是底层的，只执行一次的进程。
    Process *os.Process
    // ProcessState包含一个已经存在的进程的信息，只有在调用Wait或Run后才可用。
    ProcessState *os.ProcessState
    // 内含隐藏或非导出字段
}
```

Cmd代表一个正在准备或者在执行中的外部命令。

- 执行指定的程序

  ```go
  // 函数返回一个*Cmd，用于使用给出的参数执行name指定的程序。
  // 返回值只设定了Path和Args两个参数。
  // 如果name不含路径分隔符，将使用LookPath获取完整路径；否则直接使用name。
  // 参数arg不应包含命令名。
  func Command(name string, arg ...string) *Cmd
  ```

  Example:

  ```go
  package main
  
  import (
  	"bytes"
  	"fmt"
  	"log"
  	"os/exec"
  	"strings"
  )
  
  func main() {
  	// 执行: tr a-z A-Z "Hello World!"
  	cmd := exec.Command("tr", "a-z", "A-Z")
  	// 定义标准输入的信息为: Hello World
  	cmd.Stdin = strings.NewReader("Hello World")
  	// out用来存储标准输出信息
  	var out bytes.Buffer
  	cmd.Stdout = &out
  	// 执行命令
  	err := cmd.Run()
  	if err != nil {
  		log.Fatal(err)
  	}
  	// out.String()输出命令执行结果
  	fmt.Printf("in all caps: %q\n", out.String())
  }
  ```

  Output:

  ```bash
  $ go run main.go
  in all caps: "HELLO WORLD"
  ```

  

- 标准输入

  ```go
  // StdinPipe方法返回一个在命令Start后与命令标准输入关联的管道。
  // Wait方法获知命令结束后会关闭这个管道。
  // 必要时调用者可以调用Close方法来强行关闭管道，例如命令在输入关闭后才会执行返回时需要显式关闭管道。
  func (c *Cmd) StdinPipe() (io.WriteCloser, error)
  ```

  

- 标准输出

  ```go
  // StdoutPipe方法返回一个在命令Start后与命令标准输出关联的管道。
  // Wait方法获知命令结束后会关闭这个管道，一般不需要显式的关闭该管道。
  // 但是在从管道读取完全部数据之前调用Wait是错误的
  // 同样使用StdoutPipe方法时调用Run函数也是错误的。
  func (c *Cmd) StdoutPipe() (io.ReadCloser, error)
  ```

  Example:

  ```go
  package main
  
  import (
  	"encoding/json"
  	"fmt"
  	"log"
  	"os/exec"
  )
  
  func main() {
  	// 执行命令: echo -n {"Name": "Bob", "Age": 32}
  	cmd := exec.Command("echo", "-n", `{"Name": "Bob", "Age": 32}`)
  	// stdout接收标准输出内容
  	stdout, err := cmd.StdoutPipe()
  	if err != nil {
  		log.Fatal(err)
  	}
  	// 运行命令
  	if err := cmd.Start(); err != nil {
  		log.Fatal(err)
  	}
  	var person struct {
  		Name string
  		Age  int
  	}
  	// 将标准输出的内容解析到person对象
  	if err := json.NewDecoder(stdout).Decode(&person); err != nil {
  		log.Fatal(err)
  	}
  	// 等待标准输出读取结束
  	if err := cmd.Wait(); err != nil {
  		log.Fatal(err)
  	}
  	// 输出解析到的数据
  	fmt.Printf("%s is %d years old\n", person.Name, person.Age)
  }
  ```

  Output:

  ```bash
  $ go run main.go
  Bob is 32 years old
  ```

  

- 错误输出

  ```go
  // StderrPipe方法返回一个在命令Start后与命令标准错误输出关联的管道。
  // Wait方法获知命令结束后会关闭这个管道，一般不需要显式的关闭该管道。
  // 但是在从管道读取完全部数据之前调用Wait是错误的。
  // 同样使用StderrPipe方法时调用Run函数也是错误的。请参照StdoutPipe的例子。
  func (c *Cmd) StderrPipe() (io.ReadCloser, error)
  ```

  

- 运行cmd对象包含的命令

  ```go
  // Run执行cmd包含的命令，并阻塞直到完成。
  // 如果命令成功执行，stdin、stdout、stderr的转交没有问题，并且返回状态码为0，方法的返回值为nil；
  // 如果命令没有执行或者执行失败，会返回*ExitError类型的错误；否则返回的error可能是表示I/O问题。
  func (c *Cmd) Run() error
  ```

  

- 开始执行cmd对象包含的命令

  ```go
  // Start开始执行c包含的命令，但并不会等待该命令完成即返回。
  // Wait方法会返回命令的返回状态码并在命令返回后释放相关的资源。
  func (c *Cmd) Start() error
  ```

  Example:

  ```go
  package main
  
  import (
  	"log"
  	"os/exec"
  )
  
  func main() {
  	// 运行命令: sleep 5
  	cmd := exec.Command("sleep", "5")
  	// 开始执行命令
  	err := cmd.Start()
  	if err != nil {
  		log.Fatal(err)
  	}
  	log.Printf("Waiting for command to finish...")
  	// 等待命令执行结束
  	err = cmd.Wait()
  	log.Printf("Command finished with error: %v", err)
  }
  ```

  Output:

  ```bash
  $ go run main.go
  2020/07/20 11:10:48 Waiting for command to finish...
  2020/07/20 11:10:53 Command finished with error: <nil>
  ```

  

- 等待命令运行结束

  ```go
  // Wait会阻塞直到该命令执行完成，该命令必须是被Start方法开始执行的。
  // 如果命令成功执行，stdin、stdout、stderr的转交没有问题，并且返回状态码为0，方法的返回值为nil；
  // 如果命令没有执行或者执行失败，会返回*ExitError类型的错误；
  // 否则返回的error可能是表示I/O问题。Wait方法会在命令返回后释放相关的资源。
  func (c *Cmd) Wait() error
  ```

  

-  执行命令并返回标准输出

  ```go
  // 执行命令并返回标准输出的切片。
  func (c *Cmd) Output() ([]byte, error)
  ```

  Example:

  ```go
  package main
  
  import (
  	"fmt"
  	"log"
  	"os/exec"
  )
  
  func main() {
  	// 运行命令: date
  	// 获取标准输出信息
  	out, err := exec.Command("date").Output()
  	if err != nil {
  		log.Fatal(err)
  	}
  	fmt.Printf("The date is %s\n", out)
  }
  ```

  Output:

  ```bash
  $ go run main.go
  The date is 2020年07月20日 11:15:39
  ```

  

- 执行命令并返回标准输出和错误输出合并的切片

  ```go
  // 执行命令并返回标准输出和错误输出合并的切片。
  func (c *Cmd) CombinedOutput() ([]byte, error)
  ```

  

#### 2. 输出错误信息

- ErrNotFound

  ```go
  // 如果路径搜索没有找到可执行文件时，就会返回本错误。
  var ErrNotFound = errors.New("executable file not found in $PATH")
  ```

  

- Error

  ```go
  // Error类型记录执行失败的程序名和失败的原因。
  type Error struct {
      Name string
      Err  error
  }
  // 输出错误信息
  func (e *Error) Error() string
  ```

  

- ExitError

  ```go
  // ExitError报告某个命令的一次未成功的返回。
  type ExitError struct {
      *os.ProcessState
  }
  // 输出错误信息
  func (e *ExitError) Error() string
  ```

  

#### 3. 其他方法

exec包提供了搜索可执行文件的方法——LookPath

```go
// 在环境变量PATH指定的目录中搜索可执行文件，如file中有斜杠，则只在当前目录搜索。
// 返回完整路径或者相对于当前目录的一个相对路径。
func LookPath(file string) (string, error)
```

Example:

```go
package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	// 查找 echo 命令
	path, err := exec.LookPath("echo")
	if err != nil {
		log.Fatal("installing fortune is in your future")
	}
	// 输出命令保存的路径
	fmt.Printf("fortune is available at %s\n", path)
}
```

Output:

```bash
$ go run main.go
fortune is available at C:\Program Files\Git\usr\bin\echo.exe
```