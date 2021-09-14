package tracing

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"sync"
	"time"
)

type Span struct {
	name string

	startTime  time.Time
	finishTime time.Time

	tracer  *Tracer
	spanCtx *SpanContext

	logs    []opentracing.LogRecord
	tags    map[string]interface{}
	baggage sync.Map
}

func (s *Span) Finish() {
	opts := opentracing.FinishOptions{
		FinishTime: time.Now(),
	}

	s.FinishWithOptions(opts)
}

func (s *Span) FinishWithOptions(opts opentracing.FinishOptions) {
	s.finishTime = opts.FinishTime

	if opts.LogRecords != nil {
		s.logs = append(s.logs, opts.LogRecords...)
	}

	s.tracer.Report(s)
}

func (s *Span) SetOperationName(operationName string) opentracing.Span {
	s.name = operationName
	return s
}

func (s *Span) SetTag(key string, value interface{}) opentracing.Span {
	s.tags[key] = value
	return s
}

func (s *Span) LogFields(fields ...log.Field) {
	s.logs = append(s.logs, opentracing.LogRecord{
		Timestamp: time.Now(),
		Fields:    fields,
	})
}

func (s *Span) LogKV(alternatingKeyValues ...interface{}) {
	//TODO: add log depends on value type
	panic("implement me")
	//s.logs=append(s.logs,opentracing.LogRecord{
	//	Timestamp: time.Time{},
	//	Fields:    log.Field{},
	//})
}

func (s *Span) SetBaggageItem(restrictedKey, value string) opentracing.Span {
	s.baggage.Store(restrictedKey, value)
	return s
}

func (s *Span) BaggageItem(restrictedKey string) string {
	res, ok := s.baggage.Load(restrictedKey)
	if !ok {
		return ""
	}
	return res.(string)
}

func (s *Span) Context() opentracing.SpanContext {
	return s.spanCtx
}

func (s *Span) Tracer() opentracing.Tracer {
	return s.tracer
}

func (s *Span) LogEvent(event string) {
	panic("Deprecated API")
}

func (s *Span) LogEventWithPayload(event string, payload interface{}) {
	panic("Deprecated API")
}

func (s *Span) Log(data opentracing.LogData) {
	panic("Deprecated API")
}
