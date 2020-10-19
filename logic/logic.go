// Time : 2020/10/19 16:51
// Author : Kieran

// logic
package logic

import (
	"begonia2/dispatch"
	"begonia2/dispatch/frame"
	"begonia2/logic/containers"
	"begonia2/tool/ids"
	"context"
)

// logic.go something


type caller interface {
	CallSync(call *Call) *CallResult
	CallAsync(call *Call, callback Callback)
}


type baseLogic struct {
	dp       dispatch.Dispatcher
	waitChan *WaitChans
	coderSet *containers.CoderSet
}

func (c *baseLogic) CallSync(call *Call) *CallResult {
	ch := make(chan *CallResult)
	defer close(ch)

	c.CallAsync(call, func(result *CallResult) {
		ch <- result
	})

	return <-ch
}

type Callback func(result *CallResult)

func (c *baseLogic) CallAsync(call *Call, callback Callback) {

	reqId := ids.New()
	var f frame.Frame
	f = frame.NewRequest(reqId, call.Service, call.Fun, call.Param)

	c.waitChan.AddCallback(context.TODO(), reqId, func(cr *CallResult) {
		callback(cr)
	})

	if err := c.dp.Send(f); err != nil {
		panic(err)
		// TODO:handler error
	}

}