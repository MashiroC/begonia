package tracing

import (
	"errors"
	"fmt"
	"github.com/MashiroC/begonia/dispatch/frame"
	"github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Tracer struct {
	reportChan chan *Span
}

func NewTracer() opentracing.Tracer {
	t := &Tracer{}

	t.reportChan = make(chan *Span, 1024)

	go t.work()

	return t
}

func (t *Tracer) Report(span *Span) {
	t.reportChan <- span
}

func (t *Tracer) work() {
	for span := range t.reportChan {
		fmt.Println("span begin at:", span.startTime.Unix())
		for _, log := range span.logs {
			fmt.Println("log at:", log.Timestamp.Unix())
			fmt.Print("logs:", " ")
			for _, f := range log.Fields {
				fmt.Print(f.String(), " ")
			}
		}
		if span.spanCtx.parent != nil {
			fmt.Println("parentID:", span.spanCtx.parent.SpanID)
		}
		fmt.Println("spanID", span.spanCtx.SpanID)
		fmt.Printf("%+v\n", span.spanCtx)
		fmt.Println("================================")
	}
}

func (t *Tracer) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	sso := opentracing.StartSpanOptions{}

	for _, o := range opts {
		o.Apply(&sso)
	}

	if sso.StartTime.IsZero() {
		sso.StartTime = time.Now()
	}

	var parent *SpanContext
	var traceID, spanID string

	for _, ref := range sso.References {
		spanCtx := ref.ReferencedContext.(*SpanContext)

		if parent == nil && ref.Type == opentracing.ChildOfRef {
			parent = spanCtx
			traceID = parent.TraceID
			spanID = uuid.NewV4().String()
		}
	}

	if traceID == "" {
		traceID = uuid.NewV4().String()
		spanID = uuid.NewV4().String()
	}

	spanCtx := &SpanContext{
		parent:  parent,
		SpanID:  spanID,
		TraceID: traceID,
	}

	s := &Span{
		tracer:  t,
		spanCtx: spanCtx,

		startTime: sso.StartTime,

		logs: make([]opentracing.LogRecord, 0, 4),
		tags: sso.Tags,
	}

	return s
}

func (t *Tracer) Inject(sm opentracing.SpanContext, format interface{}, carrier interface{}) (err error) {
	err = t.checkFormat(format)
	if err != nil {
		return
	}

	req, ok := carrier.(*frame.Request)
	if !ok {
		err = opentracing.ErrInvalidCarrier
		return
	}
	req.Header = make(map[string]string)

	spanCtx, ok := sm.(*SpanContext)
	if !ok {
		err = opentracing.ErrInvalidSpanContext
		return
	}

	req.Header["parentID"] = spanCtx.SpanID
	req.Header["traceID"] = spanCtx.TraceID
	return
}

func (t *Tracer) Extract(format interface{}, carrier interface{}) (spanCtx opentracing.SpanContext, err error) {

	err = t.checkFormat(format)
	if err != nil {
		return
	}

	req, ok := carrier.(*frame.Request)
	if !ok {
		err = opentracing.ErrInvalidCarrier
		return
	}

	header := req.Header

	parentID, ok := header["parentID"]
	if !ok || len(parentID) == 0 {
		errors.New("empty parentID")
		return
	}

	traceID, ok := header["traceID"]
	if !ok || len(traceID) == 0 {
		errors.New("empty traceID")
		return
	}

	spanCtx = &SpanContext{
		parent:  nil,
		ReqID:   req.ReqID,
		SpanID:  parentID,
		TraceID: traceID,
	}

	return
}

func (t *Tracer) checkFormat(format interface{}) (err error) {
	switch format {
	case opentracing.TextMap:
		err = errors.New("not support text")
	case opentracing.HTTPHeaders:
		err = errors.New("not support http")
	case opentracing.Binary:
		err = errors.New("not support binary")
	case Begonia:
		// pass
	default:
		err = opentracing.ErrUnsupportedFormat
	}

	return
}
