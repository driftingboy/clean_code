package template

import "fmt"

// 用于扩展功能，外部系统接入...
// 1.函数式callback
type CallBack func(callId string) error

func Template(name string, fn CallBack) error {
	fmt.Printf("template start, %s \n", name)
	// do something get callId
	cid := doSomething()
	return fn(cid)
}

func doSomething() string {
	return "00001"
}

// 2.对象式callback
type TemplateTest struct {
	Name string
	Hook CallBack
}

func NewTemplateTest(Name string, Hook CallBack) *TemplateTest {
	return &TemplateTest{
		Name: Name,
		Hook: Hook,
	}
}

func (t *TemplateTest) Exec() error {
	fmt.Printf("template start, %s \n", t.Name)
	// do something get callId
	cid := doSomething()
	return t.Hook(cid)
}
