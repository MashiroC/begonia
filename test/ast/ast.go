package ast

import (
	"context"
	"errors"
	"github.com/MashiroC/begonia/app/client"
)

var service client.Service

type testRegister struct{}

func (*testRegister) Null() {
	return
}

func (*testRegister) OnlyInput(i int) {

}

func (*testRegister) OnlyOutput() (i int) {
	return 49
}

func (*testRegister) BothInAndOut(i int) (j int) {
	return i + 49
}

func (*testRegister) OutWithError(i int) (j int, err error) {
	if i < 0 {
		return 0, errors.New("error")
	}
	return i, nil
}

func (*testRegister) OutWithContext(ctx context.Context, i int) (s string, err error) {
	v := ctx.Value("info").(map[string]string)
	if v == nil {
		err = errors.New("ctx err")
	}
	return "true", nil
}

func (*testRegister) OutWithErrorAndContext(ctx context.Context, i int) (j int, err error) {
	if i < 0 {
		return 0, errors.New("error")
	}
	return 0, nil
}
