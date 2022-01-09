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
	port         = flag.Int("port", 8080, "The server port")
	dbDriverName = flag.String("dbDriverName", "mysql", "The driver name of database")
	dbUserName   = flag.String("dbUserName", "root", "The user name of database")
	dbPassword   = flag.String("dbPassword", "123456", "The password of database")
	dbHost       = flag.String("dbHost", "127.0.0.1", "The host of database")
	dbPort       = flag.Int("dbPort", 3306, "The port of database")
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

	dbConfig := db.DbConfig{
		DriverName: *dbDriverName,
		UserName:   *dbUserName,
		Password:   *dbPassword,
		Host:       *dbHost,
		Port:       *dbPort,
	}
	err = db.InitDb(dbConfig)

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
