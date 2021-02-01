> 在Go语言中我们常常会用到接口类型来编写万能程序，函数接收到参数后，我们需要分析出参数类型，这就需要用到类型反射了，这篇文章向大家介绍Go语言标准库 reflect 的使用，希望对你有帮助。

## 一、interface和反射

#### 1. Go语言中类型设计原则

学习反射前，我们先了解一下Golang关于类型设计的一些原则：

- 变量包括（type，value）两部分。
- type包括static type和concrete type，简单来说static type是在编码看得见的类型（int、string、float64...），concrete type是runtime系统看见的类型。
- 类型断言能否成功，取决于变量的concrete type，而不是static type。因此，一个reader变量如果它的concrete type也实现了write方法的话，他也可以被类型断言为writer。

#### 2. 反射

反射是建立在类型之上的，Golang的指定类型的变量的类型是静态的（也就是指定int、string这些的变量，它的type是static type），在创建变量的时候就已经确定，反射主要与Golang的interface类型相关（它的type是concrete type），只有interface类型才有反射一说。

在Golang的实现中，每个interface变量都有一个对应pair，pair中记录了实际变量的值和类型:

`(value, type)`

value是实际变量值，type是实际变量的类型。一个interface{}类型的变量包含了2个指针，一个指针指向值的类型【对应concrete type】，另外一个指针指向实际的值【对应value】。

## 二、 reflect

#### 1. 获取接口实际值的类型

```go
// TypeOf返回接口中保存的值的类型，TypeOf(nil)会返回nil。
func TypeOf(i interface{}) Type
```

我们可以通过TypeOf方法获取一个接口类型中保存的值的类型，下面举个例子：

```go
package main

import (
	"fmt"
	"reflect"
)

type Response struct {
	Code int
	Msg  string
}

func main() {
	var Resp Response
	typeOfResp := reflect.TypeOf(Resp)
	fmt.Println(fmt.Sprintf(" resp type is %s, kind is %s",
		typeOfResp, typeOfResp.Kind()))
}
```

Output:

```bash
$ go run main.go
 resp type is main.Response, kind is struct
```

在学习这里的时候，我比较疑惑 Type （类型）和 Kind（种类）的区别，如上面 Resp 的类型是 main.Response，而种类是 struct。Type 是指变量所属的类型，包括系统的原生类型（int、string等）和我们使用 type 关键字定义的类型，而 Kind 是指变量类型所属的品种，参考 reflect.Kind 中的定义，主要有以下类型：

```go
const (
	Invalid Kind = iota
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	Array
	Chan
	Func
	Interface
	Map
	Ptr
	Slice
	String
	Struct
	UnsafePointer
)
```



#### 2. 获取接口保存的值

```go
// ValueOf返回一个初始化为i接口保管的具体值的Value，ValueOf(nil)返回Value零值。
func ValueOf(i interface{}) Value
```

通过调用ValueOf我们可以获得一个Value对象，Value实现了获取接口具体值的方法：

```go
// 返回v持有的布尔值，如果v的Kind不是Bool会panic
func (v Value) Bool() bool
// 返回v持有的有符号整数（表示为int64），如果v的Kind不是Int、Int8、Int16、Int32、Int64会panic
func (v Value) Int() int64
// 返回v持有的无符号整数（表示为uint64），如v的Kind不是Uint、Uintptr、Uint8、Uint16、Uint32、Uint64会panic
func (v Value) Uint() uint64
// 返回v持有的浮点数（表示为float64），如果v的Kind不是Float32、Float64会panic
func (v Value) Float() float64
// 返回v持有的复数（表示为complex64），如果v的Kind不是Complex64、Complex128会panic
func (v Value) Complex() complex128
// 将v持有的值作为一个指针返回。如果v的Kind是Slice，返回值是指向切片第一个元素的指针。如果持有的切片为nil，返回值为0；如果持有的切片没有元素但不是nil，返回值不会是0。
func (v Value) Pointer() uintptr
// 返回v持有的[]byte类型值。如果v持有的值的类型不是[]byte会panic。
func (v Value) Bytes() []byte
// 返回v持有的值的字符串表示。
func (v Value) String() string
// 返回v[i:j]（v持有的切片的子切片的Value封装）；如果v的Kind不是Array、Slice或String会panic。如果v是一个不可寻址的数组，或者索引出界，也会panic
func (v Value) Slice(i, j int) Value
```

