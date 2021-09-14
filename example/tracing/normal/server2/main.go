package main

import (
	"context"
	"github.com/MashiroC/begonia"
	"github.com/MashiroC/begonia/app/option"
)

type TestService struct{}

func (*TestService) Echo(ctx context.Context, str string) string {
	return str + " from server 2 Echo"
}

func (*TestService) Echo2(ctx context.Context, str string) string {
	return str + " from server 2 Echo 2"
}

func main() {
	ser := begonia.NewServer(option.Addr(":12306"))

	ser.Register("TEST2", &TestService{})

	ser.Wait()
}
