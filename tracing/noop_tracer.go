package tracing

import (
	"context"
	"github.com/MashiroC/begonia/dispatch/frame"
)

//NoopTracer 默认的空实现
type NoopTracer struct{}

func (n NoopTracer) Start(ctx context.Context, operationName string, opts ...interface{}) (context.Context, Span) {
	return ctx, NoopSpan{}
}
func (n NoopTracer) Inject(sc SpanContext, carrier frame.Request) error { return nil }
func (n NoopTracer) Extract(carrier frame.Request) (SpanContext, error) {
	return NoopSpanContext{}, nil
}
func (n NoopTracer) SpanContextFromContext(ctx context.Context) SpanContext { return NoopSpanContext{} }
func (n NoopTracer) ContextWithSpanContext(ctx context.Context, context SpanContext) context.Context {
	return ctx
}

type NoopSpan struct{}

func (n NoopSpan) Context() SpanContext { return NoopSpanContext{} }
func (n NoopSpan) End()                 {}
func (n NoopSpan) Log(k, v string)      {}
func (n NoopSpan) LogError(err error)   {}

type NoopSpanContext struct{}
