package main

import "fmt"

/*
单向环形链表：
	1、与单链表非常相似，只是单向环形链表的最后一个节点再指向头节点（X.next = head）；
	2、注意，单向环形链表的头节点是要存储数据的；
	3、当链表中只有一个节点时，它需要指向自己(head.next = head)
*/

type HeroNode struct {
	ID int
	name string
	next *HeroNode
}


func InsertHeroNode(head *HeroNode, newHeroNode *HeroNode) {
	if head.next == nil {
		head.ID = newHeroNode.ID
		head.name = newHeroNode.name
		head.next = head
		return
	}
	tempNode := head
	for {
		if tempNode.next == head {
			tempNode.next = newHeroNode
			newHeroNode.next = head
			return
		}
		tempNode = tempNode.next
	}
}

func showHeroNode(head *HeroNode){
	if head.next == nil {
		fmt.Println("Empty!")
		return
	}
	tempNode := head
	for {
		fmt.Printf("[ID: %d Name: %s] -->",tempNode.ID, tempNode.name)
		if tempNode.next == head {
			break
		}
		tempNode = tempNode.next
	}
}

func main() {
	head := &HeroNode{}
	hero1 := &HeroNode{
		ID: 1,
		name: "宋江",
	}
	hero2 := &HeroNode{
		ID: 2,
		name: "卢俊义",
	}
	hero3 := &HeroNode{
		ID: 3,
		name: "林冲",
	}
	showHeroNode(head)
	InsertHeroNode(head, hero1)
	InsertHeroNode(head, hero2)
	InsertHeroNode(head, hero3)
	showHeroNode(head)

}