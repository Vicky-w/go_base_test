package main

import (
	"log"
	"os"
	"time"

	pb "../helloworld"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address2     = "localhost:50051"
	defaultName2 = "world"
)

func main() {
	// Set up a connection to the server.  创建连接
	conn, err := grpc.Dial(address2, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn) //创建对应service的cli

	// Contact the server and print out its response.
	//处理参数
	name := defaultName2
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	//连接设置
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//请求内容
	r, err := c.SayHelloAgain(ctx, &pb.HelloRequest{Name: name}) //return HelloReply
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
}
