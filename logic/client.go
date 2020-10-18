// Time : 2020/9/26 21:26
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

// client.go something

type Client interface {
	CallSync(call *Call) *CallResult
	CallAsync(call *Call, callback Callback)
	Handle()
	Close()
}

type Service interface {
	CallSync(call *Call) *CallResult
	CallAsync(call *Call, callback Callback)
	RecvMsg() (msg *Call, wf WriteFunc)
}

type MixNode interface {
	Client
	Service
}

type WriteFunc = func(result *CallResult)

func NewClient(dp dispatch.Dispatcher) Client {
	c := &client{
		dp:       dp,
		waitChan: NewWaitChans(),
	}
	// TODO: add ctx
	go c.Handle()
	return c
}

func NewService(dp dispatch.Dispatcher) Service {
	c := &client{
		dp:       dp,
		waitChan: NewWaitChans(),
	}
	// TODO: add ctx
	go c.Handle()
	return c
}

type client struct {
	dp       dispatch.Dispatcher
	waitChan *WaitChans
	coderSet *containers.CoderSet
	msgCh    chan *Call
}

func (c *client) CallSync(call *Call) *CallResult {
	ch := make(chan *CallResult)
	defer close(ch)

	c.CallAsync(call, func(result *CallResult) {
		ch <- result
	})

	return <-ch
}

type Callback func(result *CallResult)

func (c *client) CallAsync(call *Call, callback Callback) {

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

func (c *client) Handle() {
	for {
		f := c.dp.Recv()
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

func (c *client) RecvMsg() (msg *Call, wf WriteFunc) {
	f := c.dp.Recv()
	req, ok := f.(*frame.Request)
	if !ok {
		panic("request type error")
	}

	msg = &Call{
		Service: req.Service,
		Fun:     req.Fun,
		Param:   req.Params,
	}

	wf = func(result *CallResult) {
		c.dp.Send(frame.NewResponse(req.ReqId, result.Result, result.Err))
	}

	return
}

func (c *client) Close() {
	c.dp.Close()
}
