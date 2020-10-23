package main

import (
	"begonia2/app/option"
	"begonia2/app/service"
	"errors"
	"fmt"
)

const (
	mode = "center"
)

func main() {
	s := service.New(mode, option.CenterAddr(":12306"))

	helloService := &HelloService{}

	s.Register("Echo", helloService)

	s.Wait()
}

type HelloService struct {
}

func (h *HelloService) SayHello(name string) string {
	fmt.Println("sayyHello")
	return "Hello " + name
}

func (h *HelloService) SayHello2(name string) (string, error) {
	return "", errors.New("hello")
}
