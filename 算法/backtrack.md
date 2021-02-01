回溯算法

**问题描述：**

全排列

Go 语言实现：

```go
var res [][]int

func Permute(nums []int) {
	var track []int
	BackTrack(nums, track)
	return
}

// 路径：记录在 track 中
// 选择列表：nums 中不存在于 track 的那些元素
// 结束条件：nums 中的元素全部都在 track 中出现
func BackTrack(nums, track []int) {
	if len(track) == len(nums) {
		res = append(res, track)
		return
	}
	for _, num := range nums {
		if Contains(track, num) {
			continue
		}
		track = append(track, num)
		// 递归所有路径
		BackTrack(nums, track)
		// 递归结束归还最后一个选择
		track = track[:len(track) - 1]
	}
	return
}

func Contains(track []int, num int) bool {
	for _, v := range track {
		if v == num {
			return true
		}
	}
	return false
}
```



N 皇后问题