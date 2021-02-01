> os 包提供了操作系统函数的不依赖平台的接口，本文向大家介绍 os 包的使用方法，希望对你有帮助。

## 一、系统

#### 1.  获取系统环境变量

- 获取主机名

  ```go
  // Hostname返回内核提供的主机名。
  func Hostname() (name string, err error)
  ```

- 获取系统内存页尺寸

  ```go
  // Getpagesize返回底层的系统内存页的尺寸。
  func Getpagesize() int
  ```

- 获取系统环境变量

  ```go
  // Environ返回表示环境变量的格式为"key=value"的字符串的切片拷贝。
  func Environ() []string
  ```

- 获取指定的环境变量

  ```go
  // Getenv检索并返回名为key的环境变量的值。如果不存在该环境变量会返回空字符串。
  func Getenv(key string) string
  ```

- 设置系统环境变量

  ```go
  // Setenv设置名为key的环境变量。如果出错会返回该错误。
  func Setenv(key, value string) error
  ```

- 删除所有的环境变量

  ```go
  // Clearenv删除所有环境变量。
  func Clearenv()
  ```

  

#### 2. 进程管理

- 退出程序

  ```go
  // Exit让当前程序以给出的状态码code退出。
  // 一般来说，状态码0表示成功，非0表示出错。
  // 程序会立刻终止，defer的函数不会被执行。
  func Exit(code int)
  ```

- 获取Pid

  ```go
  // Getpid返回调用者所在进程的进程ID。 
  func Getpid() int
  // Getppid返回调用者所在进程的父进程的进程ID。
  func Getppid() int
  ```

  

- n

#### 3. 用户管理

- 获取用户ID

  ```go
  // Getuid返回调用者的用户ID。
  func Getuid() int
  // Geteuid返回调用者的有效用户ID。
  func Geteuid() int
  ```

- 获取用户组

  ```go
  // Getgid返回调用者的组ID。
  func Getgid() int
  // Getegid返回调用者的有效组ID。
  func Getegid() int
  // Getgroups返回调用者所属的所有用户组的组ID。
  func Getgroups() ([]int, error)
  ```

  

#### 4. 其他功能

- 变量替换

  ```go
  // Expand函数替换s中的${var}或$var为mapping(var)。
  // 例如，os.ExpandEnv(s)等价于os.Expand(s, os.Getenv)。
  func Expand(s string, mapping func(string) string) string
  ```

- 环境变量值的替换

  ```go
  // ExpandEnv函数替换s中的${var}或$var为名为var 的环境变量的值。
  // 引用未定义环境变量会被替换为空字符串。
  func ExpandEnv(s string) string
  ```

  

## 二、文件操作

#### 1. 文件的权限和模式

在 os 库中 FileMode 代表文件的模式和权限位。这些字位在所有的操作系统都有相同的含义，因此文件的信息可以在不同的操作系统之间安全的移植。不是所有的位都能用于所有的系统，唯一共有的是用于表示目录的ModeDir位。下面是源码中定义的 FileMode 不同值的含义：

```go
type FileMode uint32
const (
	// The single letters are the abbreviations
	// used by the String method's formatting.
	ModeDir        FileMode = 1 << (32 - 1 - iota) // d: is a directory
	ModeAppend                                     // a: append-only
	ModeExclusive                                  // l: exclusive use
	ModeTemporary                                  // T: temporary file; Plan 9 only
	ModeSymlink                                    // L: symbolic link
	ModeDevice                                     // D: device file
	ModeNamedPipe                                  // p: named pipe (FIFO)
	ModeSocket                                     // S: Unix domain socket
	ModeSetuid                                     // u: setuid
	ModeSetgid                                     // g: setgid
	ModeCharDevice                                 // c: Unix character device, when ModeDevice is set
	ModeSticky                                     // t: sticky
	ModeIrregular                                  // ?: non-regular file; nothing else is known about this file

	// Mask for the type bits. For regular files, none will be set.
	ModeType = ModeDir | ModeSymlink | ModeNamedPipe | ModeSocket | ModeDevice | ModeCharDevice | ModeIrregular

	ModePerm FileMode = 0777 // Unix permission bits
)
```

