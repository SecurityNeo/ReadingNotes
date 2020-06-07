package main

import "fmt"

/*
当一个数组中大部分元素为0，或者为同一个值得数组时，可以使用稀疏数组来保存该数组。
1、记录数组一共有几行几列，有多少个不同的值；
2、把不同的元素的行列及值记录在一个小规模的数组中，从而缩小程序的规模
3、row    col   value
   11     11      0
   1       2       1
   2       3       2
  这个稀疏数组表示一个11行11列的二维数组，其第二行第三列为1，第三行第四列为2，其余全为0
*/

type valueNode struct {
	row int
	col int
	val int
}

func main() {
	// 先定义一个原始数组
	var chessMap [11][11]int
	chessMap[1][2] = 1
	chessMap[2][3] = 2

	// 初始化一个稀疏数组
	var sparseArray []valueNode
	valNode := valueNode{
		row: len(chessMap),
		col: len(chessMap[0]),
		val: 0,
	}
	sparseArray = append(sparseArray, valNode)

	for idx1, v1 := range chessMap {
		for idx2, v2 := range v1 {
			if v2 != 0 {
				valNode := valueNode{
					row: idx1,
					col: idx2,
					val: v2,
				}
				sparseArray = append(sparseArray, valNode)
			}
		}
	}
	fmt.Println(sparseArray)
}
