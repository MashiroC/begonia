package main

import (
	"begonia2/app/client"
	"begonia2/app/option"
	"fmt"
)

const (
	mode = "center"
	addr = ":12306"
)

func main() {
	c := client.New(mode, option.ManagerAddr(addr))

	fmt.Println(c)

	go func() {
		//time.Sleep(5*time.Second)
		//c.Close()
	}()

	s,err:=c.Service("Echo")
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
	c.Wait()

	//s, err := c.Service("Hello")
	//if err != nil {
	//	panic(err)
	//}
	//
	//sayHello, err := s.FuncSync("SayHello")
	//if err != nil {
	//	panic(err)
	//}
	//
	//res, err := sayHello("shiina")
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println(res)

}
