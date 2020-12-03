// Package center default cluster的center节点
package center

import (
	"context"
	"github.com/MashiroC/begonia/core"
	"github.com/MashiroC/begonia/logic"
)

// Center 服务中心的接口，对外统一用接口
type Center interface {
	Run()
}

type rCenter struct {
	ctx context.Context
	cancel context.CancelFunc
	lg *logic.Service
}

func (c *rCenter) Run() {

	c.lg.Hook("dispatch.close", core.C.HandleConnClose)

	c.lg.HandleRequest = c.work

	<-c.ctx.Done()
}

func (c *rCenter) work(call *logic.Call, wf logic.ResultFunc) {
	res, err := core.C.Invoke(wf.ConnID, wf.ReqID, call.Fun, call.Param)
	if err != nil {
		wf.Result(&logic.CallResult{Err: err})
	}

	wf.Result(&logic.CallResult{
		Result: res,
	})
	return
}
