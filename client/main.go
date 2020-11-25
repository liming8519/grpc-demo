package main

import (
	"fmt"
	"github.com/liming8519/grpc-demo/consul"
	"github.com/liming8519/grpc-demo/helloworld"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"log"
	"time"
)


func main()  {
	callConsulGrpc()
}

func callConsulGrpc()  {
	schema, err := consul.GenerateAndRegisterConsulResolver(consul.Addr, consul.ServiceName)
	if err != nil {
		log.Fatal("init consul resovler err", err.Error())
	}

	//建立连接
	conn, err := grpc.Dial(fmt.Sprintf("%s:///%s", schema, consul.ServiceName), grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := helloworld.NewGreeterClient(conn)

	name := "ok"

	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		//调用远端的方法
		r, err := c.SayHello(ctx, &helloworld.HelloRequest{Name: name})
		if err != nil {
			log.Println("could not greet: %v", err)

		} else {
			log.Printf("Hello: %s", r.Message)
		}
		time.Sleep(time.Second)
	}

}