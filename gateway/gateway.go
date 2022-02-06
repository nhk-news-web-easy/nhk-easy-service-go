package gateway

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nhk-news-web-easy/nhk-easy-service-go/marshaler"
	pb "github.com/nhk-news-web-easy/nhk-easy-service-proto"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net/http"
	"strings"
)

func grpcHandler(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

func NewHttpServer(endpoint string, grpcServer *grpc.Server) *http.Server {
	ctx := context.Background()
	dialOption := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	gatewayMux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &marshaler.CustomMarshaler{
		runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				AllowPartial:    true,
				UseEnumNumbers:  true,
				EmitUnpopulated: true,
			},
		},
	}))
	err := pb.RegisterNhkServiceHandlerFromEndpoint(ctx, gatewayMux, endpoint, dialOption)

	if err != nil {
		log.Fatalf("failed to register endpoint: %v", err)
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/", gatewayMux)

	return &http.Server{
		Addr:    endpoint,
		Handler: grpcHandler(grpcServer, httpMux),
	}
}
