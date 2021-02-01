> gob包是用来管理gob流的，它可以实现在编码器（发送器）和解码器（接收器）之间进行二进制数据流的发送，一般用来传递远端程序调用的参数和结果，比如net/rpc包就有用到这个。下面我们来学习以下gob标准库的使用，希望对你有帮助。

## 一、主要函数介绍

> gob和json的pack之类的方法一样，由发送端使用Encoder对数据结构进行编码。在接收端收到消息之后，接收端使用Decoder将序列化的数据变化成本地变量。与json编码格式相比，gob编码可以实现json所不支持的struct的方法序列化，利用gob包序列化struct保存到本地也十分方便。

#### 1. 记录数据的类型和名称

- Register

  ```go
  // Register记录value下层具体值的类型和其名称。
  // 该名称将用来识别发送或接受接口类型值时下层的具体类型。
  // 本函数只应在初始化时调用，如果类型和名字的映射不是一一对应的，会panic。
  func Register(value interface{})
  ```

- RegisterName

  ```go
  // RegisterName类似Register，区别是这里使用提供的name代替类型的默认名称。
  func RegisterName(name string, value interface{})
  ```

#### 2. 编码

数据在传输时会先经过编码后再进行传输，与编码相关的有三种方法：

- NewDecoder

  ```go
  // NewEncoder返回一个将编码后数据写入w的*Encoder。
  func NewEncoder(w io.Writer) *Encoder
  ```

- Encode

  ```go
  // Encode方法将e编码后发送，并且会保证所有的类型信息都先发送。
  func (enc *Encoder) Encode(e interface{}) error
  ```

- EncodeValue

  ```go
  // EncodeValue方法将value代表的数据编码后发送，并且会保证所有的类型信息都先发送。
  func (enc *Encoder) EncodeValue(value reflect.Value) error
  ```

#### 3. 解码

接收到数据后需要对数据进行解码，与解码相关的有三种方法：

- NewDecoder

  ```go
  // 函数返回一个从r读取数据的*Decoder，如果r不满足io.ByteReader接口，则会包装r为bufio.Reader。
  func NewDecoder(r io.Reader) *Decoder
  ```

- Decode

  ```go
  // Decode从输入流读取下一个之并将该值存入e。
  // 如果e是nil，将丢弃该值,否则e必须是可接收该值的类型的指针。
  // 如果输入结束，方法会返回io.EOF并且不修改e（指向的值）。
  func (dec *Decoder) Decode(e interface{}) error
  ```

- DecodeValue

  ```go
  // DecodeValue从输入流读取下一个值，如果v是reflect.Value类型的零值（v.Kind() == Invalid），方法丢弃该值；否则它会把该值存入v。
  // 此时，v必须代表一个非nil的指向实际存在值的指针或者可写入的reflect.Value（v.CanSet()为真）。
  // 如果输入结束，方法会返回io.EOF并且不修改e（指向的值）。
  func (dec *Decoder) DecodeValue(v reflect.Value) error
  ```

## 二、gob编码的优势和局限性

#### 1. 局限性

这里需要明确一点，gob只能用在golang中，所以在实际工程开发过程中，如果与其他端，或者其他语言打交道，那么gob是不可以的，我们就要使用json了。

#### 2. 优势

>Gob流是自解码的。流中的所有数据都有前缀（采用一个预定义类型的集合）指明其类型。指针不会传递，而是传递值；也就是说数据是压平了的。递归的类型可以很好的工作，但是递归的值（比如说值内某个成员直接/间接指向该值）会出问题。这个问题将来可能会修复。

要使用gob，先要创建一个编码器，并向其一共一系列数据：可以是值，也可以是指向实际存在数据的指针。编码器会确保所有必要的类型信息都被发送。在接收端，解码器从编码数据流中恢复数据并将它们填写进本地变量里。

gob编解码并没有要求发送方和接收方的结构完全一致，下面是官方文档的翻译：

```go
// 定义一个结构体
	struct { A, B int }

// 下面类型的数据都是可以发送和接收的:

	struct { A, B int }	// the same
	*struct { A, B int }	// extra indirection of the struct
	struct { *A, **B int }	// extra indirection of the fields
	struct { A, B int64 }	// different concrete value type; see below

// 下面类型也可以接收:

	struct { A, B int }	// the same
	struct { B, A int }	// ordering doesn't matter; matching is by name
	struct { A, B, C int }	// extra field (C) ignored
	struct { B int }	// missing field (A) ignored; data will be dropped
	struct { B, C int }	// missing field (A) ignored; extra field (C) ignored.

// 下面格式是有问题的:

	struct { A int; B uint }	// change of signedness for B
	struct { A int; B float }	// change of type for B
	struct { }			// no field names in common
	struct { C, D int }		// no field names in common
```



