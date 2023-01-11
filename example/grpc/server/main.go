package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	QPS()
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server1 listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

var (
	count int32
	flag  bool
	l     sync.Mutex
)

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
		} else {
			atomic.AddInt32(&count, 1)
		}
		l.Unlock()
	} else {
		atomic.AddInt32(&count, 1)
	}
}
