package nhk_easy_service_go

import (
	"context"
	"flag"
	"fmt"
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

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterNhkServiceServer(s, &server{})

	log.Printf("server listening at %v", listener.Addr())

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
