package proxy

import (
	"log"
	"time"
)

// 内部 proxy， 可以实现同一个接口，方便替换
// IUser IUser
type IUser interface {
	Login(username, password string) (err error)
}

// User 用户
// @PROXY IUser
type User struct {
}

// Login 用户登录
func (u *User) Login(username, password string) (err error) {
	return nil
}

// UserProxy 代理类
type UserProxy struct {
	user *User
}

// NewUserProxy NewUserProxy
func NewUserProxy(user *User) *UserProxy {
	return &UserProxy{
		user: user,
	}
}

func (p *UserProxy) Login(username, password string) error {
	// before ...
	start := time.Now()

	// 业务逻辑
	if err := p.user.Login(username, password); err != nil {
		return err
	}

	// after ...
	log.Printf("user login cost time: %v", time.Now().Sub(start))

	return nil
}
