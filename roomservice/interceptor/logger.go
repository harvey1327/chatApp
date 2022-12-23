package interceptor

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

// UnaryLoggerInterceptor Intercepts a unaryRequest and logs the request, response, error
func UnaryLoggerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Printf("Request: %s", info.FullMethod)
		resp, err := handler(ctx, req)
		if err != nil {
			log.Printf("Error: %s, %s", info.FullMethod, err)
		} else {
			log.Printf("Response: %s", info.FullMethod)
		}
		return resp, err
	}
}
