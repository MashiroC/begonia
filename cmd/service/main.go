package main

import (
	"begonia2/app/service"
	"fmt"
)

const (
	mode = "center"
)

func main() {
	s:=service.New(mode,service.ManagerAddr(":12306"))

	fmt.Println(s)
	helloService := &HelloService{}

	s.Sign("Hello",helloService)

	s.Wait()
}


type HelloService struct {

}

func (h *HelloService) SayHello(name string) string {
	return "Hello " + name
}