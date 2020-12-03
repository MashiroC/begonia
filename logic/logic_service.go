package logic

import (
	"github.com/MashiroC/begonia/dispatch"
	"github.com/MashiroC/begonia/dispatch/frame"
)

// logic_service.go service节点的logic层

// NewService 创建一个实例
func NewService(dp dispatch.Dispatcher, waitChans *WaitChans) *Service {

	c := &Service{
		Client: Client{
			baseLogic: baseLogic{
				dp:       dp,
				waitChan: waitChans,
			},
		},
	}

	dp.Handle("frame", c.DpHandler)

	return c
}

type Service struct {
	Client

	// handle Func
	HandleRequest func(msg *Call, wf ResultFunc)
}

func (s *Service) DpHandler(connID string, f frame.Frame) {
	switch msg := f.(type) {
	case *frame.Request:

		call := &Call{
			Service: msg.Service,
			Fun:     msg.Fun,
			Param:   msg.Params,
		}

		wf := ResultFunc{
			Result: func(result Calls) {
				resp := result.Frame(msg.ReqID)
				s.dp.SendTo(connID, resp)
			},
			ConnID: connID,
			ReqID:  msg.ReqID,
		}

		s.HandleRequest(call, wf)

	case *frame.Response:
		s.Client.HandleResponse(msg)
	}
}
