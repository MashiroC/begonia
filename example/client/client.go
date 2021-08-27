package main

import (
	"fmt"
	"github.com/MashiroC/begonia"
	"github.com/MashiroC/begonia/app"
	"github.com/MashiroC/begonia/app/client"
	"github.com/MashiroC/begonia/app/option"
	"github.com/MashiroC/begonia/example/server/call"
	"sync"
)

func main() {
	c := begonia.NewClient(option.Addr(":12306"),option.Mode(app.Reflect))

	s, _ := c.Service("Echo")

	SayHello, _ := s.FuncSync("SayHello")
	res, err := SayHello("kieran")
	fmt.Println(res)
	fmt.Println(err)
	//time.Sleep(100*time.Minute)
	//QPS(SayHello,"kieran")
	//SayHelloAsync, _ := s.FuncAsync("SayHello")
	//SayHelloAsync(func(res interface{}, err error) {
	//	fmt.Println(res, err)
	//}, "kieran")

}

func QPS(f client.RemoteFunSync,params ...interface{}){
	wg:=sync.WaitGroup{}

	work:=30

	nums:=100*1000


	for i:=0;i<work;i++{
		wg.Add(1)
		if i%2==0{
			go func() {
				for i:=0;i<nums;i++{
					//_,err:=f(params...)
					_,err := call.SayHello("kieran")
					if err != nil {
						panic(err)
					}
				}
				wg.Done()
			}()
		}else{
			go func() {
				for i:=0;i<nums;i++{
					_,err:=f(params...)
					//_,err := call.SayHello("kieran")
					if err != nil {
						panic(err)
					}
				}
				wg.Done()
			}()
		}


	}

	wg.Wait()
}