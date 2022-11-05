package logic

import (
	"context"
	"fmt"
	"github.com/MashiroC/begonia/dispatch"
	"github.com/MashiroC/begonia/dispatch/frame"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"log"
)

// logic_service.go service节点的logic层

// NewService 创建一个实例
func NewService(dp dispatch.Dispatcher, waitChans *CallbackStore, tracer trace.Tracer) *Service {

	c := &Service{
		Client: &Client{
			Dp:        dp,
			Callbacks: waitChans,
		},
	}

	if tracer != nil {
		c.HasTracer = true
		c.Tracer = tracer
		c.PropagateBy = propagation.TraceContext{}
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

		var span trace.Span

		wf := func(result Calls) {
			if span != nil {
				span.End()
			}
			resp := result.Frame(msg.ReqID)
			if err := s.Dp.SendTo(connID, resp); err != nil {
				log.Println("err: in send to,", connID, err)
			}
		}

		if s.HasTracer {

			var m propagation.MapCarrier
			m = msg.Header
			ctx = s.PropagateBy.Extract(ctx, m)
			ctx, span = s.Tracer.Start(ctx, fmt.Sprintf("%s.%s", msg.Service, msg.Fun))
		}

		ctx = context.WithValue(ctx, "info", map[string]string{"reqID": msg.ReqID, "connID": connID})

		s.HandleRequest(ctx, call, wf)

	case *frame.Response:
		s.Client.HandleResponse(msg)
	}
}
