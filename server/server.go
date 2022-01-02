package server

import (
	"context"
	pb "github.com/nhk-news-web-easy/nhk-easy-service-proto"
)

type GrpcServer struct {
	pb.UnimplementedNhkServiceServer
}

func (server *GrpcServer) GetNews(context context.Context, request *pb.NewsRequest) (*pb.NewsReply, error) {
	return &pb.NewsReply{
		Error: nil,
		News:  nil,
	}, nil
}
