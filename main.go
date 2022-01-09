package main

import (
	"flag"
	"fmt"
	"github.com/nhk-news-web-easy/nhk-easy-service-go/db"
	"github.com/nhk-news-web-easy/nhk-easy-service-go/gateway"
	"github.com/nhk-news-web-easy/nhk-easy-service-go/server"
	pb "github.com/nhk-news-web-easy/nhk-easy-service-proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	port = flag.Int("port", 8080, "The server port")
)

func main() {
	flag.Parse()

	address := fmt.Sprintf("localhost:%d", *port)
	listener, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNhkServiceServer(grpcServer, &server.GrpcServer{})

	httpServer := gateway.NewHttpServer(address, grpcServer)

	log.Printf("server listening at %v", listener.Addr())

	err = db.InitDb()

	if err != nil {
		log.Fatalf("failed to init db %v", err)
	}

	if err = httpServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	err = db.CloseDb()

	if err != nil {
		log.Fatalf("failed to close db %v", err)
	}
}
