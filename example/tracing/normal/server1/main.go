package main

import (
	"context"
	"fmt"
	"github.com/MashiroC/begonia"
	"github.com/MashiroC/begonia/app/option"
	"os"
)

type TestService struct{}

var echo1 begonia.RemoteFunc
var echo2 begonia.RemoteFunc

func (*TestService) Echo(ctx context.Context, str string) string {
	res, err := echo1(ctx, str)
	if err != nil {
		fmt.Println(err)
		os.Exit(49)
	}
	res2, err := echo2(ctx, str)
	if err != nil {
		fmt.Println(err)
		os.Exit(49)
	}
	return res.(string) + res2.(string) + " from server1"
}

func main() {
	ser := begonia.NewServer(option.Addr(":12306"))

	ser.Register("TEST", &TestService{})

	cli := begonia.NewClient(option.Addr(":12306"))

	s, err := cli.Service("TEST2")
	if err != nil {
		panic(err)
	}
	echo1, err = s.FuncSync("Echo")
	if err != nil {
		panic(err)
	}

	echo2, err = s.FuncSync("Echo2")
	if err != nil {
		panic(err)
	}

	ser.Wait()
}
