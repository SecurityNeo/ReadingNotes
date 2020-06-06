package main

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
)

/*
思路：
	1、使用双指针分别表示队列最后一个数据（tail）和最前面一个数据（head），maxSize表示队列最大长度
	2、初始化时： head = 0,tail = 0
	3、队列为空的条件： head == tail
	4、队列满的条件：(tail + 1) % maxSize == head
	5、队列的元素个数： (tail + maxSize - head) % maxSize
*/

type CircleQueue struct {
	maxSize int
	head int		// 指向队列队首
	tail int		// 指向队列队尾
	array [5]int	// 表示此环形队列只能保存4个数据
}

//  AddQueue: 向队列中添加一个数据
func (C *CircleQueue) AddQueue(data int) error {
	if C.IsFull() {
		return errors.New("Queue is full!")
	}
	C.array[C.tail] = data
	C.tail = (C.tail + 1) % C.maxSize
	return nil
}

// PopQueue：向队列中取出一个数据
func (C *CircleQueue) PopQueue() (data int, err error){
	if C.IsEmpty() {
		err = errors.New("Queue is empty!")
		return
	}
	data = C.array[C.head]
	C.head = (C.head + 1) % C.maxSize
	return
}

// IsFull: 判断队列是否已满
func (C *CircleQueue) IsFull() bool {
	return ( C.tail + 1 ) % C.maxSize == C.head
}

// IsEmpty： 判断队列是否为空
func (C *CircleQueue) IsEmpty() bool{
	return C.tail == C.head
}

// QueueSize：查询队列中有多少个数据
func (C *CircleQueue) QueueSize() int {
	return ( C.tail + C.maxSize - C.head ) % C.maxSize
}

// ShowQueue: 展示队列中所有数据
func (C *CircleQueue) ShowQueue() {
	if C.IsEmpty() {
		fmt.Printf("Queue is empty!\n")
		return
	}
	queueSize := C.QueueSize()
	tempHead := C.head
	fmt.Println("The data of the queue is:")
	for i := 0; i < queueSize; i++ {
		fmt.Printf("arr[%d] = %d\t",tempHead, C.array[tempHead])
		tempHead = ( tempHead + 1 ) % C.maxSize
	}
	fmt.Println()
}

func main() {
	// 初始化一个环形队列
	queue := &CircleQueue{
		maxSize: 5,
		tail: 0,
		head: 0,
	}

	var input string
	var data int
	for {
		fmt.Println("1.Input 1 to add a data to the circle queue.")
		fmt.Println("2.Input 2 to pop a data to the circle queue.")
		fmt.Println("3.Input 3 to show the circle queue.")
		fmt.Println("4.Input 4 to exit.")

		fmt.Scanln(&input)
		switch input {
		case "1" :
			fmt.Println("Pls input the number:")
			fmt.Scanln(&data)
			err := queue.AddQueue(data)
			if err != nil {
				fmt.Println(err)
			}else{
				fmt.Println("Add data successfully!")
			}
		case "2":
			val, err := queue.PopQueue()
			if err != nil {
				fmt.Println(err)
			}else{
				fmt.Printf("Get a data from the circle queue: %d\n",val)
			}
		case "3":
			queue.ShowQueue()
		case "4":
			os.Exit(0)
		}
	}
}
