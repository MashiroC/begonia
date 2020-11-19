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
	//TestQPS(c)

	in:=testFunc(c, "Test", "Echo2")
	res:=in.([]interface{})
	QPS(c,"Test","Echo",res...)
}

func QPS(c client.Client, service, funName string, param ...interface{}) {
	s, err := c.Service(service)
	if err != nil {
		panic(err)
	}

	testFun, err := s.FuncSync(funName)
	if err != nil {
		panic(err)
	}

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
}

func testFunc(c client.Client, service, funName string, param ...interface{}) interface{} {
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
