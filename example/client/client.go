package main

import (
	"fmt"
	"github.com/MashiroC/begonia"
	"github.com/MashiroC/begonia/app/option"
)

func main() {
	c := begonia.NewClient(option.Addr(":12306"))

	s, err := c.Service("Echo")
	if err != nil {
		panic(err)
	}

	testFun, err := s.FuncSync("SayHello")
	if err != nil {
		panic(err)
	}

	res, err := testFun("kieran")
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}
