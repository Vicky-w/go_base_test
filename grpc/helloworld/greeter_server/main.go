package main

import (
	pb "../helloworld"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

//go:generate protoc -I ../helloworld --go_out=plugins=grpc:../helloworld ../helloworld/helloworld.proto

const (
	port2 = ":50051"
)

// server is used to implement helloworld.GreeterServer.  可以理解成proto的service载体
type server2 struct{}

//proto 对应的service内的方法
func (s *server2) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello again " + in.Name}, nil
}

func main() {
	//创建连接
	lis, err := net.Listen("tcp", port2)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server2{})
	// Register reflection service on gRPC server.   反射注册
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
