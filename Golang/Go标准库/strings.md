> 在 GO 语言中，我们可以使用 strings 标准库对字符串进行一系列操作，strings 标准库在日常的编程过程中是十分常用的，这里向大家介绍strings标准库的使用，希望对你有帮助。

## 一、字符串切割

#### 1. 字符串前后端字符清理

```go
// 返回将s前后端所有cutset包含的utf-8码值都去掉的字符串。
func Trim(s string, cutset string) string
// 返回将s前后端所有空白（unicode.IsSpace指定）都去掉的字符串。
func TrimSpace(s string) string
// 返回将s前端所有cutset包含的utf-8码值都去掉的字符串。
func TrimLeft(s string, cutset string) string
// 返回将s后端所有cutset包含的utf-8码值都去掉的字符串。
func TrimRight(s string, cutset string) string
// 返回去除s可能的前缀prefix的字符串。
func TrimPrefix(s, prefix string) string
// 返回去除s可能的后缀suffix的字符串。
func TrimSuffix(s, suffix string) string
```

#### 2. 分割字符

```go
// 返回将字符串按照空白（unicode.IsSpace确定，可以是一到多个连续的空白字符）分割的多个字符串。
// 如果字符串全部是空白或者是空字符串的话，会返回空切片。
func Fields(s string) []string
// 用去掉s中出现的sep的方式进行分割，会分割到结尾，并返回生成的所有片段组成的切片（每一个sep都会进行一次切割，即使两个sep相邻，也会进行两次切割）。
// 如果sep为空字符，Split会将s切分成每一个unicode码值一个字符串。
func Split(s, sep string) []string
```



## 二、字符串查找

#### 1. 子串sep在字符串s中第一次出现的位置，不存在则返回-1

```go
func Index(s, sep string) int
```

#### 2. 字符c在s中第一次出现的位置，不存在则返回-1

```go
func IndexByte(s string, c byte) int
```

#### 3. unicode码值r在s中第一次出现的位置，不存在则返回-1

```go
func IndexRune(s string, r rune) int
```

#### 4. 字符串chars中的任一utf-8码值在s中第一次出现的位置，如果不存在或者chars为空字符串则返回-1

```go
func IndexAny(s, chars string) int
```

#### 5. s中第一个满足函数f的位置i（该处的utf-8码值r满足f(r)==true），不存在则返回-1

```go
func IndexFunc(s string, f func(rune) bool) int
```

#### 6. 子串sep在字符串s中最后一次出现的位置，不存在则返回-1

```go
func LastIndex(s, sep string) int
```

#### 7. 字符串chars中的任一utf-8码值在s中最后一次出现的位置，如不存在或者chars为空字符串则返回-1

```go
func LastIndexAny(s, chars string) int
```

#### 8. s中最后一个满足函数f的unicode码值的位置i，不存在则返回-1

```go
func LastIndexFunc(s string, f func(rune) bool) int
```

## 三、字符串转换

#### 1. 大小写转换

```go
// 返回s中每个单词的首字母都改为标题格式的字符串拷贝，相当于首字母大写
func Title(s string) string
// 返回将所有字母都转为对应的小写版本的拷贝。
func ToLower(s string) string
// 返回将所有字母都转为对应的大写版本的拷贝。
func ToUpper(s string) string
// 返回将所有字母都转为对应的标题版本的拷贝，会将所有字母大写
func ToTitle(s string) string
```



#### 2. 字符串替换

```go
// 返回将s中前n个不重叠old子串都替换为new的新字符串，如果n<0会替换所有old子串。
func Replace(s, old, new string, n int) string
// 将s的每一个unicode码值r都替换为mapping(r)，返回这些新码值组成的字符串拷贝。
// 如果mapping返回一个负值，将会丢弃该码值而不会被替换。（返回值中对应位置将没有码值）
func Map(mapping func(rune) rune, s string) string
```



## 四、 字符串连接

#### 1. 返回count个s串联的字符串

```go
func Repeat(s string, count int) string
```

#### 2. 将一系列字符串连接为一个字符串，之间用sep来分隔

```go
func Join(a []string, sep string) string
```



## 五、其他功能

#### 1. 判断字符串s是否包含utf-8码值r

```go
func ContainsRune(s string, r rune) bool
```

#### 2. 判断字符串s是否包含字符串chars中的任一字符

```go
func ContainsAny(s, chars string) bool
```

#### 3. 返回字符串s中有几个不重复的sep子串

```go
func Count(s, sep string) int
```



## 六、 常见案例

#### 1. 删除字符串末尾的换行符

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	str := `Hello World

`
	fmt.Println(str)
	// 方法一
	fmt.Println(strings.TrimRight(str, "\n"))
	// 方法二
	fmt.Println(strings.TrimSpace(str))
	// 方法三
	fmt.Println(strings.Trim(str, "\n"))
	// 方法四
	fmt.Println(strings.Replace(str, "\n", "", -1))
}
```



#### 2. 将字符串转换为大写

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	str := `Hello World`
	// 方法一
	fmt.Println(strings.ToTitle(str))
	// 方法二
	fmt.Println(strings.ToUpper(str))
}
```



#### 3. 以冒号为分隔符，切割字符串

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	str := `root:x:0:0:root:/root:/bin/bash`
	// 方法一
	fmt.Println(strings.Split(str, ":"))
}
```



#### 4.  替换字符串

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	str := `    Hello World randow_w1    `
	// 方法一
	fmt.Println(strings.Replace(str, "randow_w1", "everyone", -1))
}
```



#### 5.  去除字符串两边的空格

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	str := `    Hello World    `
	fmt.Println(str)
	// 方法一
	fmt.Println(strings.Trim(str, " "))
	// 方法二
	fmt.Println(strings.TrimSpace(str))
	// 方法三
	fmt.Println(strings.TrimRight(str, " ")) // 去除右边空格
	// 方法四
	fmt.Println(strings.TrimLeft(str, " ")) // 去除左边空格
}
```



#### 6. 将切片以空格为连接符组成一个字符串

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	books := []string{"Math", "Chinese", "Science", "English"}
	// 方法一
	fmt.Println(strings.Join(books, " "))
}
```





