> 在 Go 语言中我们可以直接使用 container 标准库完成链表和堆操作，非常方便，我们不需要自己去实现这些方法，本文向大家介绍 container 库的使用方法，希望对你有帮助。

## 一、双向链表

#### 1. Element

Element 中保存了链表的所有信息，我们来看一下源码：

```go
type Element struct {
	// 双链接元素列表中的下一个和上一个指针。
	next, prev *Element
	// 此元素所属的链表。
	list *List
	// value 中保存了链表的值
	Value interface{}
}

// List是一个双向链表
type List struct {
	root Element // 当前节点元素
	len  int     // 链表长度
}
```

#### 2.  生成一个链表

生成一个链表很简单，我们通过下面的函数可以获得一个 List 对象的指针：

`func New() *List`

然后我们对链表进行初始化：

`func (l *List) Init() *List`

Init会清空 List 链表中的数据，生成一个全新的链表，因此我们也可以`使用 Init 方法来清空一个链表`。 

#### 3. 链表的增删改查

- 插入链表

  ```go
  // 将 v 插入链表的第一个位置
  func (l *List) PushFront(v interface{}) *Element
  // 复制 other 链表，并将 other 的最后一个元素连接到 l 的第一个位置
  func (l *List) PushFrontList(other *List)
  // 将 v 插入链表的最后一个元素
  func (l *List) PushBack(v interface{}) *Element
  // 复制 other 链表，并将 other 的第一个元素连接到 l 的最后一个位置
  func (l *List) PushBackList(other *List)
  // InsertBefore将一个值为v的新元素插入到mark前面，并返回生成的新元素。
  // 如果mark不是l的元素，l不会被修改。
  func (l *List) InsertBefore(v interface{}, mark *Element) *Element
  // InsertAfter将一个值为v的新元素插入到mark后面，并返回新生成的元素。
  // 如果mark不是l的元素，l不会被修改。
  func (l *List) InsertAfter(v interface{}, mark *Element) *Element
  ```

- 查询链表

  ```go
  // Next返回链表的后一个元素。
  func (e *Element) Next() *Element
  // Prev返回链表的前一个元素
  func (e *Element) Prev() *Element
  // 查询链表第一个元素
  func (l *List) Front() *Element
  // 查询链表最后一个元素
  func (l *List) Back() *Element
  ```

- 删除链表

  ```go
  // Remove删除链表中的元素e，并返回e.Value。
  func (l *List) Remove(e *Element) interface{}
  ```

- 修改链表

  ```go
  // MoveToFront将元素e移动到链表的第一个位置，如果e不是l的元素，l不会被修改。
  func (l *List) MoveToFront(e *Element)
  // MoveToBack将元素e移动到链表的最后一个位置，如果e不是l的元素，l不会被修改。
  func (l *List) MoveToBack(e *Element)
  // MoveBefore将元素e移动到mark的前面。
  // 如果e或mark不是l的元素，或者e==mark，l不会被修改。
  func (l *List) MoveBefore(e, mark *Element)
  // MoveAfter将元素e移动到mark的后面。如果e或mark不是l的元素，或者e==mark，l不会被修改。
  func (l *List) MoveAfter(e, mark *Element)
  ```

- 其他功能

  ```go
  // 查询链表的长度
  func (l *List) Len() int
  ```

#### 4.  我们实现一个反向链表的功能

#### 题目描述：

反转一个单链表。

**示例:**

```
输入: 1->2->3->4->5->NULL
输出: 5->4->3->2->1->NULL
```

```go
func reverseList(l *list.List) *list.List {
	// 获取 l 的头节点
	head := l.Front()
	for head.Next() != nil {
		// head 移动到最后一个位置结束
		if head.Next() == nil {
			break
		}
		// 将 head.Next 移动到第一个位置
		l.MoveToFront(head.Next())
	}
	return l
}
```

## 二、环形链表

#### 1. Ring

```go
// Ring是一个环形链表，他没有开始和结尾
type Ring struct {
	next, prev *Ring
	Value      interface{} // 用户可以自定义
}
```

#### 2. 环形链表的初始化

我们可以通过 New 函数初始化一个有 n 个节点的环形链表

`func New(n int) *Ring `

#### 3. 环形链表管理

- 查看链表长度

  ```go
  func (r *Ring) Len() int
  ```

