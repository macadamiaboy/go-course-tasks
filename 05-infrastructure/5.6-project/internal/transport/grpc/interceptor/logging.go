package interceptor

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggingUnaryInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		start := time.Now()

		resp, err = handler(ctx, req)

		st, _ := status.FromError(err)
		statusCode := st.Code().String()

		logger.Info("http request",
			"method", info.FullMethod,
			"status", statusCode,
			"duration", time.Since(start),
		)

		return resp, err
	}
}
