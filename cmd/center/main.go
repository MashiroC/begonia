package main

import (
	"github.com/MashiroC/begonia/app/center"
	"github.com/MashiroC/begonia/app/option"
)

func main() {

	addr := ":12306"
	c := center.New(option.Addr(addr))

	c.Run()
}
