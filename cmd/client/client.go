package main

import (
	"begonia2/app/client"
	"begonia2/app/option"
	"fmt"
	"sync"
	"time"
)

const (
	mode = "center"
	addr = ":12306"

	workLimit = 50
	workNums  = 100000
)

func main() {
	c := client.New(mode, option.CenterAddr(addr))


	s, err := c.Service("Echo")
	if err != nil {
		panic(err)
	}

	sayHello, err := s.FuncSync("SayHello")
	if err != nil {
		panic(err)
	}

	//sayHello2, err := s.FuncSync("SayHello2")
	//if err != nil {
	//	panic(err)
	//}

	ch := make(chan struct{}, workLimit)
	for i := 0; i < workLimit; i++ {
		ch <- struct{}{}
	}

	t := time.Now()

	wg:=sync.WaitGroup{}
	for i := 0; i < workNums; i++ {
		<-ch
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
				ch <- struct{}{}
			}()
			res, err := sayHello("shiina")
			if err != nil || res != "Hello shiina"{
				panic(err)
			}
		}()
	}

	wg.Wait()

	fmt.Println(time.Now().Sub(t).String())

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
