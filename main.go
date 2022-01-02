package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/nhk-news-web-easy/nhk-easy-service-go/gateway"
	pb "github.com/nhk-news-web-easy/nhk-easy-service-proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	port = flag.Int("port", 5000, "The server port")
)

type server struct {
	pb.UnimplementedNhkServiceServer
}

func (server *server) GetNews(context context.Context, request *pb.NewsRequest) (*pb.NewsReply, error) {
	return &pb.NewsReply{
		Error: nil,
		News:  nil,
	}, nil
}

func main() {
	flag.Parse()

	address := fmt.Sprintf("localhost:%d", *port)
	listener, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterNhkServiceServer(s, &server{})

	httpServer := gateway.NewHttpServer(address, s)
	
	log.Printf("server listening at %v", listener.Addr())

	if err = httpServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
