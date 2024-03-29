package tracing

import (
	"context"
	"github.com/MashiroC/begonia/dispatch/frame"
)

type Tracer interface {
	Start(ctx context.Context, operationName string, opts ...interface{}) (context.Context, Span)
	//Inject from the Context into the carrier
	Inject(sc SpanContext, carrier frame.Request) error
	Extract(carrier frame.Request) (SpanContext, error)
	SpanContextFromContext(ctx context.Context) SpanContext
	ContextWithSpanContext(context.Context, SpanContext) context.Context
}

type Span interface {
	Context() SpanContext
	End()
	Log(k, v string)
	LogError(err error)
}

type SpanContext interface {
}
