package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	pb "github.com/ha1o0/redis-go/protos"
)

func connectRpc() {
	// 连接到 Rust 编写的 gRPC 服务
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	// 创建一个新的 gRPC 客户端
	client := pb.NewGreeterClient(conn)

	// 调用服务端的 SayHello 方法
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "world"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	// 输出服务端返回的结果
	fmt.Printf("Greeting: %s\n", resp.Message)
}
