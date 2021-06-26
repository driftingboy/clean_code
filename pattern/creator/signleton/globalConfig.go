package signleton

import "sync"

type globalConfig map[string]string

//全局配置
var (
	session globalConfig
	once    sync.Once
)

func Get() globalConfig {
	once.Do(func() {
		session = make(map[string]string)
	})
	return session
}
func Set()       {}
func GetParams() {}
