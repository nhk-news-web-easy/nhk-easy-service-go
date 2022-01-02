package gateway

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/nhk-news-web-easy/nhk-easy-service-proto"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"strings"
)

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

func ProvideHTTP(endpoint string, grpcServer *grpc.Server) *http.Server {
	ctx := context.Background()
	//获取证书

	//添加证书
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	//新建gwmux，它是grpc-gateway的请求复用器。它将http请求与模式匹配，并调用相应的处理程序。
	gwmux := runtime.NewServeMux()
	//将服务的http处理程序注册到gwmux。处理程序通过endpoint转发请求到grpc端点
	err := pb.RegisterNhkServiceHandlerFromEndpoint(ctx, gwmux, endpoint, dopts)
	if err != nil {
		log.Fatalf("Register Endpoint err: %v", err)
	}
	//新建mux，它是http的请求复用器
	mux := http.NewServeMux()
	//注册gwmux
	mux.Handle("/", gwmux)
	//注册swagger
	log.Println(endpoint + " HTTP.Listing whth TLS and token...")
	return &http.Server{
		Addr:    endpoint,
		Handler: grpcHandlerFunc(grpcServer, mux),
	}
}
