package main

import (
	"errors"
	"log/slog"
	"os"
)

// TODO: Implement newLogger(service, env string) *slog.Logger
// - Use slog.NewJSONHandler writing to os.Stdout
// - Set level to slog.LevelInfo
// - Attach base fields: "service" and "env"

func newLogger(service, env string) *slog.Logger {
	// TODO: implement
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	return slog.New(handler).With("service", service, "env", env)
}

func main() {
	logger := newLogger("token-service", "development")

	// TODO: Log service startup at Info level with a "version" field
	logger.Info("service started", "version", "v1.0.0")

	// TODO: Log a simulated incoming request at Info level
	// with fields: method, path, request_id
	logger.Info("incoming request",
		"method", "POST",
		"path", "/api/token",
		"request_id", "req-01",
	)

	// TODO: Log a simulated error at Error level
	// using errors.New("connection timeout") as the error field
	err := errors.New("connection timeout")
	logger.Error("database connection failed", "error", err)
}
