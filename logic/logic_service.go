package logic

import (
	"context"
	"fmt"
	"github.com/MashiroC/begonia/config"
	"github.com/MashiroC/begonia/dispatch"
	"github.com/MashiroC/begonia/dispatch/frame"
	"github.com/MashiroC/begonia/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"log"
)

// logic_service.go service节点的logic层

// NewService 创建一个实例
func NewService(dp dispatch.Dispatcher, waitChans *CallbackStore) *Service {

	c := &Service{
		Client: &Client{
			Dp:        dp,
			Callbacks: waitChans,

			Tracer: tracing.NewTracer(),
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
		var ctx context.Context
		ctx=context.Background()

		call := &Call{
			Service: msg.Service,
			Fun:     msg.Fun,
			Param:   msg.Params,
		}

		var span opentracing.Span

		wf := func(result Calls) {
			span.Finish()
			resp := result.Frame(msg.ReqID)
			if err := s.Dp.SendTo(connID, resp); err != nil {
				log.Println("err: in send to,", connID, err)
			}
		}

		if config.C.Logic.Tracing.Enable{
			opts := make([]opentracing.StartSpanOption, 0, 2)
			opts = append(opts, ext.SpanKindRPCServer)

			spanCtx, err := s.Tracer.Extract(tracing.Begonia, msg)
			if err == nil && spanCtx != nil {
				opts = append(opts, opentracing.ChildOf(spanCtx))
			}

			fmt.Println("msg:",msg)
			fmt.Println("ctx:",spanCtx)

			span = s.Tracer.StartSpan(fmt.Sprintf("%s.%s", msg.Service, msg.Fun), opts...)

			ctx = context.WithValue(ctx, tracing.Begonia, span)
		}

		ctx = context.WithValue(ctx, "info", map[string]string{"reqID": msg.ReqID, "connID": connID})

		s.HandleRequest(ctx, call, wf)

	case *frame.Response:
		s.Client.HandleResponse(msg)
	}
}
