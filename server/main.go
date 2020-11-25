package main

import (
	"fmt"
	"github.com/liming8519/grpc-demo/consul"
	"github.com/liming8519/grpc-demo/helloworld"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)


const (
	grpcHost       = "127.0.0.1"
	grpcPort       = 8081
)


type server struct{}

func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	fmt.Println(">>>>>> get client request name :"+in.Name)
	return &helloworld.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	lis, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(grpcHost), grpcPort, ""})
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	cr := consul.NewConsulRegister(consul.Addr, 15)
	cr.Register(consul.RegisterInfo{
		Host:           grpcHost,
		Port:           grpcPort,
		ServiceName:    consul.ServiceName,
		UpdateInterval: time.Second})


	helloworld.RegisterGreeterServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}