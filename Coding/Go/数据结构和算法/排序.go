package main

import "fmt"

// 冒泡排序
/*
1、比较相邻的元素。如果第一个比第二个大，就交换他们两个。
2、对每一对相邻元素作同样的工作，从开始第一对到结尾的最后一对。这步做完后，最后的元素会是最大的数。
3、针对所有的元素重复以上的步骤，除了最后一个。
4、持续每次对越来越少的元素重复上面的步骤，直到没有任何一对数字需要比较。
*/
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

// 选择排序：是一种简单直观的排序算法，无论什么数据进去都是 O(n²) 的时间复杂度。所以用到它的时候，数据规模越小越好。
/*
1、首先在未排序序列中找到最小（大）元素，存放到排序序列的起始位置。
2、再从剩余未排序元素中继续寻找最小（大）元素，然后放到已排序序列的末尾。
3、重复第二步，直到所有元素均排序完毕。
*/
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

// 插入排序

/*
1、将第一待排序序列第一个元素看做一个有序序列，把第二个元素到最后一个元素当成是未排序序列。
2、从头到尾依次扫描未排序序列，将扫描到的每个元素插入有序序列的适当位置。（如果待插入的元素与有序序列中的某个元素相等，则将待插入元素插入到相等元素的后面。）
*/

func insertSort(arr *[5]int){
	for i := 1; i < len(arr);i++{
		insertIndex := i - 1
		insertVal := arr[i]
		for insertIndex >= 0 && arr[insertIndex] > insertVal {
			arr[insertIndex + 1] = arr[insertIndex]
			insertIndex -= 1
		}
		if insertIndex + 1 != i {
			arr[insertIndex+1] = insertVal
		}
	}
	fmt.Println(arr)
}

func main(){
	arr  := [5]int{4,-8,9,3,-19}
	arr2  := [5]int{2,-8,899,-3,90}
	arr3  := [5]int{5,2,-9,-20,99}
	bubbleSort(&arr)
	selectSort(&arr2)
	insertSort(&arr3)
}
