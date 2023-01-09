/*
* @Author: DengJie
* @Date:   2023/1/8 17:59
 */
package tracing

import (
	"context"
	"errors"
	"github.com/MashiroC/begonia"
	"github.com/MashiroC/begonia/app/option"
	"github.com/MashiroC/begonia/tracing"
	"testing"
)

type Hello struct {
}

func (receiver Hello) SayName(ctx context.Context, name string) string {
	//开启一个子span
	_, span := tracing.GlobalTracer().Start(ctx, "in func")
	//do something
	defer span.End()
	//通过ctx拿到span
	span = ctx.Value("span").(tracing.Span)
	span.Log("get span by ctx", "successfully")
	return "hello," + name
}

func (receiver Hello) SayNameWithError(ctx context.Context, name string) (string, error) {
	//测试返回err写入span
	return "", errors.New("test span record the call rpc err")
}

func Test_s1(t *testing.T) {
	cli := begonia.NewServer(option.Addr("127.0.0.1:12306"),
		option.Tracing(&MyTracer{}))
	cli.Register("Hello", Hello{})

	cli.Wait()
}
