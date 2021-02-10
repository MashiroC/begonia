package main

import (
	"fmt"
	"github.com/MashiroC/begonia"
	"github.com/MashiroC/begonia/app/option"
	"time"
)

func main() {
	c := begonia.NewClient(option.Addr(":12306"))
	time.Sleep(time.Hour)
	s, _ := c.Service("Echo")

	SayHello, _ := s.FuncSync("SayHello")
	res, _ := SayHello("kieran")
	fmt.Println(res.(string))

	SayHelloAsync, _ := s.FuncAsync("SayHello")
	SayHelloAsync(func(res interface{}, err error) {
		fmt.Println(res, err)
	}, "kieran")

}