官方提供了获取文件权限和模式的方法：

- 判断是否为目录

  ```go
  func (m FileMode) IsDir() bool
  ```

- 是否是一个普通文件

  ```go
  func (m FileMode) IsRegular() bool
  ```

- 获取文件权限

  ```go
  func (m FileMode) Perm() FileMode
  ```

- 将 FileMode 转换为字符串

  ```go
  func (m FileMode) String() string
  ```

我们通常不会直接调用 FileMode 的方法来获取文件信息，os 标准库中有一个 FileInfo 接口

```go
type FileInfo interface {
    Name() string       // 文件的名字（不含扩展名）
    Size() int64        // 普通文件返回值表示其大小；其他文件的返回值含义各系统不同
    Mode() FileMode     // 文件的模式位
    ModTime() time.Time // 文件的修改时间
    IsDir() bool        // 等价于Mode().IsDir()
    Sys() interface{}   // 底层数据来源（可以返回nil）
}
```

我们可以通过 Stat 方法获取一个描述 name 指定的文件对象的 FileInfo

```go
// 如果指定的文件对象是一个符号链接，返回的FileInfo描述该符号链接指向的文件的信息，本函数会尝试跳转该链接。
// 如果出错，返回的错误值为*PathError类型。
func Stat(name string) (fi FileInfo, err error)
```

还有一个 Lstat 方法与 Stat 的作用类似，不同点在于对符号链接的处理：

```go
// 如果指定的文件对象是一个符号链接，返回的FileInfo描述该符号链接的信息，本函数不会试图跳转该链接。
// 如果出错，返回的错误值为*PathError类型。
func Lstat(name string) (fi FileInfo, err error)
```

#### 3. 文件CURD操作

首先我们看一下源码中定义的文件操作权限：

```go
const (
    O_RDONLY int = syscall.O_RDONLY // 只读模式打开文件
    O_WRONLY int = syscall.O_WRONLY // 只写模式打开文件
    O_RDWR   int = syscall.O_RDWR   // 读写模式打开文件
    O_APPEND int = syscall.O_APPEND // 写操作时将数据附加到文件尾部
    O_CREATE int = syscall.O_CREAT  // 如果不存在将创建一个新文件
    O_EXCL   int = syscall.O_EXCL   // 和O_CREATE配合使用，文件必须不存在
    O_SYNC   int = syscall.O_SYNC   // 打开文件用于同步I/O
    O_TRUNC  int = syscall.O_TRUNC  // 如果可能，打开时清空文件
)
```

上面的常量在打开文件的时候需要使用，在 os 标准库中 File 结构体保存了一个打开的文件对象，下面我们看一下相关的方法：

- 创建文件
- 打开文件
- 

#### 4. 其他功能

- 判断字符是否是路径分隔符

  ```go
  // IsPathSeparator返回字符c是否是一个路径分隔符。
  func IsPathSeparator(c uint8) bool
  ```

- 判断文件置否存在

  ```go
  // 返回一个布尔值说明该错误是否表示一个文件或目录已经存在。
  func IsExist(err error) bool
  ```

- 判断文件是否不存在

  ```go
  // 返回一个布尔值说明该错误是否表示一个文件或目录不存在。
  func IsNotExist(err error) bool
  ```

- 判断是否有权限操作文件

  ```go
  // 返回一个布尔值说明该错误是否表示因权限不足要求被拒绝。
  func IsPermission(err error) bool
  ```

- 获取当前工作路径的根路径

  ```go
  // Getwd返回一个对应当前工作目录的根路径。
  // 如果当前目录可以经过多条路径抵达（因为硬链接），Getwd会返回其中一个。
  func Getwd() (dir string, err error)
  ```

