package main

import (
	"fmt"
	"github.com/MashiroC/begonia/example/server/call"
)

func main() {
	res, err := call.SayHello("kieran")
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
