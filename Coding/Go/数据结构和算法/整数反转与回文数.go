package main

import "math"

/*
整数反转(reverse)：
	题目：给出一个32位的有符号整数，你需要将这个整数中每位上的数字进行反转。假设我们的环境只能
          存储得下32位的有符号整数，则其数值范围为 [−231,  231 − 1]。请根据这个假设，如果反转后
          整数溢出那么就返回 0。
	思路：
	    1、公式: (y = y * 10 + x % 10),(x = x / 10)
		2、借助math包"math.MaxInt32"和"math.MinInt32"来进行数据溢出判断，或者位运算(1<<31 -1)和(-1<< 31)）
*/

/*
回文数（isPalindrome）：
	题目：判断一个整数是否是回文数。回文数是指正序（从左向右）和倒序（从右向左）读都是一样的整数。
	思路：
		1、借助整数反转思路将原数进行反转，然后比较反转后的数与原数是否相等；
		2、0是回文数，负数肯定不是回文数，末尾为0的数（x % 10）不是回文数（除开0）；
		3、不用将原数完全反转，实际上反转一半就能判断是否为回文数，可以节省一半的消耗
*/

func reverse(x int) int {
	if x == 0 {
		return 0
	}
	y := 0
	for x != 0 {
		y = y * 10 + x % 10
		x = x / 10
	}
	if y < math.MinInt32 || y > math.MaxInt32 {
		return 0
	}
	return y
}

func isPalindrome(x int) bool {
	if x == 0 {
		return true
	}
	y := 0
	if x < 0 || x % 10 == 0 {
		return false
	}

	for x > y {
		y = y * 10 + x % 10
		x = x / 10
	}
	if y == x || y /10 ==x {
		return true
	}
	return false
}