package main

import (
	"context"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"log"
	"sync"
)

const (
	address     = "localhost:50051"

	work = 40
	nums = 100*1000
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	name:="kieran"

	wg:=sync.WaitGroup{}
	for i:=0;i<work;i++{
		wg.Add(1)
		go func() {
			for i:=0;i<nums;i++{
				_, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
				if err!=nil{
					panic(err)
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()

}