- 返回链表下一个元素

  ```go
  func (r *Ring) Next() *Ring
  ```

- 返回链表的上一个元素

  ```go
  func (r *Ring) Prev() *Ring
  ```

- 返回指定位置的元素

  ```go
  // 返回移动n个位置（n>=0向前移动，n<0向后移动）后的元素，r不能为空。
  func (r *Ring) Move(n int) *Ring
  ```

- 连接两个环形链表

  ```go
  // Link连接r和s，并返回r原本的后继元素r.Next()。r不能为空。
  func (r *Ring) Link(s *Ring) *Ring
  ```

- 删除环形链表中的元素

  ```go
  // 删除链表中n % r.Len()个元素，从r.Next()开始删除。
  // 如果n % r.Len() == 0，不修改r。返回删除的元素构成的链表，r不能为空。
  func (r *Ring) Unlink(n int) *Ring
  ```

- 管理链表中的元素

  ```go
  // 对链表的每一个元素都执行f（正向顺序）
  func (r *Ring) Do(f func(interface{}))
  ```

#### 4. 举个例子

计算环形链表上所有值的和：

```go
package main

import (
	"container/ring"
	"fmt"
)

type SumInt struct {
	Value int
}

func (s *SumInt) add(i interface{}) {
	s.Value += i.(int)
}

func main() {
	r := ring.New(10)

	for i := 0; i < 10; i++ {
		r.Value = i
		r = r.Next()
	}

	sum := SumInt{}
	r.Do(sum.add)
	fmt.Println(sum.Value)
}
```



## 三、堆

> 堆(Heap)是计算机科学中一类特殊的数据结构的统称。堆通常是一个可以被看做一棵完全二叉树的数组对象。

#### 1.  构建堆需要满足的条件

container/heap 标准库中定义了对（最小）堆进行操作的接口：

```go
type Interface interface {
	sort.Interface
	Push(x interface{}) // 向末尾添加元素
	Pop() interface{}   // 从末尾删除元素
}
// sort.Interface
type Interface interface {
	// 元素个数
	Len() int
	// 索引为 i 的元素是否需要排在 索引为 j 的元素前面
    // Less方法的实现决定了堆是最大堆还是最小堆。
	Less(i, j int) bool
	// 交换索引为 i 和 j 的元素
	Swap(i, j int)
}
```

> 任何实现了本接口的类型都可以用于构建最小堆。最小堆，是一种经过排序的完全二叉树，其中任一非终端节点的数据值均不大于其左子节点和右子节点的值。

最小堆可以通过heap.Init建立，数据是递增顺序或者空的话也是最小堆。最小堆的约束条件是：

`!h.Less(j, i) for 0 <= i < h.Len() and 2*i+1 <= j <= 2*i+2 and j < h.Len()`

#### 2. 最小堆管理

- 插入元素

  ```go
  // 向堆h中插入元素x，并保持堆的约束性。复杂度O(log(n))，其中n等于h.Len()。
  func Push(h Interface, x interface{})
  ```

- 删除堆中最小元素

  ```go
  // 删除并返回堆h中的最小元素（不影响约束性）。
  // 复杂度O(log(n))，其中n等于h.Len()。等价于Remove(h, 0)。
  func Pop(h Interface) interface{}
  ```

- 删除堆中第一个元素

  ```go
  // 删除堆中的第i个元素，并保持堆的约束性。复杂度O(log(n))，其中n等于h.Len()。
  func Remove(h Interface, i int) interface{}
  ```

- 修改指定的元素

  ```go
  // 在修改第i个元素后，调用本函数修复堆，比删除第i个元素后插入新元素更有效率。
  // 复杂度O(log(n))，其中n等于h.Len()。
  func Fix(h Interface, i int)
  ```

#### 3. 通过最小堆实现排序

