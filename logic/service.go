// Time : 2020/10/19 16:50
// Author : Kieran

// logic
package logic

import (
	"begonia2/dispatch"
	"begonia2/dispatch/frame"
	"fmt"
)

// service.go something
type WriteFunc = func(result *CallResult, toConnID ...string)

type Service interface {
	caller
	RecvMsg() (msg *Call, wf WriteFunc)
}

func NewService(dp dispatch.Dispatcher) Service {
	c := &service{
		baseLogic: baseLogic{
			dp:       dp,
			waitChan: NewWaitChans(),
		},
	}
	// TODO: add ctx
	return c
}

type service struct {
	baseLogic
}

func (c *service) RecvMsg() (call *Call, wf WriteFunc) {

	for {

	_, f := c.dp.Recv()

	switch msg :=f.(type) {
	case *frame.Request:

		call = &Call{
			Service: msg.Service,
			Fun:     msg.Fun,
			Param:   msg.Params,
		}

		wf = func(result *CallResult, toConnID ...string) {
			resp := frame.NewResponse(msg.ReqId, result.Result, result.Err)
			if toConnID != nil {
				for _, connID := range toConnID {
					c.dp.SendTo(connID, resp)
				}
			} else {
				c.dp.Send(resp)
			}
		}
		return
	case *frame.Response:
		err := c.waitChan.Callback(msg.ReqId, &CallResult{
			Result: msg.Result,
			Err:    msg.Err,
		})
		fmt.Println("asd")
		if err!=nil{
			panic(err)
		}
	}
	}


}
