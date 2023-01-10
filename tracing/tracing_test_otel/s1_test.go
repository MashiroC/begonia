package tracing_test_otel

import (
	"context"
	"errors"
	"github.com/MashiroC/begonia"
	"github.com/MashiroC/begonia/app/option"
	"github.com/MashiroC/begonia/tracing"
	"go.opentelemetry.io/otel"
	"log"
	"testing"
)

type Hello struct {
}

func (receiver Hello) SayName(ctx context.Context, name string) string {
	//开启一个子span
	_, span1 := tracing.GlobalTracer().Start(ctx, "in func")
	//do something
	defer span1.End()
	//通过ctx拿到span
	span := ctx.Value("span").(tracing.Span)
	span.Log("get span by ctx", "successfully")
	return "hello," + name
}

func (receiver Hello) SayNameWithError(ctx context.Context, name string) (string, error) {
	return "", errors.New("test span record the call rpc err")
}

type HelloWithoutTracer struct{}

func (h HelloWithoutTracer) SayName(ctx context.Context, name string) string {
	return name
}

func Test_s1(t *testing.T) {
	tp, err := TracerProvider("http://localhost:14268/api/traces",
		"test-s1", "service", 666)
	if err != nil {
		log.Fatal(err)
	}

	otel.SetTracerProvider(tp)
	server := begonia.NewServer(option.Addr("127.0.0.1:12306"),
		option.TracingWithOtel(tp.Tracer("service")))
	server.Register("Hello", Hello{})

	go server.Wait()

	//一个没有tracer的server
	server1 := begonia.NewServer(option.Addr("127.0.0.1:12306"),
		option.TracingWithOtel(tp.Tracer("service")))
	server1.Register("HelloWithoutTracer", HelloWithoutTracer{})

	server1.Wait()
}
