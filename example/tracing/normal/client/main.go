package main

import (
	"fmt"
	"github.com/MashiroC/begonia"
	"github.com/MashiroC/begonia/app/option"
)

func main() {
	cli := begonia.NewClient(option.Addr(":12306"))

	s, err := cli.Service("TEST")
	if err != nil {
		panic(err)
	}
	echo, err := s.FuncSync("Echo")
	if err != nil {
		panic(err)
	}

	res, err := echo("kieran")

	fmt.Println(res, err)
}
