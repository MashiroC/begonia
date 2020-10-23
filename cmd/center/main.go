package main

import (
	"begonia2/app/center"
	"begonia2/app/option"
)

func main() {

	mode := "center"
	addr := ":12306"
	c := center.New(mode, option.CenterAddr(addr))

	c.Run()
}
