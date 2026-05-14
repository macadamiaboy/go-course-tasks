// Задание: Базовый трейсинг с OpenTelemetry
//
// Подключи OpenTelemetry SDK с stdout-экспортером (для dev).
// Создай span для HTTP-запроса и вложенный span для операции сервиса.
//
// Для запуска нужны зависимости:
//   go get go.opentelemetry.io/otel
//   go get go.opentelemetry.io/otel/sdk/trace
//   go get go.opentelemetry.io/otel/exporters/stdout/stdouttrace
//
// Ожидаемый результат (в stdout — JSON с данными трейса):
//   {"Name":"process-token","SpanContext":{...},...}
//   {"Name":"GET /api/v1/tokens","SpanContext":{...},...}

package main

import (
	"context"
	"fmt"
	"log"
)

// TODO: реализуй initTracer() (*sdktrace.TracerProvider, error)
// Шаги:
//   1. Создай stdout exporter: stdouttrace.New(stdouttrace.WithPrettyPrint())
//   2. Создай TracerProvider с BatchSpanProcessor
//   3. Установи глобальный провайдер: otel.SetTracerProvider(tp)
//   4. Верни TracerProvider

// TODO: вызови initTracer() в main
// Добавь defer tp.Shutdown(ctx) для сброса оставшихся span'ов

// TODO: создай tracer: otel.Tracer("token-service")

// TODO: в функции handleRequest(ctx context.Context, userID string):
//   1. Создай span "GET /api/v1/tokens": ctx, span := tracer.Start(ctx, "GET /api/v1/tokens")
//   2. Внутри вызови processToken(ctx, userID)
//   3. Закрой span: defer span.End()

// TODO: в функции processToken(ctx context.Context, userID string):
//   1. Создай вложенный span "process-token"
//   2. Добавь атрибут: span.SetAttributes(attribute.String("user.id", userID))
//   3. Закрой span

func main() {
	ctx := context.Background()

	// TODO: раскомментируй после реализации initTracer:
	// tp, err := initTracer()
	// if err != nil {
	//     log.Fatal(err)
	// }
	// defer tp.Shutdown(ctx)

	// TODO: вызови handleRequest(ctx, "user-123")

	_ = ctx
	_ = log.Fatal
	fmt.Println("TODO: implement OpenTelemetry tracing")
}
