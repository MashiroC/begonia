package main

import (
	"context"
	"fmt"
	"github.com/MashiroC/begonia/example/grpc/echo"
	"google.golang.org/grpc"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var (
	count int32
	flag  bool
	l     sync.Mutex
)

func main() {
	st, err := net.Listen("tcp", ":12306")
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	echo.RegisterGreeterServer(grpcServer, &EchoServer{})

	err = grpcServer.Serve(st)
	if err != nil {
		panic(err)
	}
}

type EchoServer struct{}

func (*EchoServer) SayHello(ctx context.Context, req *echo.HelloRequest) (reply *echo.HelloReply, err error) {
	QPS()
	return &echo.HelloReply{Message: "Hello " + req.Name}, nil
}

func QPS() {
	if !flag {
		l.Lock()
		if !flag {
			flag = true
			go func() {
				time.Sleep(1 * time.Second)
				fmt.Println(count)
				flag = false
				count = 0
			}()
		}
		l.Unlock()
	} else {
		atomic.AddInt32(&count, 1)
	}
}
