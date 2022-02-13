package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"
	pb "github.com/Ryuichi-g/meety_server/grpc/messenger"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
    _ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
    port = ":9090"
)

type server struct {
    pb.UnimplementedMessengerServer
    requests []*pb.MessageRequest
}

func (s *server) GetMessages(_ *empty.Empty, stream pb.Messenger_GetMessagesServer) error {
    for _, r := range s.requests {
        if err := stream.Send(&pb.MessageResponse{Message: r.GetMessage()}); err != nil {
            return err
        }
    }

    previousCount := len(s.requests)

    for {
        currentCount := len(s.requests)
        if previousCount < currentCount {
            r := s.requests[currentCount-1]
            log.Printf("Sent: %v", r.GetMessage())
            if err := stream.Send(&pb.MessageResponse{Message: r.GetMessage()}); err != nil {
                return err
            }
        }
        previousCount = currentCount
    }
}

func (s *server) CreateMessage(ctx context.Context, r *pb.MessageRequest) (*pb.MessageResponse, error) {
    log.Printf("Received: %v", r.GetMessage())
    fmt.Printf("%#v", time.Now())
    newR := &pb.MessageRequest{Message: r.GetMessage() + ": " + time.Now().Format("01-02 15:04")}
    s.requests = append(s.requests, newR)
    return &pb.MessageResponse{Message: r.GetMessage()}, nil
}

func main_grpc() {
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    pb.RegisterMessengerServer(s, &server{})
    reflection.Register(s)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}