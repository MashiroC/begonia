package logic

import (
	"context"
	"fmt"
	"github.com/MashiroC/begonia/dispatch"
	"github.com/MashiroC/begonia/dispatch/frame"
	"github.com/MashiroC/begonia/tracing"
	"log"
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

		ctx := context.Background()

		call := &Call{
			Service: msg.Service,
			Fun:     msg.Fun,
			Param:   msg.Params,
		}

		var span tracing.Span

		wf := func(result Calls) {
			if span != nil {
				cr, ok := result.(*CallResult)
				if ok && cr.Err != nil {
					span.LogError(cr.Err)
				}
				//span.Log("end at", time.Now().String())
				span.End()
			}
			resp := result.Frame(msg.ReqID)
			if err := s.Dp.SendTo(connID, resp); err != nil {
				log.Println("err: in send to,", connID, err)
			}
		}

		//这里也可以不判断使用NoopTracer这个空实现
		if tracing.IsGlobalTracerRegistered() {
			spanCtx, err := tracing.GlobalTracer().Extract(*msg)
			if err != nil {
				log.Println(err)
			} else {
				ctx, span = tracing.GlobalTracer().Start(tracing.GlobalTracer().ContextWithSpanContext(ctx, spanCtx), fmt.Sprintf("%s.%s", msg.Service, msg.Fun))
				ctx = context.WithValue(ctx, "span", span)
				span.Log("kind", "rpc")
				//span.Log("start at", time.Now().String())
			}
		}

		ctx = context.WithValue(ctx, "info", map[string]string{"reqID": msg.ReqID, "connID": connID})

		s.HandleRequest(ctx, call, wf)

	case *frame.Response:
		s.Client.HandleResponse(msg)
	}
}
