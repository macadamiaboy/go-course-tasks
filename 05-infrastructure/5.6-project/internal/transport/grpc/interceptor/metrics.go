package interceptor

import (
	"context"
	"strings"
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

		service, method := parseFullMethod(info.FullMethod)

		totalReqsCounter.WithLabelValues(service, method, statusCode).Inc()
		durationHistogram.WithLabelValues(service, method).Observe(time.Since(start).Seconds())

		return resp, err
	}
}

func parseFullMethod(fullMethod string) (service string, method string) {
	trimmedFullMethod := strings.TrimPrefix(fullMethod, "/")
	parts := strings.SplitN(trimmedFullMethod, "/", 2)

	if len(parts) != 2 {
		return "unknown", fullMethod
	}
	return parts[0], parts[1]
}
