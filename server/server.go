package main

// import (
// 	"context"
// 	"log"
// 	"net"

// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/reflection"

// 	// pb "./hello/proto/hello_grpc.pb.go"
// 	// pb "./hello/proto"
// 	pb "server/hello/meety"
// )

// type HelloworldHandler struct {
// 	pb.UnimplementedGreeterServer
// }

// func (h HelloworldHandler) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error) {
// 	return &pb.HelloReply{Message: "Hello " + request.Name}, nil
// }

// // func (h HelloworldHandler) SayRepeatHello(request *pb.RepeatHelloRequest, server pb.Greeter_SayRepeatHelloServer) error {
// // 	panic("implement me")
// // }


// func main() {
// 	port := ":9090"
// 	lis, err := net.Listen("tcp", port)
// 	if err != nil {
// 		log.Fatal("failed to listen %v", err)
// 	}
// 	server := grpc.NewServer()
// 	pb.RegisterGreeterServer(server, &HelloworldHandler{})
// 	reflection.Register(server)

// 	log.Printf("start gRPC server")
// 	server.Serve(lis)
// }