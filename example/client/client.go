package main

import (
	"fmt"
	"github.com/MashiroC/begonia"
	"github.com/MashiroC/begonia/app/option"
)

func main() {
	c := begonia.NewClient(option.Addr(":12306"))

	s, _ := c.Service("Echo")

	SayHello, _ := s.FuncSync("SayHello")
	res, err := SayHello("kieran")
	fmt.Println(res)
	fmt.Println(err)

	//SayHelloAsync, _ := s.FuncAsync("SayHello")
	//SayHelloAsync(func(res interface{}, err error) {
	//	fmt.Println(res, err)
	//}, "kieran")

}
