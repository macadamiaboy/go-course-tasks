package interceptor

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func MetricsUnaryInterceptor(totalReqsCounter *prometheus.CounterVec, durationHistogram *prometheus.HistogramVec) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		start := time.Now()

		resp, err = handler(ctx, req)

		st, _ := status.FromError(err)
		statusCode := st.Code().String()

		totalReqsCounter.WithLabelValues("gRPC", info.FullMethod, statusCode).Inc()
		durationHistogram.WithLabelValues("gRPC", info.FullMethod).Observe(time.Since(start).Seconds())

		return resp, err
	}
}
