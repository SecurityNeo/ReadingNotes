package main

import "fmt"

// 冒泡排序
func bubbleSort(arr *[5]int) {
	length := len(arr)
	for i := 0;i < length;i++ {
		for j := i+1;j < length;j++ {
			if arr[i] > arr[j] {
				arr[i],arr[j] = arr[j],arr[i]
			}
		}
	}
	fmt.Println(arr)
}

// 选择排序
func selectSort(arr *[5]int){
	for j:=0;j < len(arr) - 1;j++ {
		// 每一次先假定第一个数（剩余未排序部分）是最小的
		tempVal := arr[j]
		tempIndex := j
		// 遍历数组，找出真正最小的那个数
		for i := j+1;i < len(arr);i++ {
			if arr[i] < tempVal {
				tempVal = arr[i]
				tempIndex = i
			}
		}
		// 如果找出真正最小那个数不是我们最开始假定的那个数，才进行交换（优化）
		if tempIndex != j {
			arr[j],arr[tempIndex] = arr[tempIndex],arr[j]
		}
	}
	fmt.Println(arr)
}

func main(){
	arr  := [5]int{4,-8,9,3,-19}
	arr2  := [5]int{2,-8,899,-3,90}
	bubbleSort(&arr)
	selectSort(&arr2)
}
