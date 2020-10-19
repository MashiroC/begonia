// Time : 2020/10/19 16:55
// Author : Kieran

// logic
package logic

import (
	"begonia2/dispatch"
	"begonia2/dispatch/frame"
)

// mix.go something

func NewMix(dp dispatch.Dispatcher) MixNode {
	c := &mix{
		baseLogic: baseLogic{
			dp:       dp,
			waitChan: NewWaitChans(),
		},
	}
	// TODO: add ctx
	go c.Handle()
	return c
}

type MixNode interface {
	caller
	Handle()
	Close()

}

type mix struct {
	baseLogic
	reqCh  chan *frame.Request
	respCh chan *frame.Response
}

func (m *mix) work() {

}

func (m *mix) Handle() {
	for {
		msg,ok:=<-m.respCh
		if !ok {
			panic("response chan error")
		}

		reqId := msg.ReqId
		err := m.waitChan.Callback(reqId, &CallResult{
			Result: msg.Result,
			Err:    msg.Err,
		})

		if err != nil {
			panic(err)
		}
	}

}

func (m *mix) Close() {
	panic("implement me")
}

func (m *mix) RecvMsg() (msg *Call, wf WriteFunc) {
	req, ok := <-m.reqCh
	if !ok {
		panic("request chan close")
	}

	msg = &Call{
		Service: req.Service,
		Fun:     req.Fun,
		Param:   req.Params,
	}

	wf = func(result *CallResult, toConnID ...string) {
		resp := frame.NewResponse(req.ReqId, result.Result, result.Err)
		if toConnID != nil {
			for _, connID := range toConnID {
				m.dp.SendTo(connID, resp)
			}
		} else {
			m.dp.Send(resp)
		}
	}

	return
}
