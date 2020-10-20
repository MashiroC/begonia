// Time : 2020/10/19 16:55
// Author : Kieran

// logic
package logic

import (
	"begonia2/config"
	"begonia2/dispatch"
	"begonia2/dispatch/frame"
	"errors"
)

// mix.go something

func NewMix(dp dispatch.Dispatcher) MixNode {
	c := &mix{
		baseLogic: baseLogic{
			dp:       dp,
			waitChan: NewWaitChans(),
		},
		rs:newReqSet(config.C.Logic.RequestTimeOut),
		reqCh: make(chan *frame.Request),
	}
	// TODO: add ctx
	go c.Handle()
	return c
}

type MixNode interface {
	caller
	Handle()
	Close()
	RecvMsg() (msg *Call, wf WriteFunc)
}

type mix struct {
	baseLogic
	reqCh  chan *frame.Request
	rs *reqSet

}

// Handle 处理响应与请求，响应在这里被直接转发，请求则塞到管道里
func (m *mix) Handle() {
//TODO:回收过期的key
	for {
		connID, msg := m.dp.Recv()
		switch f := msg.(type) {
		case *frame.Request:

			// 如果是请求包，记录、发送给上层处理
			m.rs.Add(f.ReqId, connID)
			m.reqCh <- f

		case *frame.Response:

			// 如果是响应包，直接转发
			toID, ok := m.rs.Get(f.ReqId)
			if !ok {
				err := errors.New("connID not found")
				panic(err)
			}

			err := m.dp.SendTo(toID, f)
			if err != nil {
				panic(err)
			}

			m.rs.Remove(f.ReqId)

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

	msg = &Call {
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
			toID,ok:=m.rs.Get(req.ReqId)
			if !ok {
				panic("toID err")
			}
			m.rs.Remove(req.ReqId)

			err := m.dp.SendTo(toID,resp)
			if err!=nil{
				panic(err)
			}
		}
	}

	return
}
