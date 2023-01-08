/*
* @Author: DengJie
* @Date:   2023/1/8 17:22
 */
package conv

import (
	"context"
	"github.com/MashiroC/begonia/dispatch/frame"
	"github.com/MashiroC/begonia/tracing"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"log"
)

//这里让Otel满足我们的接口 tracing/Tracer
type OtelTracer struct {
	Tr          trace.Tracer
	PropagateBy propagation.TextMapPropagator
}

func (m *OtelTracer) SpanContextFromContext(ctx context.Context) tracing.SpanContext {
	return trace.SpanContextFromContext(ctx)
}

func (m *OtelTracer) ContextWithSpanContext(parent context.Context, spanContext tracing.SpanContext) context.Context {
	return trace.ContextWithRemoteSpanContext(parent, spanContext.(trace.SpanContext))
}

func (m *OtelTracer) Start(ctx context.Context, operationName string, opts ...interface{}) (newCtx context.Context, span tracing.Span) {
	var ms = mySpan{span: nil}
	if opts != nil {
		defer func() {
			e := recover()
			if e != nil {
				log.Printf("Start Span Error:opts error [%s]\n", e)
			}
		}()
		var realOpts = make([]trace.SpanStartOption, len(opts))
		for i := 0; i < len(realOpts); i++ {
			realOpts[i] = opts[i].(trace.SpanStartOption)
		}
		newCtx, ms.span = m.Tr.Start(ctx, operationName, realOpts...)
	} else {
		newCtx, ms.span = m.Tr.Start(ctx, operationName)
	}
	span = ms
	return
}

func (m *OtelTracer) Inject(sc tracing.SpanContext, carrier frame.Request) error {
	var mc propagation.MapCarrier
	//rpc 默认远程把
	mc = carrier.Header
	ctx := m.ContextWithSpanContext(context.Background(), sc)
	m.PropagateBy.Inject(ctx, mc)
	return nil
}

func (m *OtelTracer) Extract(carrier frame.Request) (tracing.SpanContext, error) {
	var mm propagation.MapCarrier
	mm = carrier.Header
	ctx := m.PropagateBy.Extract(context.Background(), mm)
	return m.SpanContextFromContext(ctx), nil
}

type mySpan struct {
	span trace.Span
}

func (ms mySpan) Context() tracing.SpanContext {
	return ms.span.SpanContext()
}

func (ms mySpan) End() {
	ms.span.End()
}
