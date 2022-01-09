package command

import (
	"fmt"
)

type Interface interface {
	Exec()
}

// 开始游戏
type StartCommand struct{}

func (s StartCommand) Exec() {
	fmt.Println("game start")
}

// 获取钻石
type GotDiamondCommand struct{}

func (g GotDiamondCommand) Exec() {
	fmt.Println("you got a diamond")
}

// 存档
type ArchiveCommand struct{}

func (a ArchiveCommand) Exec() {
	fmt.Println("archive success")
}
