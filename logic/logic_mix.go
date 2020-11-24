package logic

import (
	"github.com/MashiroC/begonia/dispatch"
	"github.com/MashiroC/begonia/dispatch/frame"
	"context"
)

// MixNode 混合节点
type MixNode interface {
	// logic 基础的logic接口
	logic

	// Close 释放资源
	Close()

	// RecvMsg 接收一个消息
	// 只会接收到request，response会直接被转发到对应的连接里。
	RecvCall() (msg *Call, wf ResultFunc)
}

// NewMix 创建一个mix节点
func NewMix(dp dispatch.Dispatcher) MixNode {

	c := &mix{
		baseLogic: baseLogic{
			dp:       dp,
			waitChan: NewWaitChans(),
		},
	}

	return c
}

type mix struct {
	baseLogic
}

func (m *mix) Close() {}

func (m *mix) RecvCall() (call *Call, wf ResultFunc) {
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
						res = frame.NewResponse(f.ReqID, result.Result, result.Err)
					}

					if toConnID != nil {

						for _, toID := range toConnID {

							m.waitChan.AddCallback(context.TODO(), f.ReqID, func(result *CallResult) {
								m.dp.SendTo(connID, frame.NewResponse(f.ReqID, result.Result, result.Err))
							})

							m.dp.SendTo(toID, res)

						}

					} else {

						err := m.dp.SendTo(connID, res)
						if err != nil {
							panic(err)
						}

					}

				},
				ConnID: connID,
				ReqID:  f.ReqID,
			}

			return

		case *frame.Response:

			go func() {
				// 如果是响应包，直接回调
				err := m.waitChan.Callback(f.ReqID, &CallResult{
					Result: f.Result,
					Err:    f.Err,
				})
				if err != nil {
					panic(err)
				}
			}()

			continue

		}

	}

}
