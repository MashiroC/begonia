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
	rs := NewReqSet(config.C.Logic.RequestTimeOut)
	return NewMixWithReqSet(dp, rs)
}

func NewMixWithReqSet(dp dispatch.Dispatcher, rs *ReqSet) MixNode {
	c := &mix{
		baseLogic: baseLogic{
			dp:       dp,
			waitChan: NewWaitChans(),
		},
		rs: rs,
	}
	// TODO: add ctx
	go c.Handle()
	return c
}

type MixNode interface {
	logic
	Close()
	RecvMsg() (msg *Call, wf ResultFunc)
}

type mix struct {
	baseLogic
	rs *ReqSet
}

// Handle 处理响应与请求，请求和响应在这里被直接转发，CoreService的请求被传到上层
func (m *mix) Handle() {

}

func (m *mix) Close() {
	panic("implement me")
}

func (m *mix) RecvMsg() (call *Call, wf ResultFunc) {
	//TODO:回收过期的key

	for {
		/*
			对于mix节点：
			请求：
			1.外界发来给自己处理的
			2.需要被转发的
			(需要到上层进行处理，理论上发送给自己的都是需要自己处理的)
			（logic：记录请求，等待响应，响应可以是自己发的，也可以是其他连接给的）
			（api：看需要转发还是需要直接返回，转发需要得知对方的connID）

			响应：
			1.需要被转发的
			（直接使用回调去转发到chan或者另一条连接，logic处理）
		*/

		connID, msg := m.dp.Recv()
		switch f := msg.(type) {
		case *frame.Request:

			m.rs.Add(f.ReqId, connID)

			call = &Call{
				Service: f.Service,
				Fun:     f.Fun,
				Param:   f.Params,
			}

			wf = ResultFunc{
				Result: func(result *CallResult, toConnID ...string) {

					var res frame.Frame
					if result == Redirect {
						res = f
					} else {
						res = frame.NewResponse(f.ReqId, result.Result, result.Err)
					}

					if toConnID != nil {
						for _, connID := range toConnID {
							m.dp.SendTo(connID, res)
						}
					} else {
						toID, ok := m.rs.Get(f.ReqId)
						if !ok {
							panic("toID err")
						}
						m.rs.Remove(f.ReqId)

						err := m.dp.SendTo(toID, res)
						if err != nil {
							panic(err)
						}
					}

				},
				ConnID: connID,
				ReqID:  f.ReqId,
			}

			return

		case *frame.Response:

			go func() {
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
			}()

			continue

		}

	}

}
