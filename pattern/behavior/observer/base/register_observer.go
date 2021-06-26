package base

import (
	"fmt"
)

// 观察者
type RegisterObserver interface {
	// 观察者行为
	HandleRegObsever()
}

type RegPromotionObserver struct{}

func (r *RegPromotionObserver) HandleRegObsever() {
	fmt.Println("你获得了一张优惠券")
	// logger.Info("你获得了一张体验券")
}

type RegNotifyObserver struct{}

func (r *RegNotifyObserver) HandleRegObsever() {
	fmt.Println("emil:欢迎你的注册")
	// logger.Info("你获得了一张体验券")
}

// 被观察者
type Subject interface {
	Register(RegisterObserver)
	Notify()
}

// 注册行为被观察
type RegisterSubject struct {
	ros []RegisterObserver
}

func NewRegisterSubject() *RegisterSubject {
	return &RegisterSubject{ros: make([]RegisterObserver, 0)}
}

func (rs *RegisterSubject) Subscribe(ro RegisterObserver) {
	rs.ros = append(rs.ros, ro)
}

func (rs *RegisterSubject) Notify() {
	for _, ro := range rs.ros {
		ro.HandleRegObsever()
	}
}
