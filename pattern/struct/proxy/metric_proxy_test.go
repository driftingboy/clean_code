package proxy

import (
	"fmt"
)

func Example_Generate() {
	src := generate("./user.go")
	fmt.Println(string(src))
	//Output:
	// package proxy

	// type UserProxy struct {
	// 	child *User
	// }

	// func NewUserProxy(child *User) *UserProxy {
	// 	return &UserProxy{child: child}
	// }

	// func (p *UserProxy) Login(username, password string) (err error) {
	// 	// before 这里可能会有一些统计的逻辑
	// 	start := time.Now()

	// 	err = p.child.Login(username, password)

	// 	// after 这里可能也有一些监控统计的逻辑
	// 	log.Printf("user login cost time: %s", time.Now().Sub(start))

	// 	return err
	// }
}
