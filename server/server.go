package server

import (
	"context"
	"github.com/nhk-news-web-easy/nhk-easy-service-go/service"
	pb "github.com/nhk-news-web-easy/nhk-easy-service-proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)

type GrpcServer struct {
	pb.UnimplementedNhkServiceServer
}

func (server *GrpcServer) GetNews(context context.Context, request *pb.NewsRequest) (*pb.NewsReply, error) {
	news, err := service.GetNews()

	if err != nil {
		log.Printf("failed to get news %v", err)

		return &pb.NewsReply{
			Error: &pb.Error{
				Code:    500,
				Message: "internal server error",
			},
			News: nil,
		}, nil
	}

	result := make([]*pb.News, len(news))

	for i, n := range news {
		result[i] = &pb.News{
			NewsId:          n.NewsId,
			Title:           n.Title,
			TitleWithRuby:   n.TitleWithRuby,
			OutlineWithRuby: n.OutlineWithRuby,
			Body:            n.Body,
			Url:             n.Url,
			M3U8Url:         n.M3u8Url,
			ImageUrl:        n.ImageUrl,
			PublishedAtUtc:  timestamppb.New(n.PublishedAtUtc),
		}
	}

	return &pb.NewsReply{
		Error: nil,
		News:  result,
	}, nil
}
