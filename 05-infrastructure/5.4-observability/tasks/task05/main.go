// Задание: Корреляция логов и трейсов
//
// Добавь trace_id и span_id в slog-логи, если они присутствуют в context.
// По логу должно быть возможно найти соответствующий трейс.
//
// Ожидаемый результат (JSON-лог):
//   {"level":"INFO","msg":"handling request","trace_id":"abc123...","span_id":"def456...","user_id":"u-1"}
//   {"level":"ERROR","msg":"operation failed","trace_id":"abc123...","span_id":"def456...","error":"..."}

package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

// traceFields извлекает trace_id и span_id из context (если есть активный span OTel).
//
// TODO: реализуй функцию traceFields(ctx context.Context) []any
// Шаги:
//   1. Получи SpanContext: span := trace.SpanFromContext(ctx)
//   2. Если span.SpanContext().IsValid() — добавь поля trace_id и span_id
//   3. Верни срез []any{"trace_id", "...", "span_id", "..."}
//   4. Если span невалидный — верни nil
//
// Подсказка: requires go.opentelemetry.io/otel/trace
// span.SpanContext().TraceID().String()
// span.SpanContext().SpanID().String()

func traceFields(ctx context.Context) []any {
	// TODO: implement
	return nil
}

// logWithTrace логирует сообщение, автоматически добавляя trace_id/span_id из context.
//
// TODO: реализуй logWithTrace(ctx context.Context, logger *slog.Logger, level slog.Level, msg string, args ...any)
// Объедини args и traceFields(ctx), затем вызови logger.Log(ctx, level, msg, allArgs...)

func logWithTrace(ctx context.Context, logger *slog.Logger, level slog.Level, msg string, args ...any) {
	// TODO: implement
	logger.Log(ctx, level, msg, args...)
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx := context.Background()

	// TODO: инициализируй OTel TracerProvider (аналогично task04)
	// TODO: создай span и используй его context в logWithTrace

	// Пример использования (раскомментируй после реализации):
	// ctx, span := tracer.Start(ctx, "handle-issue-token")
	// defer span.End()
	// logWithTrace(ctx, logger, slog.LevelInfo, "handling request", "user_id", "u-1")
	// logWithTrace(ctx, logger, slog.LevelError, "operation failed", "error", "db timeout")

	_ = logger
	fmt.Println("TODO: implement trace correlation in logs")
}
