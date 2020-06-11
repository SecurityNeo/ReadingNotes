package main

import (
	"fmt"
)

/*
单向环形链表：
	1、与单链表非常相似，只是单向环形链表的最后一个节点再指向头节点（X.next = head）；
	2、注意，单向环形链表的头节点是要存储数据的；
	3、当链表中只有一个节点时，它需要指向自己(head.next = head)
	4、需要定义两个辅助节点，与单项链表不一样的是，这里的头节点是要存储数据的。先将helper节点定位
       到头节点的上一个节点，后边helper节点会跟随tempNode一起移动。
	5、当需要删除的节点就是头节点时，需要将头节点的下一个节点设置为新的头节点。我们这里所有的逻辑
       都基于有一个头节点，所以头节点不能为空。我们在函数内修改了头节点，就需要将新的头节点返回去，
       后边的操作就必须用这里返回去的新头节点。
*/

type HeroNode struct {
	ID int
	name string
	next *HeroNode
}


func InsertHeroNode(head *HeroNode, newHeroNode *HeroNode) {
	// 链表为空时，头节点的下一个节点指向自己即可
	if head.next == nil {
		head.ID = newHeroNode.ID
		head.name = newHeroNode.name
		head.next = head
		return
	}
	// 同样定义一个辅助节点，用来查找谁指向了头节点
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

func deleteHeroNode(head *HeroNode, ID int) *HeroNode{
	if head.next == nil {
		fmt.Printf("Empty! \n")
		return head
	}
	// 链表中只有一个节点时，设（head.next = nil）即可
	if head.next == head {
		head.next = nil
		return head
	}
	// 这里需要定义两个辅助节点，与单项链表不一样的是，这里的头节点是要存储数据的。
	tempNode := head
	helperNode := head
	// 先将helper节点定位到头节点的上一个节点，后边helper节点会跟随tempNode一起移动。
	for {
		if helperNode.next == head {
			break
		}
		helperNode = helperNode.next
	}
	flag := true
	for {
		// 这个时候只比较到倒数第二个节点，最后一个节点还没比较，只是找到了最后一个节点
		if tempNode.next == head {
			break
		}
		if tempNode.ID == ID {
			// 这里需要特别注意，当需要删除的节点就是头节点时，需要将头节点的下一个节点设置为新的头节点。
			// 我们这里所有的逻辑都基于有一个头节点，所以头节点不能为空。我们在函数内修改了头节点，就需要
			// 将新的头节点返回去，后边的操作就必须用这里返回去的新头节点。
			if tempNode == head {
				head = head.next
			}
			helperNode.next = tempNode.next
			flag = false
			break
		}
		tempNode = tempNode.next
		helperNode = helperNode.next
	}
	// 如果flag为true，说明在上边这个for循环中没有进行节点删除。那我们要把最后一个节点再比较一次。
	if flag {
		if tempNode.ID == ID {
			helperNode.next = tempNode.next
		}else {
			// 说明比较了所有节点，仍然没有找到对应ID的节点
			fmt.Printf("ID %d does not exist!\n", ID)
		}
	}
	return  head
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
	fmt.Println()
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
	head = deleteHeroNode(head, 20)
	showHeroNode(head)
	head = deleteHeroNode(head, 1)
	showHeroNode(head)
	head = deleteHeroNode(head, 3)
	showHeroNode(head)

}