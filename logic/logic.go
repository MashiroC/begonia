// Package logic 逻辑层
package logic

import (
	"context"
	"github.com/MashiroC/begonia/dispatch"
	"github.com/MashiroC/begonia/dispatch/frame"
	"github.com/MashiroC/begonia/tool/berr"
	"github.com/MashiroC/begonia/tool/ids"
	"log"
	"strings"
)

// Callback logic层回调函数的alias
type Callback = func(result *CallResult)

// baseLogic 基础逻辑层的实现结构体
type baseLogic struct {
	dp       dispatch.Dispatcher // dispatch层的接口，供logic层向下继续调用
	waitChan *WaitChans          // 等待管道，可以在这里注册回调，调用回调
}

func (c *baseLogic) CallSync(call *Call) *CallResult {

	ch := make(chan *CallResult)
	defer close(ch)

	c.CallAsync(call, func(result *CallResult) {
		ch <- result
	})

	return <-ch
}

func (c *baseLogic) CallAsync(call *Call, callback Callback) {

	reqID := ids.New()
	var f frame.Frame
	f = frame.NewRequest(reqID, call.Service, call.Fun, call.Param)

	c.waitChan.AddCallback(context.TODO(), reqID, func(cr *CallResult) {
		callback(cr)
	})

	if err := c.dp.Send(f); err != nil {
		err = c.waitChan.Callback(reqID, &CallResult{
			Result: nil,
			Err:    berr.Warp("logic", "call", err),
		})
		if err != nil {
			// TODO:println => errorln
			log.Println(err)
		}
	}

}

func (c *baseLogic) Hook(typ string, hookFunc interface{}) {

	types := strings.Split(typ, ".")

	if len(types) == 2 {

		switch types[0] {
		case "dispatch":
			c.dp.Hook(types[1], hookFunc)
			return
		}

	} else if len(types) == 1 {
		// 现在logic还没有需要hook的
	}

	panic(berr.NewF("logic", "hook", "func name [%s] not found", typ))

}
