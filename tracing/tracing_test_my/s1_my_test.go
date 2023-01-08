/*
* @Author: DengJie
* @Date:   2023/1/8 17:59
 */
package tracing

import (
	"context"
	"github.com/MashiroC/begonia"
	"github.com/MashiroC/begonia/app/option"
	"testing"
)

type Hello struct {
}

func (receiver Hello) SayName(ctx context.Context, name string) string {
	_, span := (&MyTracer{}).Start(ctx, "in func")
	//do something
	defer span.End()
	return "hello," + name
}

func Test_s1(t *testing.T) {
	cli := begonia.NewServer(option.Addr("127.0.0.1:12306"),
		option.Tracing(&MyTracer{}))
	cli.Register("Hello", Hello{})

	cli.Wait()
}
