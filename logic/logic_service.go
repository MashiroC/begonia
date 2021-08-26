package logic

import (
	"context"
	"github.com/MashiroC/begonia/dispatch"
	"github.com/MashiroC/begonia/dispatch/frame"
)

// logic_service.go service节点的logic层

// NewService 创建一个实例
func NewService(dp dispatch.Dispatcher, waitChans *CallbackStore) *Service {

	c := &Service{
		Client: &Client{
			Dp:        dp,
			Callbacks: waitChans,
		},
	}

	dp.Handle("frame", c.DpHandler)

	return c
}

type Service struct {
	*Client

	// handle Func
	HandleRequest func(ctx context.Context, msg *Call, wf ResultFunc)
}

func (s *Service) DpHandler(connID string, f frame.Frame) {
	switch msg := f.(type) {
	case *frame.Request:

		call := &Call{
			Service: msg.Service,
			Fun:     msg.Fun,
			Param:   msg.Params,
		}

		wf := func(result Calls) {
			resp := result.Frame(msg.ReqID)
			s.Dp.SendTo(connID, resp)
		}

		ctx := context.WithValue(context.Background(), "info", map[string]string{"reqID": msg.ReqID, "connID": connID})

		s.HandleRequest(ctx, call, wf)

	case *frame.Response:
		s.Client.HandleResponse(msg)
	}
}
