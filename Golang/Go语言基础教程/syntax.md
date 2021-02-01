



**按照惯例，以一个Hello World 开始这篇文档：**

```go
//GO语言中保留C语言中的注释方法。
//每个文件必须先声明包，GO语言中以包为管理单位。
//每个工程中只能有一个main包，一个文件夹即为一个工程。
//每个程序必须包含main包，开头处标明即可。
package main
//GO语言中不需要写分号
import (   //通过import导入包，可以导入多个，分行写，用双引号括起来，导入的包必须要用。
	"fmt"
)
//定义函数
func main() {  //注意这个括号必须和函数名在同一行
	fmt.Println("Hello World!") //打印hello world
}
```

### 一、变量和常量

#### 1. 变量的命名

Go 语言中命名规范：

- 变量名由字母、下划线、数字组成。
- 变量名不能以数字开头。
- 变量名名不能是关键字
- 变量名区分大小写

#### 2. 变量基本类型

| 类型       | 描述                                                         |
| :--------- | :----------------------------------------------------------- |
| uint       | 32位或64位                                                   |
| uint8      | 无符号 8 位整型 (0 到 255)                                   |
| uint16     | 无符号 16 位整型 (0 到 65535)                                |
| uint32     | 无符号 32 位整型 (0 到 4294967295)                           |
| uint64     | 无符号 64 位整型 (0 到 18446744073709551615)                 |
| int        | 32位或64位                                                   |
| int8       | 有符号 8 位整型 (-128 到 127)                                |
| int16      | 有符号 16 位整型 (-32768 到 32767)                           |
| int32      | 有符号 32 位整型 (-2147483648 到 2147483647)                 |
| int64      | 有符号 64 位整型 (-9223372036854775808 到 9223372036854775807) |
| byte       | uint8的别名(type byte = uint8)                               |
| rune       | int32的别名(type rune = int32)，表示一个unicode码            |
| uintptr    | 无符号整型，用于存放一个指针是一种无符号的整数类型，没有指定具体的bit大小但是足以容纳指针。 uintptr类型只有在底层编程是才需要，特别是Go语言和C语言函数库或操作系统接口相交互的地方。 |
| float32    | IEEE-754 32位浮点型数                                        |
| float64    | IEEE-754 64位浮点型数                                        |
| complex64  | 32 位实数和虚数                                              |
| complex128 | 64 位实数和虚数                                              |

对于虚数我们可以通过内置方法`real(变量) `取实部，`imag(变量)` 取虚部。

在 Go 语言中通过单引号内可以包含字符，双引号中可以包含字符串。

#### 3. 变量的声明和初始化

```go
package main

import (
	"fmt"
)

func main() {
	var a int //定义变量a的类型为int
	a = 1     //给a赋值为1
	var x, y, z int
	fmt.Printf("x=%d, y=%d, z=%d\n", x, y, z)
	var b int = 2 //定义变量b的；类型为int并赋值为2
	var (         //在（）内可以同时定义多个变量
		c int     //定义c的类型为int
		d int = 4 //定义d的类型为int并初始化为4
		e     = 5 //定义变量e的值为5，系统会自动推导e的类型
	)
	fmt.Printf("a=%d, b=%d, c=%d, d=%d, e=%d\n", a, b, c, d, e)
	f := 12                         //定义变量f的值为12系统会自动推断f的类型
	fmt.Printf("f type is %T\n", f) //%T可以查看值的类型

	i, j := 1, 2 //这种方式也可以同时定义多个变量
	fmt.Printf("i=%d, j=%d\n", i, j)
	//GO语言中有一个匿名变量_
	/*
		如果一个函数调用后会返回多个值，但是你只需要其中的某个值
		匿名变量是个很好的选择
	*/
	_, i, j = a, b, c
	/*
		可以同时进行多个变量赋值，可以用于两个变量值的交换
	*/
	fmt.Printf("i=%d, j=%d\n", i, j)

	a1, a2 := 20, 30
	a1, a2 = a2, a1
	fmt.Printf("a1=%d, a2=%d\n", a1, a2)
}
```



#### 4. 常量类型

定义常量时只需要将变量定义中的var替换为const即可

```go
package main

import (
	"fmt"
)

func main() {
	const a int = 10
	const (
		c = 20
		d = 10.1
	)
	fmt.Printf("a=%d, c=%d, d=%d\n", a, c, d)

	/*
		iota常量自动生成器，用于给常量赋值
		iota遇到const，重置为0
		每使用一次iota的值自增1
	*/
	const (
		x = iota
		y = iota
		z = iota
	)
	fmt.Printf("x=%d, y=%d, z=%d\n", x, y, z)

	const (
		x1 = iota //可以只在第一个常量处赋值，其他常量的值会在上一个常量的基础上自增1
		y1
		z1
	)
	fmt.Printf("x1=%d, y1=%d, z1=%d\n", x1, y1, z1)

	const (
		x2     = iota
		y2, z2 = iota, iota //iota如果是同一行，值都一样
	)
	fmt.Printf("x2=%d, y2=%d, z2=%d\n", x2, y2, z2)
}
```

#### 5. iota

`iota`常量自动生成器，用于给常量赋值，iota遇到const，重置为0，每使用一次iota的值自增1。如下面的例子：

```go

```



### 二、运算符

#### 1. 算数运算符



\+ -  *  /  %  ++  --

#### 2. 关系运算符
==  !=  >=   <=  >  <
#### 3. 逻辑运算符
&&与 ||或 !非（取反）
非0就是真，0就是假

