package main

import (
	"begonia2/app/center"
	"begonia2/app/option"
	"fmt"
)

func main() {

	mode:="center"
	addr:=":12306"
	fmt.Println("new")
	c:=center.New(mode,option.ManagerAddr(addr))

	c.Run()
}
