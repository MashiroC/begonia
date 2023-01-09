/*
* @Author: DengJie
* @Date:   2023/1/8 17:57
 */
package tracing

import (
	"context"
	"fmt"
	"github.com/MashiroC/begonia"
	"github.com/MashiroC/begonia/app/option"
	"testing"
)

func Test_c1(t *testing.T) {
	tracer := MyTracer{}
	c := begonia.NewClient(option.Addr("127.0.0.1:12306"), option.Tracing(&tracer))

	s, err := c.Service("Hello")
	if err != nil {
		panic(err)
	}
	fun, err := s.FuncSync("SayName")
	if err != nil {
		panic(err)
	}

	funWithErr, err := s.FuncSync("SayNameWithError")
	if err != nil {
		panic(err)
	}

	ctx, span := tracer.Start(context.Background(), "step 1")
	i, err := fun(ctx, "DJ")
	span.End()
	fmt.Println(i, err)
	i, err = funWithErr(context.Background(), "DJ")
	fmt.Println(i, err)
	c.Wait()
}