#### 4. 位运算符

& 位与
| 位或
^ 异或
<< 循环左移
\>> 循环右移

### 三、控制语句

#### 1. 条件判断（if）

```go
package main

import (
	"fmt"
)

func main() {
	a := 10
	if a < 20 {
		fmt.Println("a<20")
	}

	if b := 30; b > 20 { //if 支持一个初始化语句，初始化语句与判断条件以分号隔开
		fmt.Println("b>20")
	}

	if x := 5; x < 5 {
		fmt.Println("x<5")
	} else {
		fmt.Println("x>=5")
	}

	if c := 40; c < 10 {
		fmt.Println("c<10")
	} else if c < 20 {
		fmt.Println("c<20")
	} else if c < 50 {
		fmt.Println("c<50")
	}

}
```



#### 2. 循环语句（for）

```go
package main

import (
	"fmt"
)

func main() {
	for i := 1; i < 5; i++ { //go语言中for循环不需要加括号
		fmt.Println(i)
	}

	a := "abcd"
	for index, key := range a {   //go语言可以像python中enumrate，返回a中元素的索引和键值
		fmt.Printf("a[%d]=%c  ", index, key)
	}

	for _, key := range a {  //可以使用匿名变量
		fmt.Println(key)
	}

	for index := range a {
		fmt.Printf("a[%d]=%c   ", index, a[index])
	}
}
```



#### 3. switch

```go
package main

import (
	"fmt"
)

func main() {
	a := 10
	switch { //可以使用判断语句
	case a > 10:
		fmt.Println("a > 10")
	case a == 10:
		fmt.Println("a = 10")
	case a < 10:
		fmt.Println("a < 10")
	default:
		fmt.Println("error")
	}

	c := 20
	switch c {
	case 10:
		fmt.Println("c=", 10)
	case 20:
		fmt.Println("c=", 20)
	default:
		fmt.Println("error")
	}

	switch d := 30; d { //可以有一个初始化语句，用分号隔开
	case 10:
		fmt.Println("d=", 10)
	case 20:
		fmt.Println("d=", 20)
	case 30:
		fmt.Println("d=", 30)
	default:
		fmt.Println("error")
	}

	switch e := 90; e {
	case 100, 90, 80: //一行可以写多个
		fmt.Println("优秀")
	case 70, 60:
		fmt.Println("合格")
	default:
		fmt.Println("不及格")
	}
}
```



#### 4. goto

```go
package main

import (
	"fmt"
)

func main() {
	a := 1
Key:
	fmt.Printf("Key")
	a++
	fmt.Printf("%dgoto\n", a)
	if a > 3 {
		return
	}
	goto Key
}
```



### 四、fmt 包

#### 1. 格式化输出与占位符

**通用占位符：**

```go
v     值的默认格式。
%+v   添加字段名(如结构体)
%#v　 相应值的Go语法表示 
%T    相应值的类型的Go语法表示 
%%    字面上的百分号，并非值的占位符
```

对于 `%v` 来说不同类型对应的格式如下：

```go
bool:                    %t 
int, int8 etc.:          %d 
uint, uint8 etc.:        %d, %x if printed with %#v
float32, complex64, etc: %g
string:                  %s
chan:                    %p 
pointer:                 %p
```

**布尔值：**

```go
%t   true 或 false
```

**整数值：**

```go
%b     二进制表示 
%c     相应Unicode码点所表示的字符 
%d     十进制表示 
%o     八进制表示 
%q     单引号围绕的字符字面值，由Go语法安全地转义 
%x     十六进制表示，字母形式为小写 a-f 
%X     十六进制表示，字母形式为大写 A-F 
%U     Unicode格式：U+1234，等同于 "U+%04X"
```

**浮点数及复数：**

```go
%b     无小数部分的，指数为二的幂的科学计数法，与 strconv.FormatFloat中的 'b' 转换格式一致。例如 -123456p-78 
%e     科学计数法，例如 -1234.456e+78 
%E     科学计数法，例如 -1234.456E+78 
%f     有小数点而无指数，例如 123.456 
%g     根据情况选择 %e 或 %f 以产生更紧凑的（无末尾的0）输出 
%G     根据情况选择 %E 或 %f 以产生更紧凑的（无末尾的0）输出
```

**字符串和bytes的slice表示：**

```go
%s     字符串或切片的无解译字节 
%q     双引号围绕的字符串，由Go语法安全地转义 
%x     十六进制，小写字母，每字节两个字符 
%X     十六进制，大写字母，每字节两个字符
```

**指针：**

```go
%p     十六进制表示，前缀 0x
```

我们可以通过占位符完成格式化输出的功能，如下面的例子：

```go

```



#### 2. 获取命令行参数

```go
package main

import (
	"fmt"
)

func main() {
	var a, b, c int
	fmt.Println("输入a:")
	fmt.Scanf("%d", &a) //输入一和C语言输入一样
	fmt.Println("输入b和c:")
	fmt.Scan(&b, &c) //第二种方式
	fmt.Printf("a=%d, b=%d, c=%d", a, b, c)
}
```



#### 3. 标准输入与输出

#### 4. Stderr

### 五、基本数据结构

#### 1. 指针

#### 2. 数组与切片

#### 3. map

#### 4. 结构体

### 六、函数





![点击并拖拽以移动](data:image/gif;base64,R0lGODlhAQABAPABAP///wAAACH5BAEKAAAALAAAAAABAAEAAAICRAEAOw==)