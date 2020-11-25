package main

import (
	"context"
	"fmt"
	"github.com/MashiroC/begonia/example/grpc/echo"
	"google.golang.org/grpc"
	"sync"
	"time"
)

const (
	workLimit = 50
	nodeNums = 5
	workNums  = 1000000
)

func main() {
	conn,err:=grpc.Dial(":12306",grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
	 panic(err)
	}
	client:=echo.NewGreeterClient(conn)


	wg1 := sync.WaitGroup{}

	for i:=0;i<nodeNums;i++{
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
					_,err=client.SayHello(context.Background(),&echo.HelloRequest{
						Name: "shiina",
					})
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

	res, err := client.SayHello(context.Background(), &echo.HelloRequest{Name: "shiina"})
	if err != nil {
	 panic(err)
	}
	fmt.Println(res.Message)
}
