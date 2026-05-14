// Задание: Интеграционная мини-задача — публичный и защищённый endpoint
//
// Собери API с двумя маршрутами:
//   GET /health       — публичный, без авторизации
//   GET /api/v1/me    — только с валидным Bearer-токеном
//
// Требования:
//   - цепочка middleware только на защищённом endpoint
//   - структурированное логирование запросов (slog)
//   - единый формат 401 ошибки
//   - замена верификатора без изменения handler'ов
//
// Ожидаемый результат:
//   $ curl http://localhost:8080/health
//   {"status":"ok"}
//
//   $ curl http://localhost:8080/api/v1/me
//   {"error":"authorization header is empty"}
//
//   $ curl -H "Authorization: Bearer valid-token" http://localhost:8080/api/v1/me
//   {"user_id":"user-123","role":"admin"}

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

// TODO: реализуй LoggingMiddleware(logger *slog.Logger) Middleware
// Логируй: method, path, status, duration_ms

// TODO: реализуй AuthMiddleware(verifier TokenVerifier) Middleware
// Извлекай Bearer-токен, верифицируй, клади claims в context
// На ошибке — 401 {"error":"<msg>"}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: достань claims из r.Context() и верни {"user_id":"...","role":"..."}
	_ = context.Background // подсказка
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	verifier := &mockVerifier{}

	mux := http.NewServeMux()

	// TODO: зарегистрируй маршруты:
	// /health — только с LoggingMiddleware
	// /api/v1/me — с Chain(LoggingMiddleware + AuthMiddleware)

	_ = strings.CutPrefix // убери после реализации
	_ = time.Now          // убери после реализации

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	logger.Info("server started", "addr", ":8080")
	_ = verifier // убери после реализации
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("server failed", "error", err)
	}
}
