package iterator

import "fmt"

func Example_Iterator() {
	array := []int{1, 2, 3, 4}
	list := NewArrayIntList(array)
	iterator := list.Iterator()
	for iterator.HasNext() {
		fmt.Println(iterator.CurrentItem())
		iterator.Next()
	}
	//Output:
	//1
	//2
	//3
	//4
}
