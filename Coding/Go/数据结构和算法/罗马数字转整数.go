package main

/*
题目： 罗马数字转整数
		字符          数值
		I             1
		V             5
		X             10
		L             50
		C             100
		D             500
		M             1000
		通常情况下，罗马数字中小的数字在大的数字的右边。以下六种情况特殊：
		I 可以放在 V (5) 和 X (10) 的左边，来表示 4 和 9。
		X 可以放在 L (50) 和 C (100) 的左边，来表示 40 和 90。 
		C 可以放在 D (500) 和 M (1000) 的左边，来表示 400 和 900。
思路：
	1、罗马数字每个字符对应一个十进制数字，将输入的罗马数字转换成对应数值的切片（有序的），
	   方便进行加减运算；
	2、遍历转换后的切片，如果存在前面一个数字小于后边一个数字，那就需要减去前面那个数字的两倍（因为
       正常应该是总数减去前面那个数，但是在这之前做了一次加法）。
	3、对于string，range迭代的是Unicode而不是字节，返回的值是rune，此处还需将其转换
       为string（string(data) ）
	4、在做特殊情况判断时，我们判断的是当前这个数和它的前一个数，所以，至少是从第二个数开始比
       较（if i > 0 && slice[i] > slice[i - 1]）
*/

func romanToInt(s string) int {
	var result int
	slice := make([]int,1)
	for _,data := range s {
		switch string(data) {
		case "I":
			slice = append(slice, 1)
		case "V":
			slice = append(slice, 5)
		case "X":
			slice = append(slice, 10)
		case "L":
			slice = append(slice, 50)
		case "C":
			slice = append(slice, 100)
		case "D":
			slice = append(slice, 500)
		case "M":
			slice = append(slice, 1000)
		}
	}
	for i := 0;i < len(slice);i++ {
		if i > 0 && slice[i] > slice[i - 1]{
			result = result + (slice[i] - 2 * slice[i - 1])
		}else {
			result = result + slice[i]
		}
	}

	return result
}