上面列出了常用的方法，其他方法请参考源码。

下面举个例子：

```go
package main

import (
	"fmt"
	"reflect"
)

func main() {
	data := make(map[int]interface{})
	data[1] = 1
	data[2] = "I am string type"
	for _, v := range data {
		typeOfv := reflect.TypeOf(v)
		switch typeOfv.Kind().String() {
		case "int":
			fmt.Println(fmt.Sprintf("v type is %s, kind is %s, value is %d",
				typeOfv, typeOfv.Kind(), reflect.ValueOf(v).Int()))
		case "string":
			fmt.Println(fmt.Sprintf("v type is %s, kind is %s, value is %s",
				typeOfv, typeOfv.Kind(), reflect.ValueOf(v).String()))
		}
	}
}
```

Output:

```bash
$ go run main.go
v type is int, kind is int, value is 1
v type is string, kind is string, value is I am string type
```



#### 3. 修改接口实际变量的值

在 Go 语言中，任何通过 reflect.ValueOf  获取的 Value 都无法直接进行设定变量值的，因为 reflect.ValueOf 方法处理的都是值类型，即使是 &args 也是处理指针的拷贝，获取到 Value 无法对原来的变量进行寻址，所以直接设置变量值会报错。

但是，我们可以通过 Elem 方法对 Value 进行解引用获取具有指向原变量的指针，因此是可寻址可设定变量值的。

因此要修改实际变量的值需要满足一个条件，就是该变量是可以寻址的。

我们可以通过下面的方法判断一个变量是否是可寻址或者可设置变量值的：

```go
// 判断 v 是否可寻址
func (v Value) CanAddr() bool
// 判断 v 是否可以设置变量值
func (v Value) CanSet() bool
```

下面举个例子：

```go
package main

import (
	"fmt"
	"reflect"
)

// 定义结构体 Response
type Response struct {
	Code    int    `json:"code" require:"true"`
	Message string `json:"message" require:"true"`
}

func main() {
	resp := Response{Code: 200, Message: "Success"}
	fmt.Println("Old: ", resp)
	var Iresp = &resp
	// 获取 Iresp 的 Value 对象
	v := reflect.ValueOf(Iresp).Elem()
	// v.NumField() 为结构体字段数量
	for i := 0; i < v.NumField(); i++ {
		switch v.Field(i).Kind() {
		case reflect.String:
			if v.Field(i).CanAddr() {
				v.Field(i).SetString("Update")
			}
		case reflect.Int:
			if v.Field(i).CanAddr() {
				v.Field(i).SetInt(210)
			}
		}
	}
	fmt.Println("Update: ", resp)
}
```

测试一下：

```bash
$ go run main.go
Old:  {200 Success}
Update:  {210 Update}
```

结构体变量被成功修改。



#### 4. 通过反射访问结构体成员的值

Type 接口提供了用于获取结构体域类型对象的方法，下面为三个常用的方法：

```go
// 获取一个结构体内部的字段数量
NumField() int
// 根据 index 获取结构体内的成员字段类型
Field(i int) StructField
// 根据字段名获取结构体内的成员字段类型对象
FieldByName(name string) (StructField, bool)
```

通过上面的三个方法，我们就可以获得结构体变量内所有成员字段的类型对象 reflect.StructField：

