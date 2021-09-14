package tracing

type SpanContext struct {
	parent *SpanContext
	ReqID string
	SpanID string
	TraceID string

}

func (s *SpanContext) ForeachBaggageItem(handler func(k string, v string) bool) {
	panic("implement me")
}

func NewSpanContext(){

}