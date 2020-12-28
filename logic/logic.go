// Package logic 逻辑层
package logic

import (
	"context"
	"fmt"
	"github.com/MashiroC/begonia/dispatch"
	"github.com/MashiroC/begonia/dispatch/frame"
	"github.com/MashiroC/begonia/tool/ids"
	"log"
	"strings"
)

// Callback logic层回调函数的alias
type Callback = func(result *CallResult)

// BaseLogic 基础逻辑层的实现结构体
type BaseLogic struct {
	Dp        dispatch.Dispatcher // dispatch层的接口，供logic层向下继续调用
	Callbacks *CallbackStore      // 等待管道，可以在这里注册回调，调用回调
}

func (c *BaseLogic) CallSync(call *Call) *CallResult {

	ch := make(chan *CallResult)
	defer close(ch)

	c.CallAsync(call, func(result *CallResult) {
		ch <- result
	})

	return <-ch
}

func (c *BaseLogic) CallAsync(call *Call, callback Callback) {

	reqID := ids.New()
	var f frame.Frame
	f = frame.NewRequest(reqID, call.Service, call.Fun, call.Param)

	c.Callbacks.AddCallback(context.TODO(), reqID, func(cr *CallResult) {
		callback(cr)
	})

	if err := c.Dp.Send(f); err != nil {
		err = c.Callbacks.Callback(reqID, &CallResult{
			Result: nil,
			Err:    fmt.Errorf("logic call error: %w", err),
		})
		if err != nil {
			// TODO:println => errorln
			log.Println(err)
		}
	}

}

func (c *BaseLogic) Hook(typ string, hookFunc interface{}) {

	types := strings.Split(typ, ".")

	if len(types) == 2 {

		switch types[0] {
		case "dispatch":
			c.Dp.Hook(types[1], hookFunc)
			return
		}

	} else if len(types) == 1 {
		// 现在logic还没有需要hook的
	}

	panic(fmt.Sprintf("hook func [%s] not found", typ))
}

func (c *BaseLogic) Handle(typ string, handleFunc interface{}) {
	types := strings.Split(typ, ".")

	if len(types) == 2 {

		switch types[0] {
		case "dispatch":
			c.Dp.Handle(types[1], handleFunc)
			return
		}

	} else if len(types) == 1 {
		// 现在logic还没有需要handle的
	}

	panic(fmt.Sprintf("handle func [%s] not found", typ))

}