```go
type StructField struct {
	Name      string    // 字段成员的名称
	PkgPath   string    // PkgPath是限定小写（未报告）字段名的包路径。大写（导出）字段名为空。
	Type      Type      // 字段类型
	Tag       StructTag // 字段tag
	Offset    uintptr   // 字节偏移量
	Index     []int     // 成员字段的 index
	Anonymous bool      // 成员字段是否公开
}
```

举个例子：

```go
package main

import (
	"fmt"
	"reflect"
)

func main() {
	// 定义结构体 Response
	type Response struct {
		Code    int    `json:"code" require:"true"`
		Message string `json:"message" require:"true"`
	}
	var resp Response
	resp.Code = 200
	resp.Message = "Success"

	var Iresp = resp
	// 获取 Iresp 的 Value 对象
	v := reflect.ValueOf(Iresp)
	// 获取 Iresp 的 Type 对象
	t := reflect.TypeOf(Iresp)
	// t.NumField() 为结构体字段数量
	for i := 0; i < t.NumField(); i++ {
		fmt.Println(fmt.Sprintf("{ Name: %s, Type: %s, Value: %v Tags : Json: %s Require: %s }",
			// 打印结构体字段名称
			t.Field(i).Name,
			// 打印结构体字段类型
			t.Field(i).Type.String(),
			// 打印结构体字段对应的值
			v.Field(i).Interface(),
			// 打印结构体 tag=json 的字段值
			t.Field(i).Tag.Get("json"),
			// 打印结构体 tag=require 的字段值
			t.Field(i).Tag.Get("require"),
		))
	}
}
```

Output：

```bash
$ go run main.go
{ Name: Code, Type: int, Value: 200 Tags : Json: code Require: true }
{ Name: Message, Type: string, Value: Success Tags : Json: message Require: true }
```



#### 5. 接口类型如何进行方法调用

Type 提供了获取接口下方法的方法类型对象 Method，如下：

```go
// 根据 index 查找方法
Method(int) Method
// 根据方法名查找方法
MethodByName(string) (Method, bool)
// 获取类型中公开的方法数量
NumMethod() int
```

我们来看看 Method 中都保存的那些信息：

```go
type Method struct {
    // Name是方法名。PkgPath是非导出字段的包路径，对导出字段该字段为""。
    // 结合PkgPath和Name可以从方法集中指定一个方法。
    // 参见http://golang.org/ref/spec#Uniqueness_of_identifiers
    Name    string
    PkgPath string
    Type  Type  // 方法类型
    Func  Value // 方法的值
    Index int   // 用于Type.Method的索引
}
```

举个例子：

```go
package main

import (
	"fmt"
	"reflect"
)

// 定义结构体 Response
type Response struct {
	Code    int    `json:"code" require:"true"`
	Message string `json:"message" require:"true"`
}

func (r *Response) Call0(data string) string {
	return fmt.Sprintf("Call0 【Code: %v Message: %v Data: %s】", r.Code, r.Message, data)
}

func (r *Response) Call1(data string) string {
	return fmt.Sprintf("Call1【Code: %v Message: %v Data: %s】", r.Code, r.Message, data)
}

func main() {
	var resp Response
	resp.Code = 200
	resp.Message = "Success"
	var Iresp = &resp
	// 获取 Iresp 的 Value 对象
	v := reflect.ValueOf(Iresp)
	// v.NumField() 为结构体字段数量
	for i := 0; i < v.NumMethod(); i++ {
		// 通过Call调用方法，传入参数data="Hello World!"
		ret := v.Method(i).Call([]reflect.Value{reflect.ValueOf("Hello World!")})
		// 打印方法返回的参数
		fmt.Println(len(ret), ret[0].String())
	}
}
```

测试一下：

```bash
$ go run main.go
1 Call0 【Code: 200 Message: Success Data: Hello World!】
1 Call1【Code: 200 Message: Success Data: Hello World!】
```

可以看到我们成功调用了resp的两个方法。