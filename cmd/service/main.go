package main

import "begonia2/app/service"

const (
	mode = "cluster"
)

func main() {
	s:=service.New(mode,service.ManagerAddr(":12306"))

	helloService := &HelloService{}

	s.Sign("Hello",helloService)

	s.Wait()
}


type HelloService struct {

}

func (h *HelloService) SayHello(name string) string {
	return "Hello " + name
}