- cd 到新的目录

  ```go
  // Chdir将当前工作目录修改为dir指定的目录。
  func Chdir(dir string) error
  ```

- 修改文件权限

  ```go
  // Chmod修改name指定的文件对象的mode。
  // 如果name指定的文件是一个符号链接，它会修改该链接的目的地文件的mode。
  func Chmod(name string, mode FileMode) error
  ```

- 修改文件所属用户及用户组

  ```go
  // Chmod修改name指定的文件对象的用户id和组id。
  // 如果name指定的文件是一个符号链接，它会修改该链接的目的地文件的用户id和组id。
  func Chown(name string, uid, gid int) error
  // Lchmod修改name指定的文件对象的用户id和组id。
  // 如果name指定的文件是一个符号链接，它会修改该符号链接自身的用户id和组id。
  func Lchown(name string, uid, gid int) error
  ```

- 修改文件访问时间和修改时间

  ```go
  // Chtimes修改name指定的文件对象的访问时间和修改时间，类似Unix的utime()或utimes()函数。
  func Chtimes(name string, atime time.Time, mtime time.Time) error
  ```

- 创建目录

  ```go
  // Mkdir使用指定的权限和名称创建一个目录。如果出错，会返回*PathError底层类型的错误。
  func Mkdir(name string, perm FileMode) error
  // MkdirAll使用指定的权限和名称创建一个目录，包括任何必要的上级目录，并返回nil，否则返回错误。
  // 权限位perm会应用在每一个被本函数创建的目录上。
  // 如果path指定了一个已经存在的目录，MkdirAll不做任何操作并返回nil。
  func MkdirAll(path string, perm FileMode) error
  ```

- 修改文件名称

  ```go
  // Rename修改一个文件的名字，移动一个文件。可能会有一些个操作系统特定的限制。
  func Rename(oldpath, newpath string) error
  ```

- 修改文件大小

  ```go
  // Truncate修改name指定的文件的大小。
  // 如果该文件为一个符号链接，将修改链接指向的文件的大小。
  // 如果出错，会返回*PathError底层类型的错误。
  func Truncate(name string, size int64) error
  ```

- 删除文件

  ```go
  // Remove删除name指定的文件或目录。如果出错，会返回*PathError底层类型的错误。
  func Remove(name string) error
  // RemoveAll删除path指定的文件，或目录及它包含的任何下级对象。
  // 它会尝试删除所有东西，除非遇到错误并返回。
  // 如果path指定的对象不存在，RemoveAll会返回nil而不返回错误。
  func RemoveAll(path string) error
  ```

- 获取符号链接文件的路径

  ```go
  // Readlink获取name指定的符号链接文件指向的文件的路径。如果出错，会返回*PathError底层类型的错误。
  func Readlink(name string) (string, error)
  ```

- 创建符号连接文件

  ```go
  // Symlink创建一个名为newname指向oldname的符号链接。如果出错，会返回* LinkError底层类型的错误。
  func Symlink(oldname, newname string) error
  ```

- 创建硬链接

  ```go
  // Link创建一个名为newname指向oldname的硬链接。如果出错，会返回* LinkError底层类型的错误。
  func Link(oldname, newname string) error
  ```

- 判断两个文件是否是同一个文件

  ```go
  // SameFile返回fi1和fi2是否在描述同一个文件。
  // 例如，在Unix这表示二者底层结构的设备和索引节点是相同的；在其他系统中可能是根据路径名确定的。
  // SameFile应只使用本包Stat函数返回的FileInfo类型值为参数，其他情况下，它会返回假。
  func SameFile(fi1, fi2 FileInfo) bool
  ```

- 返回一个保管临时文件的默认目录

  ```go
  // TempDir返回一个用于保管临时文件的默认目录。
  func TempDir() string
  ```

  