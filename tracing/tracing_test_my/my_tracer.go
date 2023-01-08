/*
* @Author: DengJie
* @Date:   2023/1/8 17:29
 */
package tracing

import (
	"context"
	"github.com/MashiroC/begonia/dispatch/frame"
	"github.com/MashiroC/begonia/tracing"
	"log"
	"math/rand"
	"strconv"
)

type MyTracer struct {
}

func (m *MyTracer) Start(ctx context.Context, operationName string, opts ...interface{}) (context.Context, tracing.Span) {
	//先拿ctx里面的
	spanContext := ctx.Value("MySpanContext").(MySpanContext)
	fatherID := spanContext.content["SpanID"]
	span := MySpan{content: map[string]string{
		"name":   operationName,
		"father": "my father is " + fatherID,
		"SpanID": strconv.Itoa(rand.Int()),
	}}
	return m.ContextWithSpanContext(ctx, span.Context()), span
}

func (m *MyTracer) Inject(sc tracing.SpanContext, carrier frame.Request) error {
	rsc, ok := sc.(MySpanContext)
	if ok {
		for k, v := range rsc.content {
			carrier.Header[k] = v
		}
	} else {
		carrier.Header["SpanID"] = "-1"
	}
	return nil
}

func (m *MyTracer) Extract(carrier frame.Request) (tracing.SpanContext, error) {
	sc := MySpanContext{map[string]string{}}
	for k, v := range carrier.Header {
		sc.content[k] = v
	}
	return sc, nil
}

func (m *MyTracer) SpanContextFromContext(ctx context.Context) tracing.SpanContext {
	//这里用k为“MySpanContext”
	return ctx.Value("MySpanContext")
}

func (m *MyTracer) ContextWithSpanContext(ctx context.Context, sc tracing.SpanContext) context.Context {
	return context.WithValue(ctx, "MySpanContext", sc)
}

type MySpan struct {
	content map[string]string
}

func (m MySpan) Context() tracing.SpanContext {
	return MySpanContext{
		content: map[string]string{"SpanID": m.content["SpanID"]},
	}
}

func (m MySpan) End() {
	log.Println(m)
}

type MySpanContext struct {
	content map[string]string
}
