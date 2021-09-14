package tracing

import (
	"fmt"
	"github.com/MashiroC/begonia/dispatch/frame"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTracing(t *testing.T) {
	tracer := NewTracer()
	opentracing.SetGlobalTracer(tracer)

	var span opentracing.Span
	var err error
	var spanCtx opentracing.SpanContext

	a := assert.New(t)
	/*  1 链路起点创建一个Span
	    2 记录信息
		3 头传递tracing id和span id，记录这次请求，链路到下一节点
		4 节点解包
		5 重复 2～4
		6 链路终点

	*/

	// 1
	span = opentracing.StartSpan("step1")

	a.NotNil(span.Context())
	a.NotNil(span.Tracer())

	time.Sleep(2 * time.Second)
	// 2
	span.LogFields(log.String("hello", "world"), log.Int32("step", 1))

	time.Sleep(2 * time.Second)
	span.Finish()

	// 3
	f := frame.NewRequest("testReq", "TEST", "TEST", []byte{1, 2, 3})
	req := f.(*frame.Request)
	err = opentracing.GlobalTracer().Inject(span.Context(), Begonia, f)

	a.Nil(err)
	a.NotNil(req.Header)
	a.NotZero(len(req.Header))

	// 4
	spanCtx, err = opentracing.GlobalTracer().Extract(Begonia, f)
	a.Nil(err)
	a.NotNil(spanCtx)

	// 5 - 2
	span = opentracing.GlobalTracer().StartSpan("step5", opentracing.ChildOf(spanCtx))
	span.LogFields(log.String("hello", "world"), log.Int32("step", 2))

	span.Finish()

	// 5 - 3
	f2 := frame.NewRequest("testReq2", "TEST", "TEST", []byte{1, 2, 3})
	req2 := f2.(*frame.Request)
	err = opentracing.GlobalTracer().Inject(span.Context(), Begonia, f2)

	a.Nil(err)
	a.NotNil(req2.Header)
	a.NotZero(len(req2.Header))

	// 6
	a.Equal(req2.Header["traceID"], req.Header["traceID"])
	a.NotEqual(req2.Header["parentID"], req.Header["parentID"])

	myTracer := tracer.(*Tracer)
	close(myTracer.reportChan)

	for record := range myTracer.reportChan {
		fmt.Println("record begin at:", record.startTime.Unix())
		for _, log := range record.logs {
			fmt.Println("log at:", log.Timestamp.Unix())
			fmt.Print("logs:"," ")
			for _, f := range log.Fields {
				fmt.Print(f.String(), " ")
			}
		}
		fmt.Println()
		if record.ctx.parent != nil {
			fmt.Println("parentID:",record.ctx.parent.SpanID)
		}
		fmt.Printf("%+v\n", record.ctx)
		fmt.Println("================================")
	}
}

func TestTags(t *testing.T) {
	tracer:= NewTracer()
	opentracing.SetGlobalTracer(tracer)

	span:=opentracing.StartSpan("test-tag",ext.SpanKindRPCServer)
	span.Finish()
}