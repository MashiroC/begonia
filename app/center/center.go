// Package center default cluster的center节点
package center

import (
	"github.com/MashiroC/begonia/core"
	"github.com/MashiroC/begonia/logic"
	"github.com/MashiroC/begonia/tool/berr"
)

// Center 服务中心的接口，对外统一用接口
type Center interface {
	Run()
}

type rCenter struct {
	lg logic.MixNode
}

func (c *rCenter) Run() {

	c.lg.Hook("dispatch.close", core.C.HandleConnClose)

	for {
		call, wf := c.lg.RecvCall()

		go c.work(call, wf)
	}
}

func (c *rCenter) work(call *logic.Call, wf logic.ResultFunc) {

	if call.Service == core.ServiceName {
		// 核心服务
		res, err := core.C.Invoke(wf.ConnID, wf.ReqID, call.Fun, call.Param)
		if err != nil {
			wf.Result(&logic.CallResult{Err: err})
		}

		wf.Result(&logic.CallResult{
			Result: res,
		})
		return
	}

	toID, ok := core.C.GetToID(call.Service)
	if !ok {
		wf.Result(&logic.CallResult{
			Err: berr.NewF("app.center","work","service [%s] not found",call.Service),
		})
		return
	}

	wf.Result(logic.Redirect, toID)
}