## 三、使用gob编程

这里我们看一下官网的例子https://golang.org/pkg/encoding/gob/：

#### 1. Basic

```go
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

type P struct {
	X, Y, Z int
	Name    string
}
type Q struct {
	X, Y *int32
	Name string
}

// 这是一个基础的使用用例，创建一个编码器，对数据进行编码，然后使用解码器接收数据
func main() {
	// 初始化编码器，创建一个decoder实例
	var network bytes.Buffer        // 标准输入
	enc := gob.NewEncoder(&network) // 编码
	dec := gob.NewDecoder(&network) // 解码
	// 编码器发送数据
	err := enc.Encode(P{3, 4, 5, "Pythagoras"})
	if err != nil {
		log.Fatal("encode error:", err)
	}
	err = enc.Encode(P{1782, 1841, 1922, "Treehouse"})
	if err != nil {
		log.Fatal("encode error:", err)
	}
	// 解码器接收数据
	var q Q
	err = dec.Decode(&q)
	if err != nil {
		log.Fatal("decode error 1:", err)
	}
	fmt.Printf("%q: {%d, %d}\n", q.Name, *q.X, *q.Y)
	err = dec.Decode(&q)
	if err != nil {
		log.Fatal("decode error 2:", err)
	}
	fmt.Printf("%q: {%d, %d}\n", q.Name, *q.X, *q.Y)
}
```

测试一下：

```bash
$ go run main.go
"Pythagoras": {3, 4}
"Treehouse": {1782, 1841}
```



#### 2. 自定义Encode和Decode

```go
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

// Vector 类型实现了BinaryMarshal/BinaryUnmarshal的方法，这样我们就可以发送和接受gob类型的数据。
// 我们可以等效地使用本地定义的gobcodencode/gobcodector接口
type Vector struct {
	x, y, z int
}

func (v Vector) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, _ = fmt.Fprintln(&b, v.x, v.y, v.z)
	return b.Bytes(), nil
}

// UnmarshalBinary 修改接收器，所以必须要传递指针类型
func (v *Vector) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &v.x, &v.y, &v.z)
	return err
}

// 此示例传输实现自定义编码和解码方法的值。
func main() {
	var network bytes.Buffer // 定义标准输入
	// 创建一个编码器发送数据
	enc := gob.NewEncoder(&network)
	err := enc.Encode(Vector{3, 4, 5})
	if err != nil {
		log.Fatal("encode:", err)
	}
	// 创建一个解码器接收数据
	dec := gob.NewDecoder(&network)
	var v Vector
	err = dec.Decode(&v)
	if err != nil {
		log.Fatal("decode:", err)
	}
	fmt.Println(v)
}
```

测试结果：

```bash
$ go run main.go
{3 4 5}
```



#### 3.编码interface

```go
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"math"
)

type Point struct {
	X, Y int
}

func (p Point) Hypotenuse() float64 {
	return math.Hypot(float64(p.X), float64(p.Y))
}

type Pythagoras interface {
	Hypotenuse() float64
}

// 这里展示如何对一个接口类型的值进行编码
// key与常规类型的区别是注册实现接口的具体类型。
func main() {
	var network bytes.Buffer // 标准输入
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(Point{})
	// 创建一个encoder接口，并发送值
	enc := gob.NewEncoder(&network)
	for i := 1; i <= 3; i++ {
		interfaceEncode(enc, Point{3 * i, 4 * i})
	}
	// 创建一个decoder接口，并接收值
	dec := gob.NewDecoder(&network)
	for i := 1; i <= 3; i++ {
		result := interfaceDecode(dec)
		fmt.Println(result.Hypotenuse())
	}
}

// interfaceEncode 将值编码并保存到encoder.
func interfaceEncode(enc *gob.Encoder, p Pythagoras) {
	// The encode will fail unless the concrete type has been
	// registered. We registered it in the calling function.
	// Pass pointer to interface so Encode sees (and hence sends) a value of
	// interface type.  If we passed p directly it would see the concrete type instead.
	// See the blog post, "The Laws of Reflection" for background.
	err := enc.Encode(&p)
	if err != nil {
		log.Fatal("encode:", err)
	}
}
// interfaceDecode 解码接口的值并返回
func interfaceDecode(dec *gob.Decoder) Pythagoras {
	// The decode will fail unless the concrete type on the wire has been
	// registered. We registered it in the calling function.
	var p Pythagoras
	err := dec.Decode(&p)
	if err != nil {
		log.Fatal("decode:", err)
	}
	return p
}
```

测试结果：

```go
$ go run main.go
5
10
15
```



