// Package option starter需要传入的配置相关
package option

// WriteFunc 拿到的传入参数的map
type WriteFunc func(optionMap map[string]interface{})

// CenterAddr 中心的地址
func CenterAddr(addr string) WriteFunc {
	return func(optionMap map[string]interface{}) {
		optionMap["managerAddr"] = addr
	}
}

func P2P() WriteFunc{
	return func(optionMap map[string]interface{}) {
		optionMap["dpTyp"]="p2p"
	}
}