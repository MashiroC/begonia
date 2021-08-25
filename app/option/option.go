// Package option starter需要传入的配置相关
package option

import "github.com/MashiroC/begonia/app"

// WriteFunc 拿到的传入参数的map
type WriteFunc func(optionMap map[string]interface{})

// Addr 中心的地址
func Addr(addr string) WriteFunc {
	return func(optionMap map[string]interface{}) {
		optionMap["addr"] = addr
	}
}

func P2P() WriteFunc {
	return func(optionMap map[string]interface{}) {
		optionMap["dpTyp"] = "p2p"
	}
}

// Mode 强制服务以什么方式运行，用于单测等情况下启动服务中心或其他奇怪的情况
func Mode(typ app.ServiceAppModeTyp) WriteFunc {
	return func(optionMap map[string]interface{}) {
		optionMap["mode"]=typ
	}
}