```go
package main

import (
	"container/heap"
	"fmt"
)

type myHeap []int

/* 实现排序 */
func (h *myHeap) Len() int {
	return len(*h)
}

func (h *myHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

// 最小堆实现
func (h *myHeap) Less(i, j int) bool {
	return (*h)[i] < (*h)[j]
}

/* 实现往堆中添加元素 */
func (h *myHeap) Push(v interface{}) {
	*h = append(*h, v.(int))
}

/* 实现删除堆中元素 */
func (h *myHeap) Pop() (v interface{}) {
	*h, v = (*h)[:len(*h)-1], (*h)[len(*h)-1]
	return
}

// 按层来遍历和打印堆数据，第一行只有一个元素，即堆顶元素
func (h myHeap) printHeap() {
	n := 1
	levelCount := 1
	for n <= h.Len() {
		fmt.Println(h[n-1 : n-1+levelCount])
		n += levelCount
		levelCount *= 2
	}
}

func main() {
	data := [7]int{13, 12, 45, 23, 11, 9, 20}
	aHeap := new(myHeap)
	for i := 0; i < len(data); i++ {
		aHeap.Push(data[i])
	}
	aHeap.printHeap()

	// 堆排序处理
	heap.Init(aHeap)
	fmt.Println("排序后: ")
	aHeap.printHeap()
}
```

## 四、单链表

官方的标准库中并没有实现单链表，这里我们来写一个单链表。

#### 1.  定义 List 结构体

```go
type Node struct {
	Value interface{} // 节点的值
	Next  *Node       // 指向下一个节点

}

type List struct {
	Head   *Node // 头节点
	Length int   // 链表长度
}
```

#### 2. 实现链表的增删改查

```go
func (l *List) Init() *List {
	l.Head = new(Node)
	l.Length = 0
	return l
}

func New() *List {
	return new(List).Init()
}

// Len 计算链表长度
func (l *List) Len() int {
	return l.Length
}

// Insert 链表中插入元素
func (l *List) Insert(i int, value interface{}) {
	// 链表长度小于 i - 1 直接退出
	if l.Len() < i-1 {
		return
	}
	head := l.Head
	for ; i > 1; i-- {
		head = head.Next
	}
	Next := new(Node)
	l.Length++
	Next.Value = value
	// 在末尾添加
	if head.Next == nil {
		head.Next = Next
		return
	}
	// 在链表中添加
	Next.Next = head.Next
	head.Next = Next
	return
}

// Delete 删除元素
func (l *List) Delete(i int) {
	// 链表长度小于 i 直接退出
	if l.Len() < i {
		return
	}
	head := l.Head
	for ; i > 1; i-- {
		head = head.Next
	}
	l.Length--
	if head.Next == nil {
		head.Next = nil
		return
	}
	head.Next = head.Next.Next
	return
}

// Get 获取第i个节点的值
func (l *List) Get(i int) interface{} {
	// 链表长度小于 i 直接退出
	if l.Len() < i {
		return nil
	}
	head := l.Head
	for ; i > 1; i-- {
		head = head.Next
	}
	return head.Next.Value
}

// Pop 删除链表末尾的值
func (l *List) Pop() {
	head := l.Head
	next := head.Next
	for next.Next != nil {
		head = head.Next
		next = head.Next
	}
	head.Next = nil
	l.Length--
	return
}

// 在链表末尾添加值
func (l *List) Push(value interface{}) {
	head := l.Head
	for head.Next != nil {
		head = head.Next
	}
	head.Next = new(Node)
	head.Next.Value = value
	l.Length++
	return
}
// Print 打印链表
func (l *List) Print() {
	head := l.Head
	for head.Next != nil {
		head = head.Next
		fmt.Print(head.Value, " ")
	}
	fmt.Println()
}
```

#### 3. 测试

```go
func main() {
	L := New()
	for i := 1; i < 6; i++ {
		L.Push(i)
	}
	fmt.Print("        初始化链表: ")
	L.Print()
	fmt.Print("  在链表末尾添加10: ")
	L.Push(10)
	L.Print()
	fmt.Print("  删除链表末尾元素: ")
	L.Pop()
	L.Print()
	fmt.Print("查看第三个节点的值: ")
	fmt.Println(L.Get(3))
	fmt.Print("在第三个节点插入30: ")
	L.Insert(3, 30)
	L.Print()
	fmt.Print("删除第三个节点的值: ")
	L.Delete(3)
	L.Print()
}
```

Output:

```bash
$ go run main.go
       初始化链表: 1 2 3 4 5
  在链表末尾添加10: 1 2 3 4 5 10
  删除链表末尾元素: 1 2 3 4 5
查看第三个节点的值: 3
在第三个节点插入30: 1 2 30 3 4 5
删除第三个节点的值: 1 2 3 4 5
```

