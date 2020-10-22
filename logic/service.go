// Time : 2020/10/19 16:50
// Author : Kieran

// logic
package logic

import (
	"begonia2/config"
	"begonia2/dispatch"
	"begonia2/dispatch/frame"
	"log"
	"time"
)

// service.go something
type ResultFunc struct {
	Result func(result *CallResult, toConnID ...string)
	ConnID string
	ReqID  string
}

type Service interface {
	logic
	RecvMsg() (msg *Call, wf ResultFunc)
}

func NewService(dp dispatch.Dispatcher) Service {
	c := &service{
		baseLogic: baseLogic{
			dp:       dp,
			waitChan: NewWaitChans(),
		},
	}

	if config.C.Logic.AutoReConnection {
		c.dp.Hook("close", func(connID string, err error) {
			ok := false
			if config.C.Logic.ReConnectionRetryCount <= 0 {
				for !ok {
					log.Println("connot link to server,retry...")
					time.Sleep(time.Duration(config.C.Logic.ReConnectionIntervalSecond) * time.Second)
					ok = dp.ReLink()
				}
				return
			} else {
				for i := 0; i < config.C.Logic.ReConnectionRetryCount && !ok; i++ {
					log.Println("connot link to server,retry",i,"limit",config.C.Logic.ReConnectionRetryCount)
					time.Sleep(time.Duration(config.C.Logic.ReConnectionIntervalSecond) * time.Second)
					ok = dp.ReLink()
				}
				if ok {
					return
				} else {
					panic("connect closed")
				}
			}
		})
	}

	// TODO: add ctx
	return c
}

type service struct {
	baseLogic
}

func (c *service) RecvMsg() (call *Call, wf ResultFunc) {

	for {

		connID, f := c.dp.Recv()

		switch msg := f.(type) {
		case *frame.Request:

			call = &Call{
				Service: msg.Service,
				Fun:     msg.Fun,
				Param:   msg.Params,
			}

			wf = ResultFunc{
				Result: func(result *CallResult, toConnID ...string) {
					resp := frame.NewResponse(msg.ReqId, result.Result, result.Err)
					if toConnID != nil {
						for _, connID := range toConnID {
							c.dp.SendTo(connID, resp)
						}
					} else {
						c.dp.Send(resp)
					}
				},
				ConnID: connID,
				ReqID:  msg.ReqId,
			}
			return
		case *frame.Response:
			err := c.waitChan.Callback(msg.ReqId, &CallResult{
				Result: msg.Result,
				Err:    msg.Err,
			})
			if err != nil {
				panic(err)
			}
		}
	}

}
