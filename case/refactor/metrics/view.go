package metrics

import (
	"fmt"

	"github.com/hokaccha/go-prettyjson"
)

type Viewer interface {
	Output(map[string]*RequestStatus)
}

// ===============console==============

type consoleViewer struct{}

func NewConsoleViewer() *consoleViewer {
	return &consoleViewer{}
}

func (c consoleViewer) Output(rsMap map[string]*RequestStatus) {
	data, err := prettyjson.Marshal(rsMap)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%v", string(data))
}

// ===============email==============

type emailViewer struct {
	emailAddr string
	emailPwd  string
}

func NewEmailViewer(opts ...viewerOption) *emailViewer {
	ev := &emailViewer{
		emailAddr: "1425895909@qq.com", // 默认配置移到config
		emailPwd:  "123xxxxx",
	}
	for _, o := range opts {
		o(ev)
	}
	return ev
}

func (e emailViewer) Output(rsMap map[string]*RequestStatus) {
	// report data: 邮箱的发送逻辑可能比较复杂，这时候可以抽出一个viewer interface
	// monitor send to email
	fmt.Printf("emailAddr: %s send success. %v", e.emailAddr, rsMap)
}

type viewerOption func(*emailViewer)

func WithEmailAddrAndPwd(emailAddr string, pwd string) viewerOption {
	return func(e *emailViewer) {
		e.emailAddr = emailAddr
		e.emailPwd = pwd
	}
}
