package signleton

import "sync"

type globalConfig map[string]string

//全局配置
var (
	appConfig globalConfig
	once      sync.Once
)

// 饿汉式
func init() {
	// 从配置文件读取配置到 appConfig
}

// 懒汉式
func Get() globalConfig {
	once.Do(func() {
		appConfig = make(map[string]string)
		// 从配置文件读取配置到 appConfig
	})
	return appConfig
}

// 需要并发安全，加写锁
func Set() {}

// 需要并发安全，加读锁
func GetParams() {}
