// Проект: Token Service — полная реализация
//
// Поэтапная реализация микросервиса аутентификации.
//
// Этапы:
//   Stage 1 — API Contract:    определи .proto + HTTP-контракт
//   Stage 2 — DB Design:       схема users + refresh_tokens, миграции
//   Stage 3 — Data Layer:      pgxpool + репозитории + транзакция RegisterWithToken
//   Stage 4 — HTTP + Middleware: роуты, auth-middleware, метрики-middleware
//   Stage 5 — gRPC Service:    TokenService сервер + клиент
//   Stage 6 — Observability:   slog JSON, Prometheus, OpenTelemetry
//   Stage 7 — Integration:     HTTP Auth API → gRPC TokenService
//
// Структура проекта:
//   token-service/
//     main.go                  ← точка входа (этот файл)
//     go.mod
//     proto/
//       token.proto            ← Stage 1
//     sql/
//       migrations/            ← Stage 2
//         001_users.up.sql
//         001_users.down.sql
//         002_tokens.up.sql
//         002_tokens.down.sql
//       queries/
//         users.sql
//         tokens.sql
//     internal/
//       db/                    ← Stage 3 (sqlc-generated or manual)
//         user_repo.go
//         token_repo.go
//       service/               ← Stage 5
//         token_service.go
//       middleware/            ← Stage 4
//         auth.go
//         logging.go
//         metrics.go
//       observability/         ← Stage 6
//         logger.go
//         metrics.go
//         tracer.go
//       transport/             ← Stage 7
//         http.go
//         grpc.go
//
// Быстрый старт (docker-compose для PostgreSQL):
//   docker run --rm -e POSTGRES_USER=dev -e POSTGRES_PASSWORD=dev -e POSTGRES_DB=devdb -p 5432:5432 postgres:16
//   export DATABASE_URL="postgres://dev:dev@localhost:5432/devdb"
//   go run main.go
//
// Зависимости (установи по мере реализации этапов):
//   go get github.com/jackc/pgx/v5
//   go get github.com/prometheus/client_golang/prometheus
//   go get github.com/prometheus/client_golang/prometheus/promhttp
//   go get go.opentelemetry.io/otel
//   go get go.opentelemetry.io/otel/sdk/trace
//   go get go.opentelemetry.io/otel/exporters/stdout/stdouttrace

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// =============================================================================
// Stage 1: API Contract
// =============================================================================
//
// HTTP Endpoints:
//   POST /auth/register  {"email":"...","password":"..."}   → 201 {"user_id":1}
//   POST /auth/login     {"email":"...","password":"..."}   → 200 {"access_token":"...","refresh_token":"..."}
//   POST /auth/refresh   {"refresh_token":"..."}            → 200 {"access_token":"...","refresh_token":"..."}
//   POST /auth/logout    {"refresh_token":"..."}            → 200 {"revoked":true}
//   GET  /auth/me        Header: Authorization: Bearer ...  → 200 {"user_id":1,"email":"..."}
//
// gRPC proto (proto/token.proto):
//   service TokenService {
//     rpc IssueToken(IssueTokenRequest)     returns (IssueTokenResponse);
//     rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
//     rpc RevokeToken(RevokeTokenRequest)   returns (RevokeTokenResponse);
//   }
//   message IssueTokenRequest  { string user_id = 1; }
//   message IssueTokenResponse { string token_id = 1; string user_id = 2; }
//   message ValidateTokenRequest  { string token_id = 1; }
//   message ValidateTokenResponse { string user_id = 1; bool valid = 2; }
//   message RevokeTokenRequest  { string token_id = 1; }
//   message RevokeTokenResponse { bool revoked = 1; }

// =============================================================================
// Stage 2: DB Schema (реализуй в sql/migrations/)
// =============================================================================
//
// 001_users.up.sql:
//   CREATE TABLE users (
//     id            BIGSERIAL PRIMARY KEY,
//     email         TEXT NOT NULL UNIQUE,
//     password_hash TEXT NOT NULL,
//     created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
//   );
//
// 002_tokens.up.sql:
//   CREATE TABLE refresh_tokens (
//     id         BIGSERIAL PRIMARY KEY,
//     user_id    BIGINT NOT NULL REFERENCES users(id),
//     token_hash TEXT NOT NULL UNIQUE,
//     expires_at TIMESTAMPTZ NOT NULL,
//     revoked_at TIMESTAMPTZ,
//     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
//   );
//   CREATE INDEX ON refresh_tokens(token_hash);

// =============================================================================
// Stage 3: Data Layer
// =============================================================================

// TODO: реализуй internal/db/user_repo.go — UserRepository (pgxpool)
//   CreateUser(ctx, email, passwordHash string) (User, error)
//   GetUserByEmail(ctx, email string) (User, error)

// TODO: реализуй internal/db/token_repo.go — TokenRepository (pgxpool)
//   CreateRefreshToken(ctx, userID int64, tokenHash string, expiresAt time.Time) (RefreshToken, error)
//   GetActiveToken(ctx, tokenHash string) (RefreshToken, error)
//   RevokeToken(ctx, tokenHash string) error

// TODO: реализуй RegisterWithToken(ctx, pool, email, pwdHash, tokenHash) error
// в одной транзакции: CreateUser + CreateRefreshToken

// =============================================================================
// Stage 4: HTTP Server + Middleware
// =============================================================================

// TODO: реализуй middleware/logging.go — LoggingMiddleware (метод, путь, статус, время)
// TODO: реализуй middleware/auth.go — AuthMiddleware (Bearer-токен → user_id в context)
// TODO: реализуй middleware/metrics.go — MetricsMiddleware (Prometheus counter + histogram)

// =============================================================================
// Stage 5: gRPC Service
// =============================================================================

// TODO: реализуй internal/service/token_service.go — TokenServiceServer
//   IssueToken, ValidateToken, RevokeToken с gRPC-статусами

// TODO: реализуй interceptors: LoggingInterceptor, AuthInterceptor, MetricsInterceptor

// =============================================================================
// Stage 6: Observability
// =============================================================================

// TODO: реализуй internal/observability/logger.go
//   InitLogger() *slog.Logger с полями service, env, version

// TODO: реализуй internal/observability/metrics.go
//   InitMetrics() — регистрирует httpRequestsTotal, httpRequestDuration, tokensIssued

// TODO: реализуй internal/observability/tracer.go
//   InitTracer(ctx) (*sdktrace.TracerProvider, error) — stdout или OTLP экспортер

// =============================================================================
// Stage 7: Integration
// =============================================================================

// TODO: HTTP-хендлеры вызывают gRPC TokenService:
//   handleRegister → RegisterWithToken → IssueToken
//   handleLogin    → GetUserByEmail + ValidatePassword → IssueToken
//   handleRefresh  → GetActiveToken → RevokeToken + IssueToken
//   handleLogout   → RevokeToken
//   handleMe       → ValidateToken

// =============================================================================
// Временная заглушка — заменяется по мере реализации этапов
// =============================================================================

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		"service", "token-service",
		"env", "development",
		"version", "0.1.0",
	)

	ctx := context.Background()
	_ = ctx

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)

	// TODO: подключи DATABASE_URL и инициализируй pgxpool
	// TODO: инициализируй OTel TracerProvider
	// TODO: зарегистрируй Prometheus-метрики
	// TODO: добавь роуты /auth/register, /auth/login, /auth/refresh, /auth/logout, /auth/me
	// TODO: добавь GET /metrics для Prometheus

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	logger.Info("server started", "addr", ":8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
	fmt.Println("TODO: implement Token Service — see stage comments above")
}
