package memento_test

import (
	"fmt"
	"github/litao-2071/clean_code/pattern/behavior/memento"
)

func Example_Memento() {
	io := &memento.InputOperater{}
	io.Append("hello")
	fmt.Println(io.View())
	io.Snapshot()

	io.Append("world")
	fmt.Println(io.View())
	io.Snapshot()

	io.Append("V2")
	fmt.Println(io.View())
	io.RestoreLast()
	fmt.Println(io.View())

	io.Restore(0)
	fmt.Println(io.View())
	//Output:
	// hello
	// helloworld
	// helloworldV2
	// helloworld
	// hello
}
