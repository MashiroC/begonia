package _case

import (
	"github.com/MashiroC/begonia/dispatch"
)

// 一个alias
// 这个包仅做示例，不要依赖这个包下面的所有东西！！！！！！！！！！！
type HandleFunc = func(connID string, data []byte)

// Case 如果注册router中需要调用dispatch等，可以参照这个写法。
func Case(dp dispatch.Dispatcher) func() (code int, fun HandleFunc) {
	return func() (code int, fun HandleFunc) {
		return 7, func(connID string, data []byte) {

		}
	}
}
