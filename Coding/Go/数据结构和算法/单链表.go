package main

import "fmt"

/*
单链表：
1、头节点不存储数据
2、头节点指向链表中的第一个数据
3、辅助节点，类似于辅助指针，用于链表节点的定位，辅助节点先指向头节点
4、单链表的两种插入方式：
	1）新数据总是在链表尾部插入（insertHero）
	2）新数据按照ID从小到大插入（insertHeroByID）
*/

// 示例： 使用单链表实现水浒英雄管理

type HeroNode struct {
	ID int
	name string
	next *HeroNode
}


func insertHero(head *HeroNode, newHeroNode *HeroNode){
	// 创建辅助节点，类似于辅助指针，用于链表节点的定位，辅助节点先指向头节点
	tempNode := head
	for {
		// 如果当前节点的next为空，说明已经到达链表的末尾，将tempNode.next直接指向newHeroNode即可完成新节点的加入
		if tempNode.next == nil {
			tempNode.next = newHeroNode
			break
		}else {
			// 辅助节点不断往后移
			tempNode = tempNode.next
		}
	}
}

func showHero(head *HeroNode){
	tempNode := head
	if tempNode.next == nil {
		fmt.Println("There is no hero!")
		return
	}
	for {
		// 先打印下一个节点数据，然后辅助节点往后移，当后移之后辅助节点指向为空（tempNode.next），说明已完成整个链表的遍历
		fmt.Printf("[%d %s] --> ",tempNode.next.ID,tempNode.next.name)
		tempNode = tempNode.next
		if tempNode.next == nil {
			break
		}
	}
}

func main() {
	head := &HeroNode{}
	// 定义测试数据
	hero1 := &HeroNode{
		ID: 1,
		name: "宋江",
	}
	hero2 := &HeroNode{
		ID: 2,
		name: "吴用",
	}
	hero3 := &HeroNode{
		ID: 3,
		name: "卢俊义",
	}
	insertHero(head, hero1)
	insertHero(head, hero2)
	insertHero(head, hero3)
	showHero(head)
}