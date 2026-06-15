package observability

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type HTTPMetrics struct {
	RequestsTotal    *prometheus.CounterVec
	RequestsDuration *prometheus.HistogramVec
}

type AuthMetrics struct {
	TokensIssued *prometheus.CounterVec
}

func NewHTTPMetrics() *HTTPMetrics {
	return &HTTPMetrics{
		RequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests.",
			},
			[]string{"method", "path", "status"},
		),
		RequestsDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request latency in seconds.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
	}
}

func NewAuthMetrics() *AuthMetrics {
	return &AuthMetrics{
		TokensIssued: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "tokens_issued_total",
				Help: "Total number of issued tokens.",
			},
			[]string{"token_type", "status", "reason"},
		),
	}
}
