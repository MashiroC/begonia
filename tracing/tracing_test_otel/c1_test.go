/*
* @Author: DengJie
* @Date:   2022/10/29 16:47
 */
package tracing_test_otel

import (
	"context"
	"fmt"
	"github.com/MashiroC/begonia"
	"github.com/MashiroC/begonia/app/option"
	"go.opentelemetry.io/otel"
	"log"
	"testing"
)

func Test_c1(t *testing.T) {
	client := "client"
	tp, err := TracerProvider("http://localhost:14268/api/traces",
		"test-c1", client, 777)
	if err != nil {
		log.Fatal(err)
	}
	tr := tp.Tracer(client)
	ctx, span := tr.Start(context.Background(), "foo")

	span.End()

	otel.SetTracerProvider(tp)
	c := begonia.NewClient(option.Addr("127.0.0.1:12306"), option.TracingWithOtel(tr))

	s, err := c.Service("Hello")
	if err != nil {
		panic(err)
	}
	fun, err := s.FuncSync("SayName")
	if err != nil {
		panic(err)
	}

	i, err := fun(ctx, "DJ")
	fmt.Println(i)
	i, err = fun(ctx, "DJ")
	fmt.Println(i)
	c.Wait()
}
