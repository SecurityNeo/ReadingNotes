package main

import (
	"fmt"
)

/*
双向链表：
1、头节点不存储数据
2、头节点指向链表中的第一个数据
3、辅助节点，类似于辅助指针，用于链表节点的定位，辅助节点先指向头节点
4、双向链表的两种插入方式：
	1）新数据总是在链表尾部插入（insertHero）
	2）新数据按照ID从小到大插入（insertHeroByID）
		注意：顺序插入时建议先将待插入节点的pre指向它前一个节点，next指向它后一个节点，然后再修改其前一个
		和后一个节点分别指向自己。
5、删除节点时一定要注意判断待删除节点是不是链表的最后一个节点
*/

// 示例： 使用单链表实现水浒英雄管理

type HeroNode struct {
	ID   int
	name string
	next *HeroNode
	pre  *HeroNode
}

func insertHero(head *HeroNode, newHeroNode *HeroNode) {
	// 创建辅助节点，类似于辅助指针，用于链表节点的定位，辅助节点先指向头节点
	tempNode := head
	// 首先找到最后边的节点
	for {
		if tempNode.next == nil {
			break
		}
		tempNode = tempNode.next
	}
	//
	tempNode.next = newHeroNode
	newHeroNode.pre = tempNode
}

// 按照ID从小到大的顺序插入新节点
func insertHeroByID(head *HeroNode, newHeroNode *HeroNode) {
	tempNode := head
	for {
		if tempNode.next == nil {
			tempNode.next = newHeroNode
			newHeroNode.pre = tempNode
			break
			// 如果顺序为从大到小，修改">"为"<"即可
		} else if tempNode.next.ID > newHeroNode.ID {
			newHeroNode.next = tempNode.next
			newHeroNode.pre = tempNode
			tempNode.next.pre = newHeroNode
			tempNode.next = newHeroNode
			break
		} else if tempNode.next.ID == newHeroNode.ID {
			fmt.Printf("ID %d conflict!\n", newHeroNode.ID)
			break
		}
		tempNode = tempNode.next
	}
}

// 根据ID删除链表中的数据
func deleteHeroNodeByID(head *HeroNode, id int){
	tempNode := head
	for {
		if tempNode.next == nil {
			fmt.Printf("ID %d does not exist! \n", id)
			break
		}else if tempNode.next.ID == id {
			if tempNode.next.next == nil {
				tempNode.next = nil
				return
			}else {
				tempNode.next = tempNode.next.next
				tempNode.next.pre =  tempNode
				return
			}
		}
		tempNode = tempNode.next
	}
}

// 顺序打印
func showHero(head *HeroNode) {
	tempNode := head
	if tempNode.next == nil {
		fmt.Println("There is no hero!")
		return
	}
	for {
		// 先打印下一个节点数据，然后辅助节点往后移，当后移之后辅助节点指向为空（tempNode.next），说明已完成整个链表的遍历
		fmt.Printf("[%d %s] --> ", tempNode.next.ID, tempNode.next.name)
		tempNode = tempNode.next
		if tempNode.next == nil {
			break
		}
	}
	fmt.Println()
}

// 逆序打印
func showHeroInReverse(head *HeroNode) {
	tempNode := head
	if tempNode.next == nil {
		fmt.Println("There is no hero!")
		return
	}
	// 先找到链表末尾
	for {
		if tempNode.next == nil {
			break
		}
		tempNode = tempNode.next
	}
	for {
		// 直接打印当前节点数据，然后辅助节点往前移，当前移之后辅助节点指向为空（tempNode.pre），说明已完成整个链表的遍历
		fmt.Printf("[%d %s] --> ", tempNode.ID, tempNode.name)
		tempNode = tempNode.pre
		if tempNode.pre == nil {
			break
		}
	}
	fmt.Println()
}

func main() {
	head := &HeroNode{}
	// 定义测试数据
	hero1 := &HeroNode{
		ID:   1,
		name: "宋江",
	}
	hero2 := &HeroNode{
		ID:   2,
		name: "吴用",
	}
	hero3 := &HeroNode{
		ID:   3,
		name: "卢俊义",
	}
	hero4 := &HeroNode{
		ID:   3,
		name: "林冲",
	}
	insertHeroByID(head, hero1)
	insertHeroByID(head, hero3)
	insertHeroByID(head, hero2)
	insertHeroByID(head, hero4)
	showHero(head)
	showHeroInReverse(head)
	deleteHeroNodeByID(head, 3)
	showHero(head)
	showHeroInReverse(head)
}
