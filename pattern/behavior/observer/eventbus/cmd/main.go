package main

import (
	"fmt"
	"github/litao-2071/clean_code/pattern/behavior/observer/eventbus"
	"time"
)

func Notify(topic string) {
	fmt.Printf("你订阅的%s模块更新了，快去看看\n", topic)
}

func writeLog(topic string) {
	fmt.Printf("%s 订阅日志已记录", topic)
}

func main() {
	bus := eventbus.NewSimpleBus()

	// user1
	bus.Subscribe("golang", Notify)
	bus.Subscribe("docker", Notify)

	// user2
	bus.Subscribe("golang", Notify)

	bus.Publish("golang", "golang")
	time.Sleep(1 * time.Second)
}
