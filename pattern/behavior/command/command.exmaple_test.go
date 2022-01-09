package command_test

import (
	"fmt"
	"github/litao-2071/clean_code/pattern/behavior/command"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// 模拟
	eventsChan := make(chan Event)
	defer close(eventsChan)
	events := []Event{Start, GotDiamond, Archive, Start, Start, Archive, Archive, Archive, GotDiamond, Archive}
	go func() {
		for _, event := range events {
			eventsChan <- event
		}
	}()

	// 接收请求，生成命令存储
	commandChan := make(chan command.Interface, 100)
	defer close(commandChan)
	go func() {
		for event := range eventsChan {
			command, err := NewCommandByEvent(event)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			commandChan <- command
		}
	}()

	// 执行命令(
	//	1.如果exec是一个耗时操作， 可以多个 worker 并发消费
	//  2.注意设定每天指令的超时时间，超时则命令执行失败，不要影响其他命令执行
	// )
	for {
		select {
		case cmd := <-commandChan:
			cmd.Exec()
		case <-time.After(2 * time.Second):
			fmt.Println("time out 2s, has no command")
			return
		}
	}

}

type Event int

const (
	Start Event = iota + 1
	GotDiamond
	Archive
)

// factory
func NewCommandByEvent(typ Event) (command.Interface, error) {
	switch typ {
	case Start:
		return &command.StartCommand{}, nil
	case GotDiamond:
		return &command.GotDiamondCommand{}, nil
	case Archive:
		return &command.ArchiveCommand{}, nil
	default:
		return nil, fmt.Errorf("no the command type: %v", typ)
	}
}
