package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

type Middleware func(http.Handler) http.Handler

type DBTX interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}

type Queries struct {
	db DBTX
}

type TokenService struct {
	pool           *pgxpool.Pool
	refresh_tokens Queries
	users          Queries
}

type TraceHandler struct {
	slog.Handler
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func getRoutePattern(r *http.Request) string {
	if r.Pattern != "" {
		return r.Pattern
	}

	return "unknown"
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		status := strconv.Itoa(rec.status)
		route := getRoutePattern(r)

		httpRequestsTotal.WithLabelValues(r.Method, route, status).Inc()
		httpRequestDuration.WithLabelValues(r.Method, route).Observe(time.Since(start).Seconds())
	})
}

func newLogger(service, env, version string) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return slog.New(&TraceHandler{handler}).With("service", service, "env", env, "version", version)
}

// подключение экспортера
func initTracer() (*sdktrace.TracerProvider, error) {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("my-service"),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp, nil
}

// for HTTP req
func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

		tracer := otel.Tracer("http-layer")
		ctx, span := tracer.Start(ctx, fmt.Sprintf("%s %s", r.Method, r.Pattern),
			trace.WithAttributes(
				attribute.String("url.path", r.URL.Path),
				attribute.String("http.method", r.Method)),
		)
		defer span.End()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// for the service layer
func (s *TokenService) UpdateAccToken(ctx context.Context, id string) error {
	tracer := otel.Tracer("service-layer")
	ctx, span := tracer.Start(ctx, "TokenService.UpdateAccToken")
	defer span.End()

	return s.refresh_tokens.UpdateToken(ctx, id)
}

// for the DB layer
func (r *Queries) UpdateToken(ctx context.Context, id string) error {
	tracer := otel.Tracer("db-layer")
	_, span := tracer.Start(ctx, "DB.Exec.UpdateToken")
	defer span.End()

	span.SetAttributes(
		attribute.String("db.system", "postgresql"),
		attribute.String("db.statement", "UPDATE users SET ..."),
	)

	// имитация запроса
	time.Sleep(5 * time.Millisecond)
	return nil
}

func (h *TraceHandler) Handle(ctx context.Context, r slog.Record) error {
	if ctx == nil {
		return h.Handler.Handle(ctx, r)
	}

	span := trace.SpanFromContext(ctx)
	if spanCtx := span.SpanContext(); spanCtx.IsValid() {
		r.AddAttrs(
			slog.String("trace_id", spanCtx.TraceID().String()),
			slog.String("span_id", spanCtx.SpanID().String()),
		)
	}

	return h.Handler.Handle(ctx, r)
}

func main() {
	logger := newLogger("token-service", "development", "v1.0.0")

	logger.Info("incoming request",
		"method", "POST",
		"path", "/api/token",
		"request_id", "req-01",
	)

	err := errors.New("connection timeout")
	logger.Error("database connection failed", "error", err)

	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("GET /health", MetricsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})))

	_ = http.ListenAndServe(":8080", mux)
}
