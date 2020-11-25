package main

import (
	"github.com/MashiroC/begonia/app/center"
	"github.com/MashiroC/begonia/app/option"
)



func main() {

	mode := "center"
	addr := ":12306"
	c := center.New(mode, option.CenterAddr(addr))

	c.Run()
}
