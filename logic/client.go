// Time : 2020/9/26 21:26
// Author : Kieran

// logic
package logic

import (
	"begonia2/dispatch"
	"begonia2/dispatch/frame"
)

// baseLogic.go something


type Client interface {
	logic
	Handle()
	Close()
}

type client struct {
	baseLogic
}

func NewClient(dp dispatch.Dispatcher) Client {
	c := &client{
		baseLogic:baseLogic{
			dp:       dp,
			waitChan: NewWaitChans(),
		},
	}
	// TODO: add ctx
	go c.Handle()
	return c
}

func (c *client) Handle() {
	for {
		_,f := c.dp.Recv()
		msg, ok := f.(*frame.Response)
		if !ok {
			panic("response type error")
		}

		reqId := msg.ReqId
		err := c.waitChan.Callback(reqId, &CallResult{
			Result: msg.Result,
			Err:    msg.Err,
		})

		if err != nil {
			panic(err)
		}
	}
}

func (c *client) Close() {
	c.dp.Close()
}
