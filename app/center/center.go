// Time : 2020/9/19 15:55
// Author : Kieran

// center
package center

import (
	"begonia2/logic"
)

// center.go something

// Center 服务中心的接口，对外统一用接口
type Center interface {
	Run()
}

type serviceSet interface {
	Get(service string) (connID string)
	Add(service string)
}

type rCenter struct {
	lg       logic.MixNode
	services serviceSet
	Core     CoreService
}

func (c *rCenter) Run() {
	go c.lg.Handle()
	for {
		call, wf := c.lg.RecvMsg()

		// 核心服务
		res, err := c.Core.Invoke(call.Fun, call.Param)
		if err != nil {
			panic(err)
		}

		wf(&logic.CallResult{
			Result: res,
			Err:    err,
		})

	}
}
