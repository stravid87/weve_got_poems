package utils

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

// Custom Logger to log incoming requests
func UnaryInterceptor(logger *log.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		logger.Printf("Incoming request: %s", info.FullMethod)
		resp, err := handler(ctx, req)
		return resp, err
	}
}