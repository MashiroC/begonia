// Package option starter需要传入的配置相关
package option

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

func SetLogService() WriteFunc {
	return func(optionMap map[string]interface{}) {
		optionMap["log"] = true
	}
}
