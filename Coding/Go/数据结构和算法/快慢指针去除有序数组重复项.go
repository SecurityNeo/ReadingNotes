package main

import "fmt"

/*
使用快慢指针将一个有序数组去重
1、快指针（right）：找到陌生值，一个right指针从未发现过的值。
2、慢指针（left）：告诉right指针，发现新的值就放在我这里，因为我（left）的左边都是不重复的值，或者说都是你（right）一路发现的不同的值。
3、快指针每次移动1，
4、慢指针则是每接收一次值，移动1。
5、nums[left:right]区间的值都是慢指针（left）存储过的值，所以不用担心慢指针会不会覆盖别的值。
6、因为nums[0]是第一个出现的元素，不存在与前面的元素重复的可能，所以快慢指针初始化为1。
*/


func removeDuplicates(nums []int) int {
	var right = 1
	var left = 1
	n := len(nums)
	if n < 2 {
		return n
	}
	for right < n {
		if nums[right] != nums[right - 1] {
			nums[left]  = nums[right]
			left++
		}
		right++
	}
	return left
}

func main(){
	// 初始化一个测试数组
	arr := []int{2,2,2,4,5,5,5,6,7}
	n := removeDuplicates(arr)
	fmt.Println("The length of the array: ",n)
	fmt.Println("The new array is :",arr[:n])

}