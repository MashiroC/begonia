package tracing

import "context"

type registeredTracer struct {
	tracer       Tracer
	isRegistered bool
}

var (
	globalTracer = registeredTracer{NoopTracer{}, false}
)

// SetGlobalTracer 设置一个[单例]的追踪系统。 Tracer 可以使用 GlobalTracer() 返回。
// 在调用`SetGlobalTracer`之前，任何通过`StartSpan`创建的Span都是来自noop的。
func SetGlobalTracer(tracer Tracer) {
	globalTracer = registeredTracer{tracer, true}
}

//GlobalTracer 返回`Tracer`实现的全局单例
//在调用`SetGlobalTracer()`之前，`GlobalTracer()`返回的是noop实现，它会丢掉所有的数据。
func GlobalTracer() Tracer {
	return globalTracer.tracer
}

// StartSpan 遵从 Tracer.StartSpan，见 `GlobalTracer()`。
func StartSpan(ctx context.Context, operationName string, opts ...interface{}) (context.Context, Span) {
	return globalTracer.tracer.Start(ctx, operationName, opts...)
}

// IsGlobalTracerRegistered 返回一个布尔值去判断tracer是否已经在全局注册
func IsGlobalTracerRegistered() bool {
	return globalTracer.isRegistered
}
