> unsafe包提供了一些跳过go语言类型安全限制的操作。这个标准库用的比较少，这篇博客向大家介绍以下，如何使用unsafe包，希望对你有帮助。

#### 1. unsafe包介绍 

##### （1）两种重要类型

- ArbitraryType

  ```go
  // ArbitraryType表示任意一种类型，但并非一个实际存在与unsafe包的类型。
  type ArbitraryType int
  ```

- Pointer

  ```go
  // Pointer类型用于表示任意类型的指针。
  type Pointer *ArbitraryType
  ```

  有4个特殊的只能用于Pointer类型的操作：

  - 任意类型的指针可以转换为一个Pointer类型值。
  - 一个Pointer类型值可以转换为任意类型的指针。
  - 一个uintptr类型值可以转换为一个Pointer类型值。
  - 一个Pointer类型值可以转换为一个uintptr类型值。

  因此，Pointer类型允许程序绕过类型系统读写任意内存。使用它时必须谨慎。

（2）三个方法

- Sizeof

  ```go
  // Sizeof返回类型v本身数据所占用的字节数。
  // 返回值是“顶层”的数据占有的字节数。
  // 例如，若v是一个切片，它会返回该切片描述符的大小，而非该切片底层引用的内存的大小。
  func Sizeof(v ArbitraryType) uintptr
  ```

- Alignof

  ```go
  // Alignof返回类型v的对齐方式（即类型v在内存中占用的字节数）
  // 若是结构体类型的字段的形式，它会返回字段f在该结构体中的对齐方式。
  func Alignof(v ArbitraryType) uintptr
  ```

- Offsetof

  ```go
  // Offsetof返回类型v所代表的结构体字段在结构体中的偏移量，它必须为结构体类型的字段的形式。
  // 换句话说，它返回该结构起始处与该字段起始处之间的字节数。
  func Offsetof(v ArbitraryType) uintptr
  ```

#### 2. unsafe包使用

##### （1）指针类型转换

```go
package main

import (
	"fmt"
	"unsafe"
)

func main() {
	// 定义两个不同类型的变量u和i
	u := uint32(32)
	i := int32(1)
	// 输出地址
	fmt.Println(&u, &i)
	p := &i
	// 修改p的指针类型
	p = (*int32)(unsafe.Pointer(&u))
	fmt.Println(p)
	fmt.Println(*p)
}
```

测试结果：

```bash
$ go run main.go
0xc000070058 0xc00007005c
0xc000070058
32
```



##### （2）更复杂的示例

在看代码前，我们先来了解一下结构体的一些基本概念：

- 结构体的成员变量在内存存储上是一段连续的内存
- 结构体的初始地址就是第一个成员变量的内存地址
- 基于结构体的成员地址去计算偏移量。就能够得出其他成员变量的内存地址

```go
package main

import (
	"fmt"
	"unsafe"
)

type Student struct {
	Name string
	Age  int64
}

func main() {
	s := Student{Name: "random_w", Age: 18}
	sPointer := unsafe.Pointer(&s)
	// 修改Name变量的值，因为结构体存储在内存上的一段连续内存
	// 因此我们直接取出指针后转换为 Pointer，再强制转换为字符串类型的指针值就可以获取第一个变量Name的地址
	NamePointer := (*string)(unsafe.Pointer(sPointer))
	// 修改该地址的值
	*NamePointer = "random_w_update"
	// 通过给sPointer指针的地址加上偏移量就可以获取到第二个变量的地址。
    // unsafe.Offsetof会获取回类型v所代表的结构体字段在结构体中的偏移量
    // uintptr 是 Go 的内置类型。返回无符号整数，可存储一个完整的地址，后续常用于指针运算。
	AgePointer := (*int64)(unsafe.Pointer(uintptr(sPointer) + unsafe.Offsetof(s.Age)))
	*AgePointer = 19

	fmt.Printf("s.Name: %s, s.Age: %d", s.Name, s.Age)
}
```

测试结果：

```bash
$ go run main.go
s.Name: random_w_update, s.Age: 19
```



#### 总结

-  `unsafe.Pointer` 可以让你的变量在不同的指针类型转来转去，也就是表示为任意可寻址的指针类型。
- `uintptr` 常用于与 `unsafe.Pointer` 打配合，用于做指针运算。

<font color=red>注意：没有特殊必要的话。是不建议使用 `unsafe` 标准库，它并不安全</font>