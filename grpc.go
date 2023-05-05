package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/ha1o0/redis-go/protos"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedGreeterServer // 继承自生成的服务端接口
}

// 实现服务端接口中的方法
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloResponse{Message: "Hello " + in.GetName()}, nil
}

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

func listenRpc() {
	// 监听网络连接
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// 创建 gRPC 服务器
	s := grpc.NewServer()
	// 注册服务端接口
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("Server listening at %v", lis.Addr())
	// 启动 gRPC 服务器
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
