// Time : 2020/9/19 15:55
// Author : Kieran

// center
package center

import (
	"begonia2/logic"
	"fmt"
)

// center.go something

// Center 服务中心的接口，对外统一用接口
type Center interface {
	Run()
}

type rCenter struct {
	lg logic.MixNode
}

func (c *rCenter) Run() {
	go c.lg.Handle()
	for {
		call,wf:=c.lg.RecvMsg()

		fmt.Println("call:",call)

		wf(&logic.CallResult{
			Result: []byte{1,2,3},
			Err:    "",
		})
		fmt.Println("zxc")
	}
}
