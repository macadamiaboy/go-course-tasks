package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	// TODO: save code into r.status before forwarding
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// TODO: Implement LoggingMiddleware.
// Return an http.Handler that:
// 1. Records start time (time.Now())
// 2. Wraps w in statusRecorder with default status 200
// 3. Calls next.ServeHTTP with the recorder
// 4. Logs via logger: method, path, status, duration_ms

func LoggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)
		logger.Info("incoming request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})).With("service", "api-gateway")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status":"ok"}`)
	})
	mux.HandleFunc("GET /api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[{"id":1,"name":"Alice"}]`)
	})

	handler := LoggingMiddleware(logger, mux)

	addr := ":8080"
	logger.Info("server starting", "addr", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		logger.Error("server failed", "error", err)
	}
}
