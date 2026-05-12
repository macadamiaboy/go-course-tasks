package main

import (
	"complex-task/db"
	"complex-task/handlers"
	"complex-task/service"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
)

// importing this data from env
const username string = "postgres"
const password string = "password"
const host string = "localhost"
const port int = 5432
const dbName string = "test_db"

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
	httpRequestsInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of HTTP requests being processed.",
		},
	)
)

type Middleware func(http.Handler) http.Handler

type TraceHandler struct {
	slog.Handler
}

// status interception
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// route pattern getter
func getRoutePattern(r *http.Request) string {
	if r.Pattern != "" {
		return r.Pattern
	}

	return "unknown"
}

// mw chain
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// logger stuff
func newLogger(service, env, version string) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return slog.New(&TraceHandler{handler}).With("service", service, "env", env, "version", version)
}

func loggingMiddleware(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sRecorder := statusRecorder{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(&sRecorder, r)

			logger.Info("http request",
				"method", r.Method,
				"path", getRoutePattern(r),
				"status", sRecorder.status,
			)
		})
	}
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpRequestsInFlight.Inc()
		defer httpRequestsInFlight.Dec()

		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		status := strconv.Itoa(rec.status)
		route := getRoutePattern(r)

		httpRequestsTotal.WithLabelValues(r.Method, route, status).Inc()
		httpRequestDuration.WithLabelValues(r.Method, route).Observe(time.Since(start).Seconds())
	})
}

// OTel
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
			semconv.ServiceNameKey.String("5.6-complex-service"),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp, nil
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

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?sslmode=disable",
		username,
		password,
		host,
		port,
		dbName,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		logger.Error("database connection failed", "error", err)
	}
	defer pool.Close()

	userRepo := db.NewUserRepository(pool)
	tokenRepo := db.NewRefreshTokenRepository(pool)
	tokenService := service.NewTokenService(*tokenRepo, *userRepo)
	tokenHandler := handlers.NewTokenHandler(tokenService, logger)

	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration, httpRequestsInFlight)
	tp, _ := initTracer()
	defer tp.Shutdown(context.Background())

	mux := http.NewServeMux()

	mux.Handle("POST /refresh", Chain(http.HandlerFunc(tokenHandler.Refresh), loggingMiddleware(logger), MetricsMiddleware))
	mux.Handle("GET /validate", Chain(http.HandlerFunc(tokenHandler.Validate), loggingMiddleware(logger), MetricsMiddleware))
	mux.Handle("POST /revoke", Chain(http.HandlerFunc(tokenHandler.Revoke), loggingMiddleware(logger), MetricsMiddleware))

	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("GET /health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))

	logger.Info("started", "has happened", "server started")
	_ = http.ListenAndServe(":8080", mux)
}
