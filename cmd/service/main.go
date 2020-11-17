package main

import (
	"begonia2/app/option"
	"begonia2/app/service"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const (
	mode = "center"
)

var (
	count int32
	flag  bool
	l     sync.Mutex
)

func main() {
	count = 0
	flag = false

	s := service.New(mode, option.CenterAddr(":12306"))

	helloService := &HelloService{}

	s.Register("Echo", helloService)

	s.Wait()
}

type HelloService struct {
}

func (h *HelloService) SayHello(name string) string {
	if !flag {
		l.Lock()
		if !flag {
			flag = true
			go func() {
				time.Sleep(1 * time.Second)
				fmt.Println(count)
				flag = false
				count=0
			}()
		}
		l.Unlock()
	} else {
		atomic.AddInt32(&count, 1)
	}
	//fmt.Println("sayHello")
	return "Hello " + name
}

func (h *HelloService) SayHello2(name string) (string, error) {
	fmt.Println("sayHello2")
	return "", errors.New("hello")
}
