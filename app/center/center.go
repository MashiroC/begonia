// Time : 2020/9/19 15:55
// Author : Kieran

// center
package center

import (
	"begonia2/app/core"
	"begonia2/logic"
)

// center.go something

// Center 服务中心的接口，对外统一用接口
type Center interface {
	Run()
}



type rCenter struct {
	lg       logic.MixNode
	services *serviceSet
	Core     core.SubService
	rs       *logic.ReqSet
}

func (c *rCenter) Run() {
	for {
		call, wf := c.lg.RecvMsg()

		go c.work(call, wf)
	}
}

func (c *rCenter) work(call *logic.Call, wf logic.WriteFunc) {

	if call.Service == core.ServiceName {
		// 核心服务
		res, err := c.Core.Invoke(call.Fun, call.Param)
		if err != nil {
			panic(err)
		}

		wf(&logic.CallResult{
			Result: res,
		})
		return
	}

	toID, ok := c.services.Get(call.Service)
	if !ok {
		wf(&logic.CallResult{
			Err: "service not found",
		})
		return
	}

	wf(logic.Redirect, toID)
}
