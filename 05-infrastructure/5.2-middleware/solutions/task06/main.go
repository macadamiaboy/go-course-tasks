package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type Claims struct {
	UserID string
	Role   string
}

type contextKey string

const claimsKey contextKey = "claims"

type TokenVerifier interface {
	Verify(token string) (Claims, error)
}

type mockVerifier struct{}

func (v *mockVerifier) Verify(token string) (Claims, error) {
	if token == "valid-token" {
		return Claims{UserID: "user-123", Role: "admin"}, nil
	}
	return Claims{}, errors.New("invalid token")
}

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func LoggingMiddleware(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rec, r)
			logger.Info("http request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.status,
				"duration_ms", time.Since(start).Milliseconds(),
			)
		})
	}
}

func AuthMiddleware(verifier TokenVerifier) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "authorization header is empty"})
				return
			}
			token, ok := strings.CutPrefix(header, "Bearer ")
			if !ok || strings.TrimSpace(token) == "" {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid authorization header format"})
				return
			}
			claims, err := verifier.Verify(strings.TrimSpace(token))
			if err != nil {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
				return
			}
			ctx := context.WithValue(r.Context(), claimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(claimsKey).(Claims)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"user_id": claims.UserID,
		"role":    claims.Role,
	})
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	verifier := &mockVerifier{}

	mux := http.NewServeMux()
	mux.Handle("GET /health", Chain(
		http.HandlerFunc(healthHandler),
		LoggingMiddleware(logger),
	))
	mux.Handle("GET /api/v1/me", Chain(
		http.HandlerFunc(meHandler),
		LoggingMiddleware(logger),
		AuthMiddleware(verifier),
	))

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	logger.Info("server started", "addr", ":8080")
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("server failed", "error", err)
	}
}
