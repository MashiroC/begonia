/*
* @Author: DengJie
* @Date:   2022/10/31 17:59
 */
package tracing_test_otel

import (
	"context"
	"github.com/MashiroC/begonia"
	"github.com/MashiroC/begonia/app/option"
	"go.opentelemetry.io/otel"
	"log"
	"testing"
)

type Hello struct {
}

func (receiver Hello) SayName(ctx context.Context, name string) string {
	_, span := otel.Tracer("service").Start(ctx, "in func")
	//do something
	defer span.End()
	return "hello," + name
}

func Test_s1(t *testing.T) {
	tp, err := TracerProvider("http://localhost:14268/api/traces",
		"test-s1", "service", 666)
	if err != nil {
		log.Fatal(err)
	}

	otel.SetTracerProvider(tp)
	cli := begonia.NewServer(option.Addr("127.0.0.1:12306"),
		option.TracingWithOtel(tp.Tracer("service")))
	cli.Register("Hello", Hello{})

	cli.Wait()
}
