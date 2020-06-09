package main

import "strings"

/*
题目：查找字符串数组中的最长公共前缀。如果不存在公共前缀，返回空字符串 ""。
思路：
	1、借助"strings.Index(s, sep string) int"来判断一个字符串是否以"sep string"开头。函数说明：
		子串sep在字符串s中第一次出现的位置，不存在则返回-1。
	2、先假定字符串数组中的某个字符串（prefix）为公共前缀，那么一定存在其余字符串data
		满足"strings.Index(data, prefix) == 0"。
	3、如果不满足上面的条件，那我们把prefix的最后一个字符串去掉（prefix = prefix[:len(prefix) - 1]），
	   再进行上面的比较。
	4、如果在循环比较过程中出现prefix长度为0（len(prefix) == 0）,说明这个字符串数据中没有用公共前缀。
*/

func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	prefix := strs[0]
	for _,data := range strs {
		for strings.Index(data, prefix) != 0 {
			if len(prefix) == 0 {
				return ""
			}
			prefix = prefix[:len(prefix) -1 ]
		}
	}
	return prefix
}