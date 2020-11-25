package main

import (
	"errors"
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
	grpcPort       = 8081
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	fmt.Println(">>>>>> get client request name :"+in.Name)
	return &helloworld.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	grpcHost, err := getLocalIP()
	if err != nil {
		log.Fatalf("failed to getlocalip: %v", err)
	}
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
	fmt.Println("starting ", grpcHost, grpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func getLocalIP() (ipv4 string, err error) {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet // IP地址
		isIpNet bool
	)
	// 获取所有网卡
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}
	// 取第一个非lo的网卡IP
	for _, addr = range addrs {
		// 这个网络地址是IP地址: ipv4, ipv6
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			// 跳过IPV6
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String() // 192.168.1.1
				return
			}
		}
	}

	err = errors.New("not found")
	return
}