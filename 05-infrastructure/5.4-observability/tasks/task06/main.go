// Задание: Интеграционная мини-задача — observability для Token Service
//
// Собери минимальную observability-конфигурацию:
//   - структурированные логи slog (JSON, с service/env/version)
//   - Prometheus-метрики на /metrics
//   - OpenTelemetry-трейсы для сценария "выдача токена"
//
// Требования:
//   - минимум 3 полезные метрики
//   - минимум 2 span на один запрос
//   - лог ошибки содержит trace_id
//
// Зависимости:
//   go get github.com/prometheus/client_golang/prometheus
//   go get github.com/prometheus/client_golang/prometheus/promhttp
//   go get go.opentelemetry.io/otel
//   go get go.opentelemetry.io/otel/sdk/trace
//
// Ожидаемый результат:
//   $ go run main.go
//   {"level":"INFO","msg":"server started","service":"token-service","addr":":8080"}
//
//   $ curl -X POST http://localhost:8080/api/v1/tokens -d '{"user_id":"u-1"}'
//   {"token_id":"tok-...","user_id":"u-1"}
//
//   $ curl http://localhost:8080/metrics
//   token_service_http_requests_total{method="POST",route="/api/v1/tokens",status="201"} 1
//   token_service_tokens_issued_total 1

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

// TODO: инициализируй Prometheus-метрики:
//   httpRequestsTotal   — counter vec по method, route, status
//   httpRequestDuration — histogram vec по method, route
//   tokensIssued        — counter (сколько токенов выдано)

// TODO: инициализируй OTel TracerProvider с stdout-экспортером

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// TODO: реализуй metricsMiddleware — оборачивает handler и записывает метрики
// после каждого запроса: httpRequestsTotal и httpRequestDuration

// TODO: реализуй handleIssueToken(w, r):
//   1. Декодируй {"user_id":"..."}
//   2. Создай span "issue-token" с вложенным span "validate-user"
//   3. Инкрементируй tokensIssued
//   4. Залогируй с trace_id
//   5. Верни {"token_id":"tok-<ts>","user_id":"..."}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		"service", "token-service",
		"env", "development",
		"version", "0.1.0",
	)

	ctx := context.Background()

	// TODO: инициализируй OTel TracerProvider
	// TODO: зарегистрируй Prometheus-метрики (prometheus.MustRegister)

	mux := http.NewServeMux()

	// TODO: добавь маршруты:
	// mux.Handle("GET /health", metricsMiddleware(http.HandlerFunc(healthHandler)))
	// mux.Handle("POST /api/v1/tokens", metricsMiddleware(http.HandlerFunc(handleIssueToken)))
	// mux.Handle("GET /metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	logger.Info("server started", "addr", ":8080")
	_ = ctx
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server failed", "error", err)
	}
	fmt.Println("TODO: implement full observability stack")
}
