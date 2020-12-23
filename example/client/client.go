package main

import (
	"fmt"
	"github.com/MashiroC/begonia"
	"github.com/MashiroC/begonia/app/option"
	"sync"
	"time"
)

const (
	addr = ":12306"

	workLimit = 50
	nodeNums  = 5
	workNums  = 1000000
)

func main() {
	c := begonia.NewClient(option.Addr(addr))
	//TestQPS(c)

	//res := testFunc(c, "Test", "Echo2")
	//fmt.Println(res.([]interface{})[12])
	//in := testFunc(c, "Test", "Echo2")
	//res := in.([]interface{})
	//fmt.Println(res)
	//fmt.Println(testFunc(c, "Test", "Echo", res...))
	//QPS(c,"Test","Echo",res...)
	//QPS(c, "Echo", "SayHello", "shiina")
	fmt.Println(testFunc(c, "Echo", "SayHello", "shiina"))
}

func QPS(c begonia.Client, service, funName string, param ...interface{}) {
	s, err := c.Service(service)
	if err != nil {
		panic(err)
	}

	testFun, err := s.FuncSync(funName)
	testFunAsync, err := s.FuncAsync(funName)
	testFunAsync(func(res interface{}, err error) {
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
	}, "kieran")
	if err != nil {
		panic(err)
	}

	wg1 := sync.WaitGroup{}

	for i := 0; i < nodeNums; i++ {
		wg1.Add(1)
		go func() {
			ch := make(chan struct{}, workLimit)
			for i := 0; i < workLimit; i++ {
				ch <- struct{}{}
			}

			t := time.Now()

			wg := sync.WaitGroup{}
			for i := 0; i < workNums; i++ {
				<-ch
				wg.Add(1)
				go func() {
					defer func() {
						wg.Done()
						ch <- struct{}{}
					}()
					if len(param) != 0 {
						_, err = testFun(param...)
					} else {
						_, err = testFun()
					}
					if err != nil {
						panic(err)
					}
				}()
			}

			wg.Wait()

			fmt.Println(time.Now().Sub(t).String())

			wg1.Done()
		}()
	}
	wg1.Wait()
}

func testFunc(c begonia.Client, service, funName string, param ...interface{}) interface{} {
	s, err := c.Service(service)
	if err != nil {
		panic(err)
	}

	testFun, err := s.FuncSync(funName)
	if err != nil {
		panic(err)
	}

	var res interface{}
	if len(param) != 0 {
		res, err = testFun(param...)
	} else {
		res, err = testFun()
	}
	if err != nil {
		panic(err)
	}

	return res
}
