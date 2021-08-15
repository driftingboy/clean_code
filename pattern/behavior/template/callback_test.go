package template

import (
	"fmt"
)

func ExampleTemplate() {
	_ = Template("test", func(callId string) error {
		fmt.Printf("do callback func... call id = %s", callId)
		return nil
	})
	//output:
	//template start, test
	//do callback func... call id = 00001
}
