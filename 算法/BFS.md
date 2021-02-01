BFS

#### 1. 树的最小高度



```go
// 二叉树
type TreeNode struct {
	Value interface{}
	Left,Right *TreeNode
}

// 队列
type Queue []*TreeNode

func (q *Queue)Size()int {
	return len(*q)
}

// 出队
func (q *Queue)Poll() *TreeNode{
	var length = len(*q)
	if length == 0 {
		return nil
	}
	var data = (*q)[0]
	*q = (*q)[1:length]
	return data
}

//入队
func (q *Queue) Offer(node *TreeNode) {
	*q = append(*q, node)
	return
}

// 计算树最小高度
func MinDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}
	// 初始化一个队列记录树的节点
	var q = new(Queue)
	// 根节点入队
	q.Offer(root)
	var depth = 1
	for q.Size() != 0 {
		sz := q.Size()
		for i := 0 ; i < sz ; i ++ {
			node := q.Poll()
			if node.Left == nil && node.Right == nil {
				return depth
			}
			if node.Left != nil {
				q.Offer(node.Left)
			}
			if node.Right != nil {
				q.Offer(node.Right)
			}
		}
		depth++
	}
	return depth
}
```

