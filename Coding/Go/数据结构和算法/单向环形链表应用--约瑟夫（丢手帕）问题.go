package main

import "fmt"

/*
单向环形链表的应用：约瑟夫问题（丢手帕问题）
 */


type Student struct {
	ID int
	Next *Student
}


// AddStudent: 添加人数，构成单向环形链表。num表示链表中的个数，函数返回链表头节点指针
func AddStudent(num int) *Student{
	var first = &Student{}
	helper := first
	if num < 1 {
		fmt.Printf("数据无效，退出\n")
		return nil
	}
	for i := 1;i <= num; i++{
		// 构造节点
		stu := &Student{
			ID : i,
		}
		// 第一个节点比较特殊，需单独处理
		if i == 1 {
			first = stu
			helper = stu
			helper.Next = first
		}else {
			helper.Next = stu
			// 辅助节点往后移动一下
			helper = stu
			// 最后一个节点指向头节点构成环形
			helper.Next = first
		}
	}
	return first
}

// 展示链表中的所有元素
func ShowStudent(first *Student){
	if first == nil {
		fmt.Printf("空链表！\n")
		return
	}
	helper := first
	for {
		if helper.Next == first {
			// 把最后一个数据打印之后再退出
			fmt.Printf("ID %d --> ", helper.ID)
			break
		}
		fmt.Printf("ID %d --> ", helper.ID)
		helper = helper.Next
	}
	fmt.Println()
}


// PlayGame: startNum表示从哪一个开始报数，count表示数多少个数
func PlayGame(first *Student,startNum int, count int){
	nodeNum := 0

	// 定义两个辅助节点
	head := first
	tail := first
	// 将tail定位到链表尾巴，同时计算链表中有nodeNum和节点
	for {
		if tail.Next == first {
			nodeNum ++
			break
		}
		tail = tail.Next
		nodeNum ++
	}
	if startNum > nodeNum {
		fmt.Printf("节点总数为 %d，输入超出节点总数!\n", nodeNum)
		return
	}
	// 将辅助节点head（头）定位到开始那个数数那个节点
	for i := 1;i <= startNum; i++ {
		// 如果每次从1号开始数的话，相当于head和tail都不动
		if i == 1 {
			continue
		}
		head = head.Next
		tail = tail.Next
	}
	for {
		if head == tail {
			// 将最后留下的节点打印后才结束
			fmt.Printf("ID %d 为最后一个\n",head.ID)
			return
		}
		for i := 1;i <= count;i++ {
			if i == 1{
				continue
			}else {
				head = head.Next
				tail = tail.Next
			}
		}
		fmt.Printf("ID %d 出列 \n",head.ID)
		head = head.Next
		tail.Next = head
	}
}

func main(){
	first := AddStudent(5)
	ShowStudent(first)
	PlayGame(first,6,2)